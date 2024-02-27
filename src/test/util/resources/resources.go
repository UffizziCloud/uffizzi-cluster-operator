package resources

import (
	"github.com/UffizziCloud/uffizzi-cluster-operator/src/api/v1alpha1"
	"github.com/UffizziCloud/uffizzi-cluster-operator/src/pkg/constants"
	"github.com/UffizziCloud/uffizzi-cluster-operator/src/test/util"
	"github.com/fluxcd/helm-controller/api/v2beta1"
	"github.com/fluxcd/source-controller/api/v1beta2"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"strings"
)

func CreateTestNamespace(name string) *v1.Namespace {
	return &v1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: strings.Join([]string{name, util.RandomString(5)}, "-"),
		},
	}
}

func CreateTestUffizziCluster(name, ns string) *v1alpha1.UffizziCluster {
	return &v1alpha1.UffizziCluster{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: ns,
		},
	}
}

func CreateTestUffizziClusterWithSpec(name, ns string, spec v1alpha1.UffizziClusterSpec) *v1alpha1.UffizziCluster {
	uc := CreateTestUffizziCluster(name, ns)
	uc.Spec = spec
	return uc
}

func GetHelmReleaseFromUffizziCluster(uc *v1alpha1.UffizziCluster) *v2beta1.HelmRelease {
	return &v2beta1.HelmRelease{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "uc-" + uc.Name,
			Namespace: uc.Namespace,
		},
	}
}

func GetETCDHelmReleaseFromUffizziCluster(uc *v1alpha1.UffizziCluster) *v2beta1.HelmRelease {
	return &v2beta1.HelmRelease{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "uc-etcd-" + uc.Name,
			Namespace: uc.Namespace,
		},
	}
}

func GetHelmRepositoryFromUffizziCluster(uc *v1alpha1.UffizziCluster) *v1beta2.HelmRepository {
	return &v1beta2.HelmRepository{
		ObjectMeta: metav1.ObjectMeta{
			Name:      constants.LOFT_HELM_REPO,
			Namespace: uc.Namespace,
		},
	}
}

func CreateNamespacedName(name, namespace string) types.NamespacedName {
	return types.NamespacedName{
		Name:      name,
		Namespace: namespace,
	}
}
