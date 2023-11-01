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

## Installation

This operator is best installed with Helm as a dependency of the Uffizzi Helm chart: https://github.com/UffizziCloud/uffizzi/tree/develop/charts/uffizzi-app

Alternatively, if you're installing the Uffizzi control plane on a separate cluster, you may install this operator as a dependency of the Uffizzi controller Helm chart: https://github.com/UffizziCloud/uffizzi_controller/tree/uffizzi-controller-2.2.5/charts/uffizzi-controller

Lastly, a Helm chart for the operator itself is provided: https://github.com/UffizziCloud/uffizzi-cluster-operator/tree/main/chart

## Dependencies

### Flux Components
We need the Helm Controller and Source Controller which are components of flux cd. Use the following command to install 
them (if you are using GKE).
```bash
flux install --namespace=flux-system --components="source-controller,helm-controller" --toleration-keys="sandbox.gke.io/runtime"
```

### Nginx Ingress Controller

The NGINX Ingress controller also need to be installed in the said cluster with SSL Passthrough enabled (disabled by default) https://kubernetes.github.io/ingress-nginx/user-guide/tls/#ssl-passthrough.

<!-- ### Keycloak

Install Keycloak with the Ingress. -->

<!-- ```bash
kubectl create -f config/helmreleases/dex.yaml
``` -->

Then configure it based on https://oauth2-proxy.github.io/oauth2-proxy/docs/configuration/oauth_provider/#keycloak-oidc-auth-provider


<!-- ### OAuth2 Proxy

```bash
kubectl create -f config/helmreleases/oauth2proxy.yaml
``` -->

## Development

We're most often using our shared `qa` cluster for development. Before "taking over", shut down the operator that's active there. If that was installed via the Helm chart, this command will look like this:

```bash
kubectl scale --replicas=0 deployment uffizzi-uffizzi-cluster-operator --namespace uffizzi
```

Be sure so scale back up after you're finished testing!
```bash
kubectl scale --replicas=1 deployment uffizzi-uffizzi-cluster-operator --namespace uffizzi
```

To install the CRD for UffizziCluster and run the operator locally, use the following command:

```bash
make install && make run
```

## Usage

Once installed, use the following command to create a sample UffizziCluster:

```bash
kubectl apply -f examples/helm-basic.yml
```

The VCluster will be created with the Helm and Source Controllers installed as well.

## Cleanup

```bash
kubectl delete UffizziCluster,helmrelease,helmrepository --all && make uninstall
```
