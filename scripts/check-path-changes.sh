#!/bin/bash

# This script determines if the latest commit includes changes in directories
# relevant to a specified service. It's used in CI/CD pipelines to trigger
# actions like tests or deployments for the affected service.

# First argument is the service name.
service=$1

# Associative array mapping services to their respective path patterns.
declare -A paths_to_check=(
    ["index"]="cmd/index/ services/index/ pkg/"
    ["library"]="cmd/library/ services/library/ pkg/"
    ["geoip"]="cmd/geoip/ services/geoip/ pkg/"
    ["validation"]="cmd/validation/ services/validation/ pkg/"
    ["dataproxy"]="services/dataproxy/ pkg/"
    ["nodecleaner"]="cmd/nodecleaner/ services/cronjob/nodecleaner/ pkg/"
    ["revalidatenode"]="cmd/revalidatenode/ services/cronjob/revalidatenode/ pkg/"
    ["schemaparser"]="cmd/schemaparser/ services/cronjob/schemaparser/ pkg/"
    ["dataproxyupdater"]="services/cronjob/dataproxyupdater/ pkg/"
    ["dataproxyrefresher"]="cmd/dataproxyrefresher/ services/cronjob/dataproxyrefresher/ pkg/"
)

# Get paths for the specified service.
service_paths=${paths_to_check[$service]}

# Flag to detect changes.
changes_detected=false

# Loop through files changed in the last commit.
for file in $(git diff --name-only ${{ github.event.before }} ${{ github.event.after }}); do
    for path in $service_paths; do
        if echo "$file" | grep -q "^$path"; then
            changes_detected=true
            break 2
        fi
    done
done

# Set the output variable.
echo "::set-output name=${service}_changes_detected::$changes_detected"
