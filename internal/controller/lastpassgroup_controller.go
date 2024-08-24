package controller

import (
	"bytes"
	"context"
	"slices"
	"text/template"
	"time"

	corev1 "k8s.io/api/core/v1"

	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	edgelevelv1alpha1 "github.com/edgelevel/lastpass-operator/api/v1alpha1"
	"github.com/edgelevel/lastpass-operator/pkg/lastpass"
	"github.com/edgelevel/lastpass-operator/pkg/utils"
	"github.com/go-logr/logr"
)

var (
	secretOwnerKey = ".metadata.controller"
	apiGVStr       = edgelevelv1alpha1.GroupVersion.String()
)

// LastPassGroupReconciler reconciles a LastPassGroup object
type LastPassGroupReconciler struct {
	client.Client
	Log                logr.Logger
	Scheme             *runtime.Scheme
	SecretNameTemplate *template.Template
}

//+kubebuilder:rbac:groups=edgelevel.com,resources=lastpassgroups,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=edgelevel.com,resources=lastpassgroups/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=edgelevel.com,resources=lastpassgroups/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the LastPassGroup object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.17.3/pkg/reconcile
func (r *LastPassGroupReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	var log = logf.Log.WithName("controller_lastpassgroup")
	reqLogger := log.WithValues("Request.Namespace", req.Namespace, "Request.Name", req.Name)
	reqLogger.Info("Reconciling LastPass")

	lastPassUsername := utils.GetEnvOrDie("LASTPASS_USERNAME")
	lastPassPassword := utils.GetEnvOrDie("LASTPASS_PASSWORD")

	// Check that "lpass" binary is available or exit
	lastpass.VerifyCliExistsOrDie()

	// Login to LastPass
	if err := lastpass.Login(lastPassUsername, lastPassPassword); err != nil {
		// Attempt login again, sometimes it fails even if the credentials are valid - requeue the request.
		return reconcile.Result{}, err
	}

	instance := &edgelevelv1alpha1.LastPassGroup{}
	err := r.Client.Get(context.TODO(), req.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}

	lastPassSecrets, err := lastpass.RequestSecretsGroup(instance.Spec.GroupRef.Group)
	if err != nil {
		// Error parsing the response - requeue the request.
		return reconcile.Result{}, err
	}

	desiredSecrets := r.newGroupSecretsForCR(instance, lastPassSecrets)

	existingSecretslist := &corev1.SecretList{}
	r.Client.List(context.TODO(), existingSecretslist, client.InNamespace(instance.Namespace), client.MatchingFields{
		secretOwnerKey: req.Name,
	})

	deletedSecrets := []corev1.Secret{}
	for _, sec := range existingSecretslist.Items {
		contains := slices.ContainsFunc(desiredSecrets, func(s *corev1.Secret) bool { return s.Name == sec.Name })
		if !contains {
			deletedSecrets = append(deletedSecrets, sec)
		}
	}

	for _, del := range deletedSecrets {
		reqLogger.Info("Deleting Secret", "Secret.Namespace", del.Namespace, "Secret.Name", del.Name)

		if err := r.Client.Delete(context.TODO(), &del); err != nil {
			reqLogger.Error(err, "Failed to delete Secret", "Secret.Namespace", del.Namespace, "Secret.Name", del.Name)
			return reconcile.Result{}, err
		}
	}

	for _, desired := range desiredSecrets {
		reqLogger.Info("Verify LastPassSecret", "Secret.Namespace", desired.Namespace, "Secret.Name", desired.Name)

		// Set LastPassSecret instance as the owner and controller
		if err := controllerutil.SetControllerReference(instance, desired, r.Scheme); err != nil {
			reqLogger.Error(err, "Failed to set LastPassSecret instance as the owner and controller")
			return reconcile.Result{}, err
		}

		// Check if this Secret already exists
		current := &corev1.Secret{}
		err = r.Client.Get(context.TODO(), types.NamespacedName{Name: desired.Name, Namespace: desired.Namespace}, current)
		if err != nil && errors.IsNotFound(err) {
			reqLogger.Info("Creating Secret", "Secret.Namespace", desired.Namespace, "Secret.Name", desired.Name)
			err = r.Client.Create(context.TODO(), desired)
			if err != nil {
				reqLogger.Error(err, "Failed to create Secret", "Secret.Namespace", desired.Namespace, "Secret.Name", desired.Name)
				return reconcile.Result{}, err
			}
			// Secret created successfully - don't requeue
			continue
		} else if err != nil {
			reqLogger.Error(err, "Failed to get Secret", "Secret.Namespace", desired.Namespace, "Secret.Name", desired.Name)
			return reconcile.Result{}, err
		}

		// Check if this Secret is changed
		if current.Annotations["lastModifiedGmt"] != desired.Annotations["lastModifiedGmt"] || current.Annotations["lastTouch"] != desired.Annotations["lastTouch"] {
			reqLogger.Info("Updating Secret",
				"Secret.Namespace", desired.Namespace,
				"Secret.Name", desired.Name,
				"Current:LastModifiedGmt", current.Annotations["lastModifiedGmt"],
				"Desired:LastModifiedGmt", desired.Annotations["lastModifiedGmt"],
				"Current:LastTouch", current.Annotations["lastTouch"],
				"Desired:LastTouch", desired.Annotations["lastTouch"])
			err = r.Client.Update(context.TODO(), desired)
			if err != nil {
				reqLogger.Error(err, "Failed to update Secret", "Secret.Namespace", desired.Namespace, "Secret.Name", desired.Name)
				return reconcile.Result{}, err
			}
			// Secret updated successfully - don't requeue
			continue
		}

		reqLogger.Info("Skip reconcile: Secret already exists and is up to date", "Secret.Namespace", current.Namespace, "Secret.Name", current.Name)
	}

	if instance.Spec.SyncPolicy.Enabled {
		return reconcile.Result{RequeueAfter: time.Second * instance.Spec.SyncPolicy.Refresh}, nil
	}

	return reconcile.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *LastPassGroupReconciler) SetupWithManager(mgr ctrl.Manager) error {

	if err := mgr.GetFieldIndexer().IndexField(context.Background(), &corev1.Secret{}, secretOwnerKey, func(rawObj client.Object) []string {
		owner := metav1.GetControllerOf(rawObj.(*corev1.Secret))
		if owner == nil || owner.APIVersion != apiGVStr || owner.Kind != "LastPassGroup" {
			return nil
		}

		return []string{owner.Name}
	}); err != nil {
		return err
	}

	return ctrl.NewControllerManagedBy(mgr).
		For(&edgelevelv1alpha1.LastPassGroup{}).
		Complete(r)
}

