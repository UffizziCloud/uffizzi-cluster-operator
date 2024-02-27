package e2e

import (
	context "context"
	"github.com/UffizziCloud/uffizzi-cluster-operator/src/api/v1alpha1"
	"github.com/UffizziCloud/uffizzi-cluster-operator/src/controllers/uffizzicluster"
	"github.com/UffizziCloud/uffizzi-cluster-operator/src/pkg/constants"
	"github.com/UffizziCloud/uffizzi-cluster-operator/src/test/util/conditions"
	"github.com/UffizziCloud/uffizzi-cluster-operator/src/test/util/diff"
	"github.com/UffizziCloud/uffizzi-cluster-operator/src/test/util/resources"
	"github.com/google/go-cmp/cmp"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func wrapUffizziClusterLifecycleTest(ctx context.Context, ns *v1.Namespace, uc *v1alpha1.UffizziCluster) {
	helmRelease := resources.GetHelmReleaseFromUffizziCluster(uc)
	helmRepo := resources.GetHelmRepositoryFromUffizziCluster(uc)
	defer Context("When deleting UffizziCluster", func() {
		It("Should delete the UffizziCluster", func() {
			By("By deleting the UffizziCluster")
			Expect(k8sClient.Delete(ctx, uc)).Should(Succeed())
		})
		It("Should delete the Namespace", func() {
			By("By deleting the Namespace")
			Expect(deleteTestNamespace(ns.Name)).Should(Succeed())
		})
	})

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
				if err := k8sClient.Get(ctx, resources.CreateNamespacedName(helmRelease.Name, ns.Name), helmRelease); err != nil {
					return false
				}
				return true
			})
			//
			By("Checking if the Loft HelmRepository was created")
			Eventually(func() bool {
				if err := k8sClient.Get(ctx, resources.CreateNamespacedName(constants.LOFT_HELM_REPO, ns.Name), helmRepo); err != nil {
					return false
				}
				return true
			})

		})

		It("Should initialize correctly", func() {
			expectedConditions := []metav1.Condition{}
			uffizziClusterNSN := resources.CreateNamespacedName(uc.Name, ns.Name)
			By("Check if UffizziCluster initializes correctly")
			Eventually(func() bool {
				if err := k8sClient.Get(ctx, uffizziClusterNSN, uc); err != nil {
					return false
				}
				expectedConditions = uffizzicluster.GetAllInitializingConditions()
				return conditions.ContainsAllConditions(expectedConditions, uc.Status.Conditions)
			}, timeout, pollingTimeout).Should(BeTrue())
			d := cmp.Diff(expectedConditions, uc.Status.Conditions)
			GinkgoWriter.Printf(diff.PrintWantGot(d))
		})

		It("Should be in a Ready State", func() {
			expectedConditions := []metav1.Condition{}
			uffizziClusterNSN := resources.CreateNamespacedName(uc.Name, ns.Name)
			By("Check if UffizziCluster has the correct Ready conditions")
			Eventually(func() bool {
				if err := k8sClient.Get(ctx, uffizziClusterNSN, uc); err != nil {
					return false
				}
				expectedConditions = uffizzicluster.GetAllReadyConditions()
				return conditions.ContainsAllConditions(expectedConditions, uc.Status.Conditions)
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
			uffizziClusterNSN := resources.CreateNamespacedName(uc.Name, ns.Name)
			By("Check if UffizziCluster has the correct Sleep conditions")
			Eventually(func() bool {
				if err := k8sClient.Get(ctx, uffizziClusterNSN, uc); err != nil {
					return false
				}
				return conditions.ContainsAllConditions(expectedConditions, uc.Status.Conditions)
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
			uffizziClusterNSN := resources.CreateNamespacedName(uc.Name, ns.Name)
			By("Check if UffizziCluster has the correct Sleep conditions")
			Eventually(func() bool {
				if err := k8sClient.Get(ctx, uffizziClusterNSN, uc); err != nil {
					return false
				}
				return conditions.ContainsAllConditions(expectedConditions, uc.Status.Conditions)
			}, timeout, pollingTimeout).Should(BeTrue())
			d := cmp.Diff(expectedConditions, uc.Status.Conditions)
			GinkgoWriter.Printf(diff.PrintWantGot(d))
		})
	})
}

func deleteTestNamespace(name string) error {
	return k8sClient.Delete(ctx, &v1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
	})
}
