#!/bin/bash

# This script determines which services should be excluded from actions like
# rebuilding images or deployments based on the changes in the latest commit.

# Arguments passed from the workflow: previous and current commit SHAs.
BEFORE_COMMIT=$1
AFTER_COMMIT=$2

# Associative array mapping services to their directory patterns.
# Each service is associated with directories and files it depends on.
declare -A SERVICE_PATHS=(
    ["index"]="go.mod pkg/ cmd/index/ services/index/"
    ["library"]="go.mod pkg/ cmd/library/ services/library/"
    ["validation"]="go.mod pkg/ cmd/validation/ services/validation/"
    ["dataproxy"]="go.mod pkg/ cmd/dataproxy/ services/dataproxy/"
    ["nodecleaner"]="go.mod pkg/ cmd/nodecleaner/ services/nodecleaner/"
    ["revalidatenode"]="go.mod pkg/ cmd/revalidatenode/ services/revalidatenode/"
    ["schemaparser"]="go.mod pkg/ cmd/schemaparser/ services/schemaparser/"
    ["dataproxyupdater"]="go.mod pkg/ cmd/dataproxyupdater/ services/dataproxyupdater/"
    ["dataproxyrefresher"]="go.mod pkg/ cmd/dataproxyrefresher/ services/dataproxyrefresher/"
    ["maintenance"]="build/maintenance"
)

# Function to check if any files changed for a service.
# It iterates over changed files and checks if they match any service path.
check_changes_for_service() {
    local service=$1
    local service_paths=${SERVICE_PATHS[$service]}
    local changes_detected=false

    # Get the list of changed files between two commits.
    local changed_files=$(git diff --name-only "$BEFORE_COMMIT" "$AFTER_COMMIT")

    for file in $changed_files; do
        for path in $service_paths; do
            if [[ "$file" == $path* ]]; then
                changes_detected=true
                break 2
            fi
        done
    done

    echo $changes_detected
}

# Main function to determine services to exclude based on changes.
main() {
    local exclude_services=()

    for service in "${!SERVICE_PATHS[@]}"; do
        if [[ $(check_changes_for_service "$service") == "false" ]]; then
            exclude_services+=("{\"service\": \"$service\"}")
        fi
    done

    # Convert the exclusion list to a JSON array and output it.
    printf -v joined '%s,' "${exclude_services[@]}"
    echo "[${joined%,}]"
}

# Execute the main function.
main