// newSecretForCR creates a new secret
func (r *LastPassGroupReconciler) newGroupSecretsForCR(cr *edgelevelv1alpha1.LastPassGroup, secrets []lastpass.LastPassSecret) []*corev1.Secret {
	labels := map[string]string{
		"app": "lastpass-operator",
	}

	desiredSecrets := []*corev1.Secret{}
	for _, secret := range secrets {
		annotations := map[string]string{
			"id":              secret.ID,
			"group":           secret.Group,
			"name":            secret.Name,
			"fullname":        secret.Fullname,
			"lastModifiedGmt": secret.LastModifiedGmt,
			"lastTouch":       secret.LastTouch,
		}

		data := map[string]string{}
		if cr.Spec.GroupRef.WithUsername {
			data["USERNAME"] = secret.Username
		}
		if cr.Spec.GroupRef.WithPassword {
			data["PASSWORD"] = secret.Password
		}
		if cr.Spec.GroupRef.WithUrl {
			data["URL"] = secret.URL
		}
		if cr.Spec.GroupRef.WithNote {
			data["NOTE"] = secret.Note
		}

		var secretName bytes.Buffer
		r.SecretNameTemplate.Execute(&secretName, struct {
			LastPass       *edgelevelv1alpha1.LastPassGroup
			LastPassSecret lastpass.LastPassSecret
		}{cr, secret})

		desiredSecrets = append(desiredSecrets, &corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:        secretName.String(),
				Namespace:   cr.Namespace,
				Labels:      labels,
				Annotations: annotations,
			},
			StringData: data,
		})
	}

	return desiredSecrets
}
