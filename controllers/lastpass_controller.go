package controllers

import (
	"context"
	"time"

	edgelevelv1alpha1 "github.com/edgelevel/lastpass-operator/api/v1alpha1"
	"github.com/edgelevel/lastpass-operator/pkg/lastpass"
	"github.com/edgelevel/lastpass-operator/pkg/utils"
	"github.com/go-logr/logr"
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
)

var log = logf.Log.WithName("controller_lastpass")

// LastPassReconciler reconciles a LastPass object
type LastPassReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=edgelevel.com,resources=lastpasses,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=edgelevel.com,resources=lastpasses/status,verbs=get;update;patch

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the LastPass object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.6.4/pkg/reconcile
func (r *LastPassReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", req.Namespace, "Request.Name", req.Name)
	reqLogger.Info("Reconciling LastPass")

	// Check that the environment variables are defined or exit. See also lastpass-master-secret
	lastPassUsername := utils.GetEnvOrDie("LASTPASS_USERNAME")
	lastPassPassword := utils.GetEnvOrDie("LASTPASS_PASSWORD")

	// Check that "lpass" binary is available or exit
	lastpass.VerifyCliExistsOrDie()

	// Login to LastPass
	if err := lastpass.Login(lastPassUsername, lastPassPassword); err != nil {
		// Attempt login again, sometimes it fails even if the credentials are valid - requeue the request.
		return reconcile.Result{}, err
	}

	// Fetch the LastPass instance
	instance := &edgelevelv1alpha1.LastPass{}
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

	// Request LastPass secrets
	lastPassSecrets, err := lastpass.RequestSecrets(instance.Spec.SecretRef.Group, instance.Spec.SecretRef.Name)
	// Logout
	lastpass.Logout()
	if err != nil {
		// Error parsing the response - requeue the request.
		return reconcile.Result{}, err
	}

	for index := range lastPassSecrets {

		// Define a new Secret object
		desired := newSecretForCR(instance, lastPassSecrets[index])

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

	// Periodically reconcile the Custom Resource
	if instance.Spec.SyncPolicy.Enabled {
		return reconcile.Result{RequeueAfter: time.Second * instance.Spec.SyncPolicy.Refresh}, nil
	}

	// Reconcile only if something happens inside the cluster: ignore if the Secret changes externally
	return reconcile.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *LastPassReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&edgelevelv1alpha1.LastPass{}).
		Complete(r)
}

// newSecretForCR creates a new secret
func newSecretForCR(cr *edgelevelv1alpha1.LastPass, secret lastpass.LastPassSecret) *corev1.Secret {
	labels := map[string]string{
		"app": "lastpass-operator",
	}
	annotations := map[string]string{
		"id":              secret.ID,
		"group":           secret.Group,
		"name":            secret.Name,
		"fullname":        secret.Fullname,
		"lastModifiedGmt": secret.LastModifiedGmt,
		"lastTouch":       secret.LastTouch,
	}

	data := map[string]string{}
	if cr.Spec.SecretRef.WithUsername {
		data["USERNAME"] = secret.Username
	}
	if cr.Spec.SecretRef.WithPassword {
		data["PASSWORD"] = secret.Password
	}
	if cr.Spec.SecretRef.WithUrl {
		data["URL"] = secret.URL
	}
	if cr.Spec.SecretRef.WithNote {
		data["NOTE"] = secret.Note
	}

	return &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:        cr.Name + "-" + secret.ID,
			Namespace:   cr.Namespace,
			Labels:      labels,
			Annotations: annotations,
		},
		StringData: data,
	}
}
