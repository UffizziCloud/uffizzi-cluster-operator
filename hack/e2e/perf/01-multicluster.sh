#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

LOOP_BOUND="${1:-3}"

for i in $(seq 1 $LOOP_BOUND); do
    # Generate a unique namespace name
    NAMESPACE="uffizzi-cluster-$i-$(date +%s)"
    # Create the namespace
    kubectl create namespace "$NAMESPACE" > /dev/null
    # Label the namespace
    kubectl label namespace "$NAMESPACE" app=uffizzi > /dev/null
    # Deploy the UffizziCluster resource to the unique namespace
    kubectl create -f hack/e2e/perf/manifests/gen.yaml --namespace="$NAMESPACE" > /dev/null
done

namespaces=($(kubectl get ns --selector='app=uffizzi' -o jsonpath='{.items[*].metadata.name}'))

# Function to check the APIReady condition of a UffizziCluster within a namespace
check_api_ready() {
    local namespace=$1
    local api_ready=$(kubectl get uffizzicluster --namespace="$namespace" -o jsonpath='{.items[0].status.conditions[?(@.type=="APIReady")].status}')
    echo "$api_ready"
}

# Monitor each UffizziCluster for readiness
start_time=$(date +%s)
for ns in "${namespaces[@]}"; do
#    echo "Monitoring UffizziCluster in namespace $ns"
    while true; do
        api_ready=$(check_api_ready "$ns")
        if [ "$api_ready" == "True" ]; then
#            echo "UffizziCluster in namespace $ns is ready."
            break
        else
#            echo "Waiting for UffizziCluster in namespace $ns to become ready..."
            sleep 5
        fi
    done
done
end_time=$(date +%s)

# Calculate the total time taken for all UffizziClusters to become ready
total_time=$((end_time - start_time))
echo "$total_time"

# Cleanup
for ns in "${namespaces[@]}"; do
    kubectl delete namespace "$ns" > /dev/null
done
