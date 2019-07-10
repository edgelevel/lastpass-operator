package lastpass

import (
	"context"
	"time"

	niqdevv1alpha1 "github.com/niqdev/lastpass-operator/pkg/apis/niqdev/v1alpha1"
	"github.com/niqdev/lastpass-operator/pkg/lastpass"
	"github.com/niqdev/lastpass-operator/pkg/utils"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

var log = logf.Log.WithName("controller_lastpass")

// Add creates a new LastPass Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileLastPass{client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("lastpass-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource LastPass
	err = c.Watch(&source.Kind{Type: &niqdevv1alpha1.LastPass{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// Watch for changes to secondary resource Secrets and requeue the owner LastPass
	err = c.Watch(&source.Kind{Type: &corev1.Secret{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &niqdevv1alpha1.LastPass{},
	})
	if err != nil {
		return err
	}

	return nil
}

// blank assignment to verify that ReconcileLastPass implements reconcile.Reconciler
var _ reconcile.Reconciler = &ReconcileLastPass{}

// ReconcileLastPass reconciles a LastPass object
type ReconcileLastPass struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client client.Client
	scheme *runtime.Scheme
}

// Reconcile reads that state of the cluster for a LastPass object and makes changes based on the state read
// and what is in the LastPass.Spec
// Note:
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcileLastPass) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
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
	instance := &niqdevv1alpha1.LastPass{}
	err := r.client.Get(context.TODO(), request.NamespacedName, instance)
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
		if err := controllerutil.SetControllerReference(instance, desired, r.scheme); err != nil {
			return reconcile.Result{}, err
		}

		// Check if this Secret already exists
		current := &corev1.Secret{}
		err = r.client.Get(context.TODO(), types.NamespacedName{Name: desired.Name, Namespace: desired.Namespace}, current)
		if err != nil && errors.IsNotFound(err) {
			reqLogger.Info("Creating a new Secret", "Secret.Namespace", desired.Namespace, "Secret.Name", desired.Name)
			err = r.client.Create(context.TODO(), desired)
			if err != nil {
				return reconcile.Result{}, err
			}
			// Secret created successfully - don't requeue
			continue
		} else if err != nil {
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
			err = r.client.Update(context.TODO(), desired)
			if err != nil {
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

// newSecretForCR creates a new secret
func newSecretForCR(cr *niqdevv1alpha1.LastPass, secret lastpass.LastPassSecret) *corev1.Secret {
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
