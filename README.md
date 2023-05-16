# Uffizzi Cluster Operator

A Kubernetes operator for creating ephemeral clusters. 

- [x] Create a VCluster.
- [x] Install Helm and Source Controller in the VCluster.
- [x] Create Helm charts when mentioned in the UffizziCluster CRD.
- [x] Expose Ingress for the VCluster.
- [ ] Expose Ingreess which point to applications in the VCluster.
- [ ] Expose Ingress which gives commandline access to the VCluster environment (run vcluster connect in a terminal and give webterminal access ?)
- [ ] Suspend the VCluster if it is not being used for a certain period of time.

## Requirements
We need the Helm Controller and Source Controller which are components of flux cd. Use the following command to install 
them (if you are using GKE).
```azure
flux install --namespace=flux-system --components="source-controller,helm-controller" --toleration-keys="sandbox.gke.io/runtime"
```

The NGINX Ingress controller also need to be installed in the said cluster with SSL Passthrough enabled (disabled by default) https://kubernetes.github.io/ingress-nginx/user-guide/tls/#ssl-passthrough.

## Development

To install the CRD for UffizziCluster and run the operator locally, use the following command:

```azure
make install && make run
```

You also need to have the fluxcd binary installed on your system which the operator uses to get the yaml
to inject in the HelmRelease values.

## Usage

To create a sample UffizziCluster, use the following command:

```
kubectl apply -f examples/ucluster.yml
```

The VCluster will be created with the Helm and Source Controllers installed as well.

## Cleanup

```
kubectl delete UffizziCluster,helmrelease,helmrepository --all && make uninstall
```
