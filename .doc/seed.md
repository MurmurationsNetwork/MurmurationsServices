# Seed data
1. List the pods and find the data-proxy-app. `kubectl get pods`
2. Connect to the data-proxy-app pod. `kubectl exec -it [POD_NAME] -- sh`
3. Change directory to seeder folder. `cd services/dataproxy/cmd/seeder`
4. Execute the command with the format `go run main.go "EXCEL_URL" "SCHEMA_NAME" "FROM" "TO"`
   The following is an example:
   ```
   go run main.go "https://docs.google.com/uc?export=download&id=[EXCEL_ID]" "karte_von_morgen-v1.0.0" 2 2
   ```
