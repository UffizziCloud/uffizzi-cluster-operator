kubectl get pods --all-namespaces -o jsonpath="{range .items[*]}{.metadata.namespace}{' '}{.metadata.name}{'\n'}{end}" | while read -r line; do
  namespace=$(echo "$line" | cut -d' ' -f1)
  pod=$(echo "$line" | cut -d' ' -f2)
  kubectl patch pod "$pod" -n "$namespace" --type='json' -p='[{"op": "add", "path": "/spec/tolerations", "value": [{"key": "sandbox.gke.io/runtime", "operator": "Equal", "value": "gvisor", "effect": "NoSchedule"}]}]' || true
done