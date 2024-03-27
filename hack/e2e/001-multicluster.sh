#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

# create multiple uffizzi clusters

for i in {1..10}; do
  kubectl create -f hack/e2e/manifests/001-multicluster.yaml
done

# Retrieve the names of the newly created UffizziCluster resources
# Assuming 'uffizzicluster' is the kind for the UffizziCluster resources
uffizzi_clusters=($(kubectl get uffizzicluster -o jsonpath='{.items[*].metadata.name}'))

# Check if we have any clusters to monitor
if [ ${#uffizzi_clusters[@]} -eq 0 ]; then
    echo "No UffizziClusters found. Exiting..."
    exit 1
fi

echo "Monitoring the following UffizziClusters for readiness: ${uffizzi_clusters[@]}"

# Function to check the APIReady condition of a UffizziCluster
check_api_ready() {
    local cluster_name=$1
    api_ready=$(kubectl get uffizzicluster "$cluster_name" -o jsonpath='{.status.conditions[?(@.type=="APIReady")].status}')
    echo "$api_ready"
}

# Monitor each UffizziCluster until the APIReady condition is True
start_time=$(date +%s)
for cluster in "${uffizzi_clusters[@]}"; do
    while true; do
        api_ready=$(check_api_ready "$cluster")
        if [ "$api_ready" == "True" ]; then
            echo "UffizziCluster $cluster is ready."
            break
        else
            echo "Waiting for UffizziCluster $cluster to become ready..."
            sleep 5
        fi
    done
done
end_time=$(date +%s)

# Calculate the total time taken for all UffizziClusters to become ready
total_time=$((end_time - start_time))
echo "All UffizziClusters are ready. Total time taken: $total_time seconds."
