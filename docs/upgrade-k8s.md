# Upgrade k8s from 1.21.5 - 1.22.8
1. Upgrade to 1.21.11 - 30 mins
2. Upgrade to 1.22.8 - 30 mins
3. Run command.
   ```
   helm upgrade \
   cert-manager jetstack/cert-manager \
   --namespace cert-manager \
   --create-namespace \
   --version v1.8.0 \
   --set installCRDs=true
   ```
4. Delete all dead pods - two cert-manager pods and one prometheus pod (<ctrl+d> in k9s)
5. Update Makefile's env and SPECIFIC_TAG(check [tags](https://hub.docker.com/r/murmurations/index/tags)).
6. Redeploy all services.
   ```
   make manually-deploy-mq
   make manually-deploy-index
   make manually-deploy-validation
   make manually-deploy-library
   make manually-deploy-nodecleaner
   make manually-deploy-schemaparser
   make manually-deploy-revalidatenode
   ```
7. Run Postman tests to make sure all services are running as expected.
