#!/bin/bash

# This script determines which services should be excluded from actions like
# rebuild images or deployments based on the latest commit changes.

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

# Array to hold services with no changes detected.
exclude_services=()

# Function to check if changes are detected for a service.
function check_changes_for_service() {
    local service_paths=${paths_to_check[$1]}
    local changes_detected=false

    for file in $(git diff --name-only ${{ github.event.before }} ${{ github.event.after }}); do
        for path in $service_paths; do
            if echo "$file" | grep -q "^$path"; then
                changes_detected=true
                break 2
            fi
        done
    done

    echo $changes_detected
}

# Check each service and build the exclusion list.
for service in "${!paths_to_check[@]}"; do
    if [[ $(check_changes_for_service $service) == "false" ]]; then
        exclude_services+=("{\"service\": \"$service\"}")
    fi
done

# Convert the exclusion list to a JSON array and output it.
printf -v joined '%s,' "${exclude_services[@]}"
echo "[${joined%,}]"
