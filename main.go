/*
Copyright 2023.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"flag"
	uclusteruffizzicomv1alpha1 "github.com/UffizziCloud/uffizzi-cluster-operator/src/api/v1alpha1"
	"github.com/UffizziCloud/uffizzi-cluster-operator/src/controllers/etcd"
	"github.com/UffizziCloud/uffizzi-cluster-operator/src/controllers/uffizzicluster"
	"k8s.io/utils/pointer"
	"os"
	ctrlcfg "sigs.k8s.io/controller-runtime/pkg/config/v1alpha1"

	// Import all Kubernetes client auth plugins (e.g. Azure, GCP, OIDC, etc.)
	// to ensure that exec-entrypoint and run can make use of them.
	_ "k8s.io/client-go/plugin/pkg/client/auth"

	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"

	fluxhelmv2beta1 "github.com/fluxcd/helm-controller/api/v2beta1"
	fluxsourcev1beta2 "github.com/fluxcd/source-controller/api/v1beta2"
	//+kubebuilder:scaffold:imports
)

var (
	scheme   = runtime.NewScheme()
	setupLog = ctrl.Log.WithName("setup")
)

func init() {
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))

	utilruntime.Must(uclusteruffizzicomv1alpha1.AddToScheme(scheme))

	utilruntime.Must(fluxhelmv2beta1.AddToScheme(scheme))

	utilruntime.Must(fluxsourcev1beta2.AddToScheme(scheme))
	//+kubebuilder:scaffold:scheme
}

func main() {
	var (
		metricsAddr               string
		probeAddr                 string
		enableLeaderElection      bool
		concurrentReconciliations int
	)

	flag.StringVar(&metricsAddr, "metrics-bind-address", ":8080", "The address the metric endpoint binds to.")
	flag.StringVar(&probeAddr, "health-probe-bind-address", ":8081", "The address the probe endpoint binds to.")
	flag.BoolVar(&enableLeaderElection, "leader-elect", false,
		"Enable leader election for controller manager. "+
			"Enabling this will ensure there is only one active controller manager.")
	flag.IntVar(&concurrentReconciliations, "concurrent", 5, "The number of concurrent reconciles per controller.")
	opts := zap.Options{
		Development: true,
	}
	opts.BindFlags(flag.CommandLine)
	flag.Parse()

	ctrl.SetLogger(zap.New(zap.UseFlagOptions(&opts)))

	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme:                 scheme,
		MetricsBindAddress:     metricsAddr,
		Port:                   9443,
		HealthProbeBindAddress: probeAddr,
		LeaderElection:         enableLeaderElection,
		LeaderElectionID:       "bfbf93f2.uffizzi.com",
		Controller: ctrlcfg.ControllerConfigurationSpec{
			GroupKindConcurrency: map[string]int{
				uclusteruffizzicomv1alpha1.SchemaGroupVersion.WithKind("UffizziCluster").GroupKind().String(): concurrentReconciliations,
			},
			RecoverPanic: pointer.Bool(true),
		},
	})
	if err != nil {
		setupLog.Error(err, "unable to start manager")
		os.Exit(1)
	}
	// Setup UffizziClusterReconciler
	if err = (&uffizzicluster.UffizziClusterReconciler{
		Client: mgr.GetClient(),
		Scheme: mgr.GetScheme(),
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "UffizziCluster")
		os.Exit(1)
	}
	// Setup UffizziClusterEtcdReconciler
	if err = (&etcd.UffizziClusterEtcdReconciler{
		Client: mgr.GetClient(),
		Scheme: mgr.GetScheme(),
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "UffizziClusterEtcd")
		os.Exit(1)
	}
	//+kubebuilder:scaffold:builder

	if err := mgr.AddHealthzCheck("healthz", healthz.Ping); err != nil {
		setupLog.Error(err, "unable to set up health check")
		os.Exit(1)
	}
	if err := mgr.AddReadyzCheck("readyz", healthz.Ping); err != nil {
		setupLog.Error(err, "unable to set up ready check")
		os.Exit(1)
	}

	setupLog.Info("starting manager")
	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		setupLog.Error(err, "problem running manager")
		os.Exit(1)
	}
}
