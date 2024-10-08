package main

import (
	"flag"
	"fmt"
	"os"
	goruntime "runtime"
	"text/template"

	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	metricsserver "sigs.k8s.io/controller-runtime/pkg/metrics/server"

	edgelevelcomv1alpha1 "github.com/edgelevel/lastpass-operator/api/v1alpha1"
	"github.com/edgelevel/lastpass-operator/internal/controller"
	"github.com/edgelevel/lastpass-operator/version"
	"github.com/rs/zerolog/log"
	// +kubebuilder:scaffold:imports
)

var (
	scheme   = runtime.NewScheme()
	setupLog = ctrl.Log.WithName("setup")
)

func init() {
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))

	utilruntime.Must(edgelevelcomv1alpha1.AddToScheme(scheme))
	// +kubebuilder:scaffold:scheme
}

func printVersion() {
	log.Info().Msgf("Go Version: %s", goruntime.Version())
	log.Info().Msgf(fmt.Sprintf("Go OS/Arch: %s/%s", goruntime.GOOS, goruntime.GOARCH))
	log.Info().Msgf(fmt.Sprintf("Version of lastpass-operator: %s", version.Version))
}

func main() {
	var metricsAddr string
	var secretNameTemplateStr string
	var enableLeaderElection bool
	flag.StringVar(&metricsAddr, "metrics-addr", ":8080", "The address the metric endpoint binds to.")
	flag.StringVar(&secretNameTemplateStr, "secret-name-template", "{{.LastPass.ObjectMeta.Name}}-{{.LastPassSecret.ID}}", "The go template to generate secrets name from LastPass and LastPassSecret objects.")
	flag.BoolVar(&enableLeaderElection, "enable-leader-election", false,
		"Enable leader election for controller manager. "+
			"Enabling this will ensure there is only one active controller manager.")
	flag.Parse()

	ctrl.SetLogger(zap.New(zap.UseDevMode(true)))

	printVersion()

	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme: scheme,
		Metrics: metricsserver.Options{
			BindAddress: metricsAddr,
		},
		LeaderElection:   enableLeaderElection,
		LeaderElectionID: "e9330328.edgelevel.com",
	})
	if err != nil {
		setupLog.Error(err, "unable to start manager")
		os.Exit(1)
	}

	secretNameTemplate, err := template.New("secretName").Parse(secretNameTemplateStr)
	if err != nil {
		panic(err)
	}

	if err = (&controller.LastPassReconciler{
		Client:             mgr.GetClient(),
		Log:                ctrl.Log.WithName("controllers").WithName("LastPass"),
		Scheme:             mgr.GetScheme(),
		SecretNameTemplate: secretNameTemplate,
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "LastPass")
		os.Exit(1)
	}
	if err = (&controller.LastPassGroupReconciler{
		Client:             mgr.GetClient(),
		Log:                ctrl.Log.WithName("controllers").WithName("LastPass"),
		Scheme:             mgr.GetScheme(),
		SecretNameTemplate: secretNameTemplate,
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "LastPassGroup")
		os.Exit(1)
	}
	// +kubebuilder:scaffold:builder

	setupLog.Info("starting manager")
	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		setupLog.Error(err, "problem running manager")
		os.Exit(1)
	}
}
