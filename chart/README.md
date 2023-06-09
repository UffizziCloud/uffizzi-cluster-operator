## Installation

If this is your first time using Helm, consult their documentation: https://helm.sh/docs/intro/quickstart/

Begin by adding our Helm repository:

```
helm repo add uffizzi https://uffizzicloud.github.io/uffizzi/
```

Then install the lastest version as a new release using the values you specified earlier. We recommend isolating Uffizzi in its own Namespace.

```
helm install my-uffizzi-cluster-operator uffizzi/uffizzi-cluster-operator --namespace uffizzi --create-namespace
```

If you encounter any errors here, tell us about them in [our Slack](https://join.slack.com/t/uffizzi/shared_invite/zt-ffr4o3x0-J~0yVT6qgFV~wmGm19Ux9A).

You should then see the release is installed:
```
helm list --namespace uffizzi
```
