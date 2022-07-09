# Seed data
1. List the pods and find the data-proxy-app. `kubectl get pods`
2. Connect to the data-proxy-app pod. `kubectl exec -it [POD_NAME] -- sh`
3. Change directory to seeder folder. `cd services/dataproxy/cmd/seeder`
