# Uffizzi Cluster Operator

A Kubernetes operator for creating fully managed virtual clusters. 

- [x] Create a VCluster.
- [x] Install Helm and Source Controller in the VCluster.
- [x] Create Helm charts when mentioned in the UffizziCluster CRD.
- [x] Expose Ingress for the VCluster to connect via the `vcluster connect <name> --server=<vcluster-ingress>` command.
- [x] Expose Services from within the VCluster.
- [x] Enable authentication for the Ingresses.
- [ ] Expose Ingress which gives commandline access to the VCluster environment (run vcluster connect in a terminal and give webterminal access ?)
- [ ] Suspend the VCluster if it is not being used for a certain period of time.

## Requirements

### Flux Components
We need the Helm Controller and Source Controller which are components of flux cd. Use the following command to install 
them (if you are using GKE).
```bash
flux install --namespace=flux-system --components="source-controller,helm-controller" --toleration-keys="sandbox.gke.io/runtime"
```

### Nginx Ingress Controller

The NGINX Ingress controller also need to be installed in the said cluster with SSL Passthrough enabled (disabled by default) https://kubernetes.github.io/ingress-nginx/user-guide/tls/#ssl-passthrough.

### Keycloak

Install Keycloak with the Ingress.

```bash
kubectl create -f config/helmreleases/dex.yaml
```

Then configure it based on Then configure it based on https://oauth2-proxy.github.io/oauth2-proxy/docs/configuration/oauth_provider/#keycloak-oidc-auth-provider


### OAuth2 Proxy

```bash
kubectl create -f config/helmreleases/oauth2proxy.yaml
```

## Development

To install the CRD for UffizziCluster and run the operator locally, use the following command:

```bash
make install && make run
```

You also need to have the fluxcd binary installed on your system which the operator uses to get the yaml
to inject in the HelmRelease values.

## Usage

To create a sample UffizziCluster, use the following command:

```bash
kubectl apply -f examples/basic-ucluster.yml
```

The VCluster will be created with the Helm and Source Controllers installed as well.

## Cleanup

```bash
kubectl delete UffizziCluster,helmrelease,helmrepository --all && make uninstall
```
