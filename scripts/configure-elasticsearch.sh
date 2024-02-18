#!/bin/bash

# This script configures Elasticsearch pods to access S3 storage by adding the
# necessary access and secret keys to each pod's Elasticsearch keystore.

read -p "Enter your ACCESS_KEY: " ACCESS_KEY
read -p "Enter your SECRET_KEY: " SECRET_KEY
read -p "Enter your NAMESPACE (default if empty): " NAMESPACE
NAMESPACE=${NAMESPACE:-default}

PODS="index-es-cluster-0 index-es-cluster-1 index-es-cluster-2"

for POD in $PODS; do
    echo "Configuring pod: $POD"

    kubectl exec -n ${NAMESPACE} $POD -- bash -c "echo ${ACCESS_KEY} | /usr/share/elasticsearch/bin/elasticsearch-keystore add --stdin s3.client.default.access_key"

    kubectl exec -n ${NAMESPACE} $POD -- bash -c "echo ${SECRET_KEY} | /usr/share/elasticsearch/bin/elasticsearch-keystore add --stdin s3.client.default.secret_key"

    echo "Configuration completed for pod: $POD"
done

echo "All pods configured."
