package lifecycle

import (
	"context"
	"github.com/UffizziCloud/uffizzi-cluster-operator/src/pkg/constants"
	"github.com/UffizziCloud/uffizzi-cluster-operator/src/test/util/conditions"
	"github.com/UffizziCloud/uffizzi-cluster-operator/src/test/util/resources"
	"github.com/fluxcd/pkg/apis/meta"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (td *LifecycleTestDefinition) Run(ctx context.Context) {
	var (
		ns = resources.CreateTestNamespace(td.Name)
		uc = resources.CreateTestUffizziCluster(td.Name, ns.Name)
	)
	uc.Spec = td.Spec
	var (
		k8sClient       = td.K8SClient
		timeout         = "10m"
		pollingTimeout  = "100ms"
		helmRelease     = resources.GetHelmReleaseFromUffizziCluster(uc)
		etcdHelmRelease = resources.GetETCDHelmReleaseFromUffizziCluster(uc)
		helmRepo        = resources.GetHelmRepositoryFromUffizziCluster(uc)
	)

	// defer deletion of the uffizzi cluster and namespace
	defer Context("When deleting UffizziCluster", func() {
		It("Should delete the UffizziCluster", func() {
			By("By deleting the UffizziCluster")
			Expect(k8sClient.Delete(ctx, uc)).Should(Succeed())
		})
		It("Should delete the Namespace", func() {
			By("By deleting the Namespace")
			Expect(k8sClient.Delete(ctx, ns)).Should(Succeed())
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
			By("Checking if the Loft HelmRepository was created and the repository is ready")
			Eventually(func() bool {
				if err := k8sClient.Get(ctx, resources.CreateNamespacedName(constants.LOFT_HELM_REPO, ns.Name), helmRepo); err != nil {
					return false
				}
				for _, c := range helmRepo.Status.Conditions {
					if c.Type == meta.ReadyCondition {
						return c.Status == metav1.ConditionTrue
					}
				}
				return true
			})
			//
			By("Checking if the HelmRelease was created")
			Eventually(func() bool {
				if err := k8sClient.Get(ctx, resources.CreateNamespacedName(helmRelease.Name, ns.Name), helmRelease); err != nil {
					return false
				}
				return true
			})
		})

		if uc.Spec.ExternalDatastore == constants.ETCD {
			It("Should create a Bitnami HelmRepository", func() {
				//
				By("Checking if the Bitnami HelmRepository was created")
				Eventually(func() bool {
					if err := k8sClient.Get(ctx, resources.CreateNamespacedName(constants.BITNAMI_HELM_REPO, ns.Name), helmRepo); err != nil {
						return false
					}
					return true
				})
			})
			It("Should create a HelmRelease for ETCD", func() {
				//
				By("Checking if the HelmRelease for ETCD was created")
				Eventually(func() bool {
					if err := k8sClient.Get(ctx, resources.CreateNamespacedName(etcdHelmRelease.Name, ns.Name), etcdHelmRelease); err != nil {
						return false
					}
					return true
				})
			})
		}

		// Initializing
		It("Should initialize correctly", func() {
			expectedConditions := td.ExpectedStatus.Initializing.Conditions
			uffizziClusterNSN := resources.CreateNamespacedName(uc.Name, ns.Name)
			By("Check if UffizziCluster initializes correctly")
			Eventually(func() bool {
				if err := k8sClient.Get(ctx, uffizziClusterNSN, uc); err != nil {
					return false
				}
				return conditions.ContainsAllConditions(expectedConditions, uc.Status.Conditions)
			}, timeout, pollingTimeout).Should(BeTrue())

			//GinkgoWriter.Printf(conditions.CreateConditionsCmpDiff(expectedConditions, uc.Status.Conditions))
		})

		It("Should be in a Ready State", func() {
			expectedConditions := td.ExpectedStatus.Ready.Conditions
			uffizziClusterNSN := resources.CreateNamespacedName(uc.Name, ns.Name)
			By("Check if UffizziCluster has the correct Ready conditions")
			Eventually(func() bool {
				if err := k8sClient.Get(ctx, uffizziClusterNSN, uc); err != nil {
					return false
				}
				return conditions.ContainsAllConditions(expectedConditions, uc.Status.Conditions)
			}, timeout, pollingTimeout).Should(BeTrue())

			//GinkgoWriter.Printf(conditions.CreateConditionsCmpDiff(expectedConditions, uc.Status.Conditions))
		})
	})

	Context("When putting a cluster to sleep", func() {
		It("Should put the cluster to sleep", func() {
			By("By putting the UffizziCluster to sleep")
			uc.Spec.Sleep = true
			Expect(k8sClient.Update(ctx, uc)).Should(Succeed())
		})

		It("Should be in a Sleep State", func() {
			expectedConditions := td.ExpectedStatus.Sleeping.Conditions
			uffizziClusterNSN := resources.CreateNamespacedName(uc.Name, ns.Name)
			By("Check if UffizziCluster has the correct Sleep conditions")
			Eventually(func() bool {
				if err := k8sClient.Get(ctx, uffizziClusterNSN, uc); err != nil {
					return false
				}
				return conditions.ContainsAllConditions(expectedConditions, uc.Status.Conditions)
			}, timeout, pollingTimeout).Should(BeTrue())

			//GinkgoWriter.Printf(conditions.CreateConditionsCmpDiff(expectedConditions, uc.Status.Conditions))
		})
	})

	Context("When waking a cluster up", func() {
		It("Should wake the cluster up", func() {
			By("By waking the UffizziCluster up")
			uc.Spec.Sleep = false
			Expect(k8sClient.Update(ctx, uc)).Should(Succeed())
		})

		It("Should be Awoken", func() {
			expectedConditions := td.ExpectedStatus.Awoken.Conditions
			uffizziClusterNSN := resources.CreateNamespacedName(uc.Name, ns.Name)
			By("Check if UffizziCluster has the correct Sleep conditions")
			Eventually(func() bool {
				if err := k8sClient.Get(ctx, uffizziClusterNSN, uc); err != nil {
					return false
				}
				return conditions.ContainsAllConditions(expectedConditions, uc.Status.Conditions)
			}, timeout, pollingTimeout).Should(BeTrue())

			//GinkgoWriter.Printf(conditions.CreateConditionsCmpDiff(expectedConditions, uc.Status.Conditions))
		})
	})
}
