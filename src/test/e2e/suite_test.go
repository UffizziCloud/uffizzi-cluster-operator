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

package e2e

import (
	"context"
	uffizziv1alpha1 "github.com/UffizziCloud/uffizzi-cluster-operator/src/api/v1alpha1"
	"github.com/UffizziCloud/uffizzi-cluster-operator/src/controllers/etcd"
	"github.com/UffizziCloud/uffizzi-cluster-operator/src/controllers/uffizzicluster"
	"github.com/UffizziCloud/uffizzi-cluster-operator/src/pkg/utils/exec"
	fluxhelmv2beta1 "github.com/fluxcd/helm-controller/api/v2beta1"
	fluxsourcev1beta2 "github.com/fluxcd/source-controller/api/v1beta2"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"k8s.io/utils/pointer"
	"math/rand"
	"os"
	"path/filepath"
	ctrl "sigs.k8s.io/controller-runtime"
	ctrlcfg "sigs.k8s.io/controller-runtime/pkg/config/v1alpha1"
	"testing"

	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	//+kubebuilder:scaffold:imports
)

// These test use Ginkgo (BDD-style Go testing framework). Refer to
// http://onsi.github.io/ginkgo/ to learn more about Ginkgo.

type UffizziClusterE2E struct {
	IsTainted          bool
	UseExistingCluster bool
	K8SManager         ctrl.Manager
}

var (
	cfg       *rest.Config
	k8sClient client.Client
	testEnv   *envtest.Environment
	ctx       context.Context
	cancel    context.CancelFunc
	e2e       UffizziClusterE2E
)

func TestAPIs(t *testing.T) {
	RegisterFailHandler(Fail)

	RunSpecs(t, "Controller Suite")
}

var _ = BeforeSuite(func() {
	logf.SetLogger(zap.New(zap.WriteTo(GinkgoWriter), zap.UseDevMode(true)))
	e2e = UffizziClusterE2E{

		IsTainted:          getTaintedTestClusterEnvVar(),
		UseExistingCluster: getEnvtestRemoteEnvVar(),
	}

	By("bootstrapping test environment")
	testEnv = &envtest.Environment{
		CRDDirectoryPaths:     []string{filepath.Join("..", "..", "config", "crd", "bases")},
		ErrorIfCRDPathMissing: true,
		UseExistingCluster:    &e2e.UseExistingCluster,
	}

	var err error
	// cfg is defined in this file globally.
	cfg, err = testEnv.Start()
	Expect(err).NotTo(HaveOccurred())
	Expect(cfg).NotTo(BeNil())

	err = uffizziv1alpha1.AddToScheme(scheme.Scheme)
	Expect(err).NotTo(HaveOccurred())

	err = fluxhelmv2beta1.AddToScheme(scheme.Scheme)
	Expect(err).NotTo(HaveOccurred())

	err = fluxsourcev1beta2.AddToScheme(scheme.Scheme)
	Expect(err).NotTo(HaveOccurred())

	//+kubebuilder:scaffold:scheme

	k8sClient, err = client.New(cfg, client.Options{Scheme: scheme.Scheme})
	Expect(err).NotTo(HaveOccurred())
	Expect(k8sClient).NotTo(BeNil())

	// set the k8s manager which the tests will use
	e2e.K8SManager = NewTestK8SManager(1)
})

var _ = AfterSuite(func() {
	cancel()
	By("tearing down the test environment")
	err := testEnv.Stop()
	Expect(err).NotTo(HaveOccurred())
})

func testingSeed() {
	// Seed the test environment with some data
	rand.Seed(12345)
}

// Kubectl executes a Kubectl command.
func (e2e *UffizziClusterE2E) Kubectl(args ...string) (string, string, error) {
	// run kubectl with args and return stdout, stderr, and error
	return exec.CmdWithContext(context.TODO(), exec.PrintCfg(), "kubectl", args...)
}

func (e2e *UffizziClusterE2E) StartReconcilerWithArgs(concurrent int) {
	e2e.K8SManager = NewTestK8SManager(concurrent)
	err := (&uffizzicluster.UffizziClusterReconciler{
		Client: e2e.K8SManager.GetClient(),
		Scheme: e2e.K8SManager.GetScheme(),
	}).SetupWithManager(e2e.K8SManager)
	Expect(err).ToNot(HaveOccurred())

	err = (&etcd.UffizziClusterEtcdReconciler{
		Client: e2e.K8SManager.GetClient(),
		Scheme: e2e.K8SManager.GetScheme(),
	}).SetupWithManager(e2e.K8SManager)
	Expect(err).ToNot(HaveOccurred())

	defer GinkgoRecover()
	err = e2e.K8SManager.Start(ctx)
	Expect(err).ToNot(HaveOccurred(), "failed to run manager")

	go e2e.K8SManager.Start(ctx)
}

func NewTestK8SManager(concurrent int) ctrl.Manager {
	ctx, cancel = context.WithCancel(context.TODO())
	k8sManager, err := ctrl.NewManager(cfg, ctrl.Options{
		Scheme: scheme.Scheme,
		Controller: ctrlcfg.ControllerConfigurationSpec{
			GroupKindConcurrency: map[string]int{
				uffizziv1alpha1.SchemaGroupVersion.WithKind("UffizziCluster").GroupKind().String(): concurrent,
			},
			RecoverPanic: pointer.Bool(true),
		},
	})
	Expect(err).ToNot(HaveOccurred())
	return k8sManager
}

// CompareSlices compares two slices for equality. It returns true if both slices have the same elements in the same order.
func compareSlices[T comparable](slice1, slice2 []T) bool {
	if len(slice1) != len(slice2) {
		return false
	}
	for i, v := range slice1 {
		if v != slice2[i] {
			return false
		}
	}
	return true
}

func getEnvtestRemoteEnvVar() bool {
	if os.Getenv("ENVTEST_REMOTE") == "true" {
		return true
	}
	return false
}

func getTaintedTestClusterEnvVar() bool {
	if os.Getenv("E2E_ARG_IS_TAINTED") == "true" {
		return true
	}
	return false
}
