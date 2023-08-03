# Uffizzi Cluster Operator

A Kubernetes operator for creating fully managed virtual clusters. The `UffizziCluster` resource is a crucial component in the architecture of the Uffizzi application. It serves as a representation of a cluster that can be created and managed within Uffizzi. 

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

To install the CRD for UffizziCluster and run the operator locally, use the following command:

```bash
make install && make run
```

You also need to have the `fluxcd` binary installed on your system which the operator uses to get the yaml
to inject in the HelmRelease values.

## Usage

To create a sample UffizziCluster, use the following command:

```bash
kubectl apply -f examples/helm-basic.yml
```

The VCluster will be created with the Helm and Source Controllers installed as well.

## Cleanup

```bash
kubectl delete UffizziCluster,helmrelease,helmrepository --all && make uninstall
```

## Architecture Overview

The general architecture of the Uffizzi system can be summarized as follows:

`
CLI -> App -> Controller -> (UffizziCluster) -> UffizziClusterOperator -> HelmReleaseCRD -> HelmController
`

Here's a breakdown of the components and their interactions:

1. **CLI**: The Command-Line Interface (CLI) is used by users to interact with the Uffizzi application. It calls the Uffizzi app API to create an Uffizzi cluster.

2. **App**: The Uffizzi app receives requests from the CLI and communicates with the Uffizzi controller to create a new Uffizzi cluster.

3. **Controller**: The Uffizzi controller manages the lifecycle of Uffizzi clusters. When notified by the app about a new cluster, it creates an `UffizziCluster` resource and monitors its status for updates.

4. **UffizziCluster**: The `UffizziCluster` resource represents an individual cluster within the Uffizzi system. The controller creates it and contains information about the desired state and status of the cluster.

5. **UffizziClusterOperator**: The UffizziClusterOperator component reconciles on the `UffizziCluster` resource. It interacts with the HelmReleaseCRD to create HelmReleases with vCluster definitions.

6. **HelmReleaseCRD**: The HelmReleaseCRD (Custom Resource Definition) defines the structure of the HelmRelease resource. It is used to create HelmReleases that define the desired state of the vCluster resources.

7. **HelmController**: The HelmController continuously monitors the HelmRelease resources. Once a HelmRelease for a vCluster is created, the HelmController creates all the resources specified in the Helm chart. These resources have an owner reference set to the HelmRelease, establishing their relationship.

8. **vCluster**: The vCluster resources are created by the HelmController and owned by the corresponding HelmRelease CR. They represent the virtual clusters managed by Uffizzi.

## UffizziCluster Creation and Management Flow

The following steps outline the flow of operations involved in creating and managing an UffizziCluster:

1. The user interacts with the Uffizzi application through the CLI, requesting the creation of a new Uffizzi cluster.

2. The Uffizzi app receives the request and communicates with the Uffizzi controller, notifying it of the new cluster creation request.

3. The Uffizzi controller creates an `UffizziCluster` resource, which represents the desired state and status of the cluster.

4. The UffizziClusterOperator component reconciles on the `UffizziCluster` resource and interacts with the HelmReleaseCRD.

5. The UffizziClusterOperator creates a HelmRelease, defining the vCluster resources required for the Uffizzi cluster.

6. The HelmController continuously monitors the HelmRelease resources.

7. Once the HelmController detects the creation of a HelmRelease for a vCluster, it creates all the resources specified in the Helm chart.

8. All the resources created by the HelmController have an owner reference set to the corresponding HelmRelease, establishing the ownership relationship.

9. The UffizziClusterController monitors the status of the UffizziCluster resource and updates it when the cluster is ready and has an assigned ingress address.
