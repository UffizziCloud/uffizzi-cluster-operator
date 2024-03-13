package e2e

import (
	"context"
	"k8s.io/client-go/rest"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
)

var (
	cfg     *rest.Config
	testEnv *envtest.Environment
	ctx     context.Context
	cancel  context.CancelFunc
	e2e     UffizziClusterE2E
)

type UffizziClusterE2E struct {
	IsTainted          bool
	UseExistingCluster bool
	K8SManager         ctrl.Manager
	K8SClient          client.Client
	Ctx                context.Context
	Cancel             context.CancelFunc
}

func GetE2E() *UffizziClusterE2E {
	return &e2e
}
