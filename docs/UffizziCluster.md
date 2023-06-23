# UffizziCluster Resource

The `UffizziCluster` resource is a key component in the architecture of the Uffizzi application. It serves as a representation of a cluster that can be created and managed within Uffizzi. 

## Architecture Overview

The general architecture of the Uffizzi system can be summarized as follows:

`
CLI -> App -> Controller -> (UffizziCluster) -> UffizziClusterOperator -> HelmReleaseCRD -> HelmController
`

Here's a breakdown of the components and their interactions:

1. **CLI**: The Command-Line Interface (CLI) is used by users to interact with the Uffizzi application. It calls the Uffizzi app API to create an Uffizzi cluster.

2. **App**: The Uffizzi app receives requests from the CLI and communicates with the Uffizzi controller to initiate the creation of a new Uffizzi cluster.

3. **Controller**: The Uffizzi controller is responsible for managing the lifecycle of Uffizzi clusters. When notified by the app about a new cluster, it creates an `UffizziCluster` resource and monitors its status for updates.

4. **UffizziCluster**: The `UffizziCluster` resource represents an individual cluster within the Uffizzi system. It is created by the controller and contains information about the desired state and status of the cluster.

5. **UffizziClusterOperator**: The UffizziClusterOperator component reconciles on the `UffizziCluster` resource. It interacts with the HelmReleaseCRD to create HelmReleases with vCluster definitions.

6. **HelmReleaseCRD**: The HelmReleaseCRD (Custom Resource Definition) defines the structure of the HelmRelease resource. It is used to create HelmReleases that define the desired state of the vCluster resources.

7. **HelmController**: The HelmController continuously monitors the HelmRelease resources. Once a HelmRelease for a vCluster is created, the HelmController creates all the resources specified in the Helm chart. These resources have an owner reference set to the HelmRelease, establishing the relationship between them.

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

Through this flow, the UffizziCluster resource serves as a central point of control and coordination for the creation and management of Uffizzi clusters within the Uffizzi system.

For more detailed information about the Uffizzi system and its components, please refer to the relevant documentation and resources provided.
