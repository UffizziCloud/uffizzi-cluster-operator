kubectl get deployments --all-namespaces -o jsonpath="{range .items[*]}{.metadata.namespace}{' '}{.metadata.name}{'\n'}{end}" | while read -r line; do
  namespace=$(echo "$line" | cut -d' ' -f1)
  deployment=$(echo "$line" | cut -d' ' -f2)
  kubectl patch deployment "$deployment" --patch-file="./.github/kubectl-patch/nodeselector-toleration.yaml" -n "$namespace"
done

kubectl get statefulset --all-namespaces -o jsonpath="{range .items[*]}{.metadata.namespace}{' '}{.metadata.name}{'\n'}{end}" | while read -r line; do
  namespace=$(echo "$line" | cut -d' ' -f1)
  statefulset=$(echo "$line" | cut -d' ' -f2)
  kubectl patch statefulset "$statefulset" --patch-file="./.github/kubectl-patch/nodeselector-toleration.yaml" -n "$namespace"
done
