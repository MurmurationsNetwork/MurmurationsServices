#!/bin/bash

# This script configures Elasticsearch pods for S3 access by updating
# the necessary access and secret keys in each pod's Elasticsearch keystore.

# Check if the correct number of arguments are passed
if [ "$#" -ne 3 ]; then
  echo "Usage: $0 ACCESS_KEY SECRET_KEY NAMESPACE"
  exit 1
fi

# Assign command line arguments to variables
ACCESS_KEY=$1
SECRET_KEY=$2
NAMESPACE=${3:-default} # Use provided namespace or default to 'default'

# Define Elasticsearch pods to configure.
PODS="index-es-cluster-0 index-es-cluster-1 index-es-cluster-2"

for POD in $PODS; do
    echo "Configuring pod: $POD"

    # Add/overwrite access key in the keystore.
    kubectl exec -n ${NAMESPACE} $POD -- bash -c "echo ${ACCESS_KEY} | \
      /usr/share/elasticsearch/bin/elasticsearch-keystore add -f \
      s3.client.default.access_key --stdin"

    # Add/overwrite secret key in the keystore.
    kubectl exec -n ${NAMESPACE} $POD -- bash -c "echo ${SECRET_KEY} | \
      /usr/share/elasticsearch/bin/elasticsearch-keystore add -f \
      s3.client.default.secret_key --stdin"

    echo "Configuration completed for pod: $POD"
done

echo "All pods configured."
