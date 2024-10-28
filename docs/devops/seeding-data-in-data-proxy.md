# Seeding Data in the Data Proxy

This guide provides a step-by-step process for mapping data fields and seeding data in the data proxy.

This process was built for the [Karte von Morgen](https://www.kartevonmorgen.org/) dataset, but it can be adapted for other datasets that are stored in a CSV file and accessed via a URL.

## Create Mappings

The purpose of the mappings is to map the data fields in the CSV file to the data fields defined in the schema used for posting the nodes to the data proxy.

1. Post mappings to `{{dataProxyBaseUrl}}/mappings`

Example request body:

```json
{
    "schema": "karte_von_morgen-v1.0.0",
    "name": "title",
    "description": "description",
    "latitude": "lat",
    "longitude": "lng",
    "primary_url": "homepage",
    "tags": "tags",
    "image": "image_url",
    "kvm_category": "categories",
    "email": "contact_email",
    "region": "city",
    "country_name": "country",
    "locality": "state"
}
```

## Seed Data in a Production Environment

1. List the pods and find the data-proxy-app. `kubectl get pods`
2. Connect to the data-proxy-app pod. `kubectl exec -it [POD_NAME] -- sh`
3. Change directory to seeder folder. `cd app`
4. Execute the command with the format `./seeder "EXCEL_URL" "SCHEMA_NAME" "FROM" "TO"` (`FROM` and `TO` are the page numbers in the CSV file that you want to seed).
   The following is an example:

   ```bash
   ./seeder "https://docs.google.com/uc?export=download&id=[EXCEL_ID]" "karte_von_morgen-v1.0.0" 2 101
   ```

## Seed Data in a Local Development Environment

1. List the pods and find the data-proxy-app. `kubectl get pods`
2. Connect to the data-proxy-app pod. `kubectl exec -it [POD_NAME] -- sh`
3. Change directory to seeder folder. `cd cmd/dataproxy/seeder`
4. Execute the command with the format `go run main.go "EXCEL_URL" "SCHEMA_NAME" "FROM" "TO"`
   The following is an example:

   ```bash
   go run main.go "https://docs.google.com/uc?export=download&id=[EXCEL_ID]" "karte_von_morgen-v1.0.0" 2 101
   ```
