#!/bin/bash

# The goal of this script is to clean up pods that remain indefinitely in a
# terminating state due to node failures. It forcefully clears these pods.

NAMESPACES=("default" "murm-queue" "murm-logging")

# Iterate over each namespace to check for specific conditions
for NAMESPACE in "${NAMESPACES[@]}"; do
    echo "Checking pods in namespace: $NAMESPACE"

    # Identify and handle pods stuck in the terminating state.
    kubectl get pods -n $NAMESPACE --no-headers | grep Terminating | \
    awk '{print $1}' | while read pod_name; do
        echo "Deleting pod $pod_name stuck in terminating state in $NAMESPACE..."
        kubectl delete pod $pod_name -n $NAMESPACE --force --grace-period=0
    done
done
