#!/bin/bash

# This script configures Elasticsearch pods for S3 access by updating
# the necessary access and secret keys in each pod's Elasticsearch keystore.

# Ensure we read input from the terminal even when input is redirected
if [ -t 0 ]; then
  # Terminal input available, proceed as normal
  echo "Reading input from the terminal..."
else
  # Input redirected, force read from /dev/tty (the terminal)
  exec < /dev/tty
fi

# Prompt for S3 access credentials and namespace.
read -p "Enter your ACCESS_KEY: " ACCESS_KEY
read -p "Enter your SECRET_KEY: " SECRET_KEY
read -p "Enter your NAMESPACE (default if empty): " NAMESPACE
NAMESPACE=${NAMESPACE:-default} # Default namespace if not specified.

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
