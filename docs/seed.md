# Data Proxy Seeding

## Create mappings

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

## Seed data in production
1. List the pods and find the data-proxy-app. `kubectl get pods`
2. Connect to the data-proxy-app pod. `kubectl exec -it [POD_NAME] -- sh`
3. Change directory to seeder folder. `cd app`
4. Execute the command with the format `./seeder "EXCEL_URL" "SCHEMA_NAME" "FROM" "TO"`
   The following is an example:
   ```
   ./seeder "https://docs.google.com/uc?export=download&id=[EXCEL_ID]" "karte_von_morgen-v1.0.0" 2 2
   ```

## Seed data in local
1. List the pods and find the data-proxy-app. `kubectl get pods`
2. Connect to the data-proxy-app pod. `kubectl exec -it [POD_NAME] -- sh`
3. Change directory to seeder folder. `cd cmd/dataproxy/seeder`
4. Execute the command with the format `go run main.go "EXCEL_URL" "SCHEMA_NAME" "FROM" "TO"`
   The following is an example:
   ```
   go run main.go "https://docs.google.com/uc?export=download&id=[EXCEL_ID]" "karte_von_morgen-v1.0.0" 2 2
   ```
