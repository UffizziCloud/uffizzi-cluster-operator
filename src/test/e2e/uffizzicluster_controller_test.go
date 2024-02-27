package e2e

import (
	"context"
	"github.com/UffizziCloud/uffizzi-cluster-operator/src/api/v1alpha1"
	"github.com/UffizziCloud/uffizzi-cluster-operator/src/controllers/uffizzicluster"
	"github.com/UffizziCloud/uffizzi-cluster-operator/src/pkg/constants"
	"github.com/UffizziCloud/uffizzi-cluster-operator/src/test/util/diff"
	fluxhelmv2beta1 "github.com/fluxcd/helm-controller/api/v2beta1"
	fluxsourcev1beta2 "github.com/fluxcd/source-controller/api/v1beta2"
	"github.com/google/go-cmp/cmp"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"strings"
)

var _ = Describe("UffizziCluster Controller", func() {
	ctx := context.Background()
	testingSeed()
	name := "basic"
	timeout := "1m"
	pollingTimeout := "100ms"

	ns := &v1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: strings.Join([]string{name, randomString(5)}, "-"),
		},
	}
	uc := &v1alpha1.UffizziCluster{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: ns.Name,
		},
	}
	helmRelease := &fluxhelmv2beta1.HelmRelease{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "uc-" + name,
			Namespace: ns.Name,
		},
	}
	helmRepo := &fluxsourcev1beta2.HelmRepository{
		ObjectMeta: metav1.ObjectMeta{
			Name:      constants.LOFT_HELM_REPO,
			Namespace: ns.Name,
		},
	}

	Context("When creating UffizziCluster", func() {
		It("Should create a UffizziCluster", func() {
			//
			By("By Creating Namespace for the UffizziCluster")
			Expect(k8sClient.Create(ctx, ns)).Should(Succeed())

			//
			By("By creating a new UffizziCluster")
			Expect(k8sClient.Create(ctx, uc)).Should(Succeed())
		})

		It("Should create a HelmRelease and HelmRepository", func() {
			//
			By("Checking if the HelmRelease was created")
			Eventually(func() bool {
				if err := k8sClient.Get(ctx, createNamespacesName(helmRelease.Name, ns.Name), helmRelease); err != nil {
					return false
				}
				return true
			})
			//
			By("Checking if the Loft HelmRepository was created")
			Eventually(func() bool {
				if err := k8sClient.Get(ctx, createNamespacesName(constants.LOFT_HELM_REPO, ns.Name), helmRepo); err != nil {
					return false
				}
				return true
			})

		})

		It("Should initialize correctly", func() {
			expectedConditions := []metav1.Condition{}
			uffizziClusterNSN := createNamespacesName(uc.Name, ns.Name)
			By("Check if UffizziCluster initializes correctly")
			Eventually(func() bool {
				if err := k8sClient.Get(ctx, uffizziClusterNSN, uc); err != nil {
					return false
				}
				expectedConditions = uffizzicluster.GetAllInitializingConditions()
				return containsAllConditions(expectedConditions, uc.Status.Conditions)
			}, timeout, pollingTimeout).Should(BeTrue())
			d := cmp.Diff(expectedConditions, uc.Status.Conditions)
			GinkgoWriter.Printf(diff.PrintWantGot(d))
		})

		It("Should be in a Ready State", func() {
			expectedConditions := []metav1.Condition{}
			uffizziClusterNSN := createNamespacesName(uc.Name, ns.Name)
			By("Check if UffizziCluster has the correct Ready conditions")
			Eventually(func() bool {
				if err := k8sClient.Get(ctx, uffizziClusterNSN, uc); err != nil {
					return false
				}
				expectedConditions = uffizzicluster.GetAllReadyConditions()
				return containsAllConditions(expectedConditions, uc.Status.Conditions)
			}, timeout, pollingTimeout).Should(BeTrue())
			d := cmp.Diff(expectedConditions, uc.Status.Conditions)
			GinkgoWriter.Printf(diff.PrintWantGot(d))
		})
	})

	Context("When putting a cluster to sleep", func() {
		It("Should put the cluster to sleep", func() {
			By("By putting the UffizziCluster to sleep")
			uc.Spec.Sleep = true
			Expect(k8sClient.Update(ctx, uc)).Should(Succeed())
		})

		It("Should be in a Sleep State", func() {
			expectedConditions := []metav1.Condition{
				uffizzicluster.APINotReady(),
				uffizzicluster.DataStoreNotReady(),
			}
			uffizziClusterNSN := createNamespacesName(uc.Name, ns.Name)
			By("Check if UffizziCluster has the correct Sleep conditions")
			Eventually(func() bool {
				if err := k8sClient.Get(ctx, uffizziClusterNSN, uc); err != nil {
					return false
				}
				return containsAllConditions(expectedConditions, uc.Status.Conditions)
			}, timeout, pollingTimeout).Should(BeTrue())
			d := cmp.Diff(expectedConditions, uc.Status.Conditions)
			GinkgoWriter.Printf(diff.PrintWantGot(d))
		})
	})

	Context("When waking a cluster up", func() {
		It("Should wake the cluster up", func() {
			By("By waking the UffizziCluster up")
			uc.Spec.Sleep = false
			Expect(k8sClient.Update(ctx, uc)).Should(Succeed())
		})

		It("Should be Awoken", func() {
			expectedConditions := []metav1.Condition{
				uffizzicluster.APIReady(),
				uffizzicluster.DataStoreReady(),
			}
			uffizziClusterNSN := createNamespacesName(uc.Name, ns.Name)
			By("Check if UffizziCluster has the correct Sleep conditions")
			Eventually(func() bool {
				if err := k8sClient.Get(ctx, uffizziClusterNSN, uc); err != nil {
					return false
				}
				return containsAllConditions(expectedConditions, uc.Status.Conditions)
			}, timeout, pollingTimeout).Should(BeTrue())
			d := cmp.Diff(expectedConditions, uc.Status.Conditions)
			GinkgoWriter.Printf(diff.PrintWantGot(d))
		})
	})

	Context("When deleting UffizziCluster", func() {
		It("Should delete the UffizziCluster", func() {
			By("By deleting the UffizziCluster")
			Expect(k8sClient.Delete(ctx, uc)).Should(Succeed())
		})
		It("Should delete the Namespace", func() {
			By("By deleting the Namespace")
			Expect(deleteTestNamespace(ns.Name)).Should(Succeed())
		})
	})
})

func createNamespacesName(name, namespace string) types.NamespacedName {
	return types.NamespacedName{
		Name:      name,
		Namespace: namespace,
	}
}

func deleteTestNamespace(name string) error {
	return k8sClient.Delete(ctx, &v1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
	})
}
