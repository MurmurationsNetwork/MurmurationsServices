#!/bin/bash

NAMESPACES=("default" "murm-queue" "murm-logging")

# Declare associative array to track which pods need to be restarted
declare -A RESTART_PODS_PREFIX

restart_pods() {
    echo "Pausing for 1 minute before restarting pods..."
    sleep 60
    for namespace in "${NAMESPACES[@]}"; do
        for prefix in "${!RESTART_PODS_PREFIX[@]}"; do
            echo "Restarting pods with prefix $prefix in $namespace due to prior deletion..."
            kubectl get pods -n $namespace --no-headers | grep "$prefix" | awk '{print $1}' | while read pod_name; do
                echo "Deleting pod $pod_name in $namespace to trigger restart..."
                kubectl delete pod $pod_name -n $namespace
            done
        done
    done
}

# Check each namespace for specific conditions
for NAMESPACE in "${NAMESPACES[@]}"; do
    echo "Checking pods in namespace: $NAMESPACE"

    # Handle pods stuck in terminating state and record deletions
    kubectl get pods -n $NAMESPACE --no-headers | grep Terminating | awk '{print $1}' | while read pod_name; do
        echo "Deleting pod $pod_name stuck in terminating in $NAMESPACE..."
        kubectl delete pod $pod_name -n $NAMESPACE --force --grace-period=0
        # Check for murm-queue namespace to restart index-app and validation-app
        if [[ "$NAMESPACE" == "murm-queue" ]]; then
            RESTART_PODS_PREFIX["index-app"]=1
            RESTART_PODS_PREFIX["validation-app"]=1
        elif [[ "$NAMESPACE" == "default" ]]; then
            # Check prefix and record for later restart if it matches
            if [[ "$pod_name" == index-es-cluster* ]]; then
                RESTART_PODS_PREFIX["index-app"]=1
            elif [[ "$pod_name" == library-mongo* ]]; then
                RESTART_PODS_PREFIX["library-app"]=1
            elif [[ "$pod_name" == data-proxy-mongo* ]]; then
                RESTART_PODS_PREFIX["data-proxy-app"]=1
            fi
        fi
    done
done

restart_pods
