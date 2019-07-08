package lastpasssecret

import (
	"context"

	niqdevv1alpha1 "github.com/niqdev/lastpass-operator/pkg/apis/niqdev/v1alpha1"
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

	"github.com/niqdev/lastpass-operator/pkg/lastpass"
	"github.com/niqdev/lastpass-operator/pkg/utils"
)

var log = logf.Log.WithName("controller_lastpasssecret")

// Add creates a new LastPassSecret Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileLastPassSecret{client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("lastpasssecret-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource LastPassSecret
	err = c.Watch(&source.Kind{Type: &niqdevv1alpha1.LastPassSecret{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// Watch for changes to secondary resource Pods and requeue the owner LastPassSecret
	err = c.Watch(&source.Kind{Type: &corev1.Secret{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &niqdevv1alpha1.LastPassSecret{},
	})
	if err != nil {
		return err
	}

	return nil
}

// blank assignment to verify that ReconcileLastPassSecret implements reconcile.Reconciler
var _ reconcile.Reconciler = &ReconcileLastPassSecret{}

// ReconcileLastPassSecret reconciles a LastPassSecret object
type ReconcileLastPassSecret struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client client.Client
	scheme *runtime.Scheme
}

// Reconcile reads that state of the cluster for a LastPassSecret object and makes changes based on the state read
// and what is in the LastPassSecret.Spec
// TODO update secret if "last_modified_gmt" or "last_touch" change (?)
func (r *ReconcileLastPassSecret) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling LastPassSecret")

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

	// Fetch the LastPassSecret instance
	instance := &niqdevv1alpha1.LastPassSecret{}
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

	// Request secrets
	internalSecrets, err := lastpass.RequestSecrets(instance.Spec.ItemRef.Group, instance.Spec.ItemRef.Name)
	// Logout
	lastpass.Logout()
	if err != nil {
		// Error parsing the response - requeue the request.
		return reconcile.Result{}, err
	}

	for index := range internalSecrets {
		reqLogger.Info("Verify internal secret", "id", internalSecrets[index].ID)

		// Define a new Secret object
		secret := newSecretForCR(instance, internalSecrets[index])

		// Set LastPassSecret instance as the owner and controller
		if err := controllerutil.SetControllerReference(instance, secret, r.scheme); err != nil {
			// Interrupt and exit the loop (skip other secrets)
			return reconcile.Result{}, err
		}

		// Check if this Secret already exists
		found := &corev1.Secret{}
		err = r.client.Get(context.TODO(), types.NamespacedName{Name: secret.Name, Namespace: secret.Namespace}, found)
		if err != nil && errors.IsNotFound(err) {
			reqLogger.Info("Creating a new Secret", "Secret.Namespace", secret.Namespace, "Secret.Name", secret.Name)
			err = r.client.Create(context.TODO(), secret)
			if err != nil {
				return reconcile.Result{}, err
			}
			// Secret created successfully - don't requeue
			continue
		} else if err != nil {
			return reconcile.Result{}, err
		}

		// Secret already exists - don't requeue
		reqLogger.Info("Skip reconcile: Secret already exists", "Pod.Namespace", found.Namespace, "Secret.Name", found.Name)
	}

	return reconcile.Result{}, nil
}

func newSecretForCR(cr *niqdevv1alpha1.LastPassSecret, secret lastpass.InternalSecret) *corev1.Secret {
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
	if cr.Spec.ItemRef.WithUsername {
		data["USERNAME"] = secret.Username
	}
	if cr.Spec.ItemRef.WithPassword {
		data["PASSWORD"] = secret.Password
	}
	if cr.Spec.ItemRef.WithUrl {
		data["URL"] = secret.URL
	}
	if cr.Spec.ItemRef.WithNote {
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
