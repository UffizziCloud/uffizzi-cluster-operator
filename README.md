# Ephemeral Cluster Operator

A Kubernetes operator for creating ephemeral clusters. 

- [x] Create a VCluster
- [x] Install Helm and Source Controller in the VCluster
- [ ] Helm Chart in the VCluster by specifying it in the EphemeralCluster CRD

## Requirements
We need the Helm Controller and Source Controller which are components of flux cd. Use the following command to install 
them (if you are using GKE).
```azure
flux install --namespace=flux-system --components="source-controller,helm-controller" --toleration-keys="sandbox.gke.io/runtime"
```

## Development

To install the CRD for EphemeralCluster and run the operator locally, use the following command:

```azure
make install && make run
```

You also need to have the fluxcd binary installed on your system which the operator uses to get the yaml
to inject in the HelmRelease values.

## Usage

To create a sample EphemeralCluster, use the following command:

```
kubectl apply -f examples/ecluster.yml
```

The VCluster will be created with the Helm and Source Controllers installed as well.

## Cleanup

```
kubectl delete ephemeralcluster,helmrelease,helmrepository --all && make uninstall
```