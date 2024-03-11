#!/bin/bash

# Check if the correct number of arguments are passed
if [ "$#" -lt 2 ]; then
  echo "Usage: $0 ACCESS_KEY SECRET_KEY"
  exit 1
fi

# Assign the command line arguments to variables
ACCESS_KEY=$1
SECRET_KEY=$2

# Define Elasticsearch pods to configure
PODS="index-es-cluster-0 index-es-cluster-1 index-es-cluster-2"

for POD in $PODS; do
    echo "Configuring pod: $POD"

    # Add/overwrite access key in the keystore
    kubectl exec -n default $POD -- bash -c "echo ${ACCESS_KEY} | \
      /usr/share/elasticsearch/bin/elasticsearch-keystore add -f \
      s3.client.default.access_key --stdin"

    # Add/overwrite secret key in the keystore
    kubectl exec -n default $POD -- bash -c "echo ${SECRET_KEY} | \
      /usr/share/elasticsearch/bin/elasticsearch-keystore add -f \
      s3.client.default.secret_key --stdin"

    echo "Configuration completed for pod: $POD"
done

echo "All pods configured."
