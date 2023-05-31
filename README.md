# Murmurations Services

_This project is licensed under the terms of the GNU General Public License v3.0_

[![GitHub Workflow Status](https://img.shields.io/github/actions/workflow/status/MurmurationsNetwork/MurmurationsServices/main.yaml?branch=main&style=flat-square)](https://github.com/MurmurationsNetwork/MurmurationsServices/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/MurmurationsNetwork/MurmurationsServices?style=flat-square)](https://goreportcard.com/report/github.com/MurmurationsNetwork/MurmurationsServices)

## Run locally

*Setup*

1. Install [Docker Desktop](https://www.docker.com/products/docker-desktop) and enable Kubernetes

2. Install [Helm](https://helm.sh/docs/intro/install/)

3. Install [NGINX Ingress Controller](docs/ingress-nginx)

4. Install [Skaffold](https://skaffold.dev/docs/install/)

5. Download large docker files

```
docker pull elasticsearch:7.17.5
docker pull kibana:7.17.5
```

6. [Create secrets](docs/secrets.md) for each service

7. Add the following to your host file `sudo vim /etc/hosts`

```
127.0.0.1   index.murmurations.dev
127.0.0.1   library.murmurations.dev
127.0.0.1   data-proxy.murmurations.dev
```

*After completing the setup*

1. Run `make dev` to start services

2. Try `index.murmurations.dev/v2/ping`, `library.murmurations.dev/v2/ping` and `data-proxy.murmurations.dev/v1/ping`

## Using Pre-commit and custom git hooks

1. Install pre-commit on your Mac by running `brew install pre-commit`.

2. Add pre-commit file and change permission.
   ```
   touch .git/hooks/pre-commit
   chmod +x .git/hooks/pre-commit
   ```
   
3. Use `vim .git/hooks/pre-commit` to edit the pre-commit file.
   ```
   #!/bin/sh
   
   PASS=true
   
   # Run Newman
   make newman-test
   if [[ $? != 0 ]]; then
       printf "\t\033[31mNewman\033[0m \033[0;30m\033[41mFAILURE!\033[0m\n"
       PASS=false
   else
       printf "\t\033[32mNewman\033[0m \033[0;30m\033[42mpass\033[0m\n"
   fi
   
   if ! $PASS; then
       printf "\033[0;30m\033[41mCOMMIT FAILED\033[0m\n"
       exit 1
   else
       printf "\033[0;30m\033[42mCOMMIT SUCCEEDED\033[0m\n"
   fi
   
   exit 0
   ```

4. Run `pre-commit install` to set up the git hook scripts in your local repository. You can safely ignore the "`Running in migration mode with existing hooks`" message.

Now, pre-commit will run automatically on `git commit`. If you want to manually run all pre-commit hooks on a repository, run `pre-commit run --all-files`.

## Running E2E Tests Locally

1. Ensure that you have Newman installed. If not, install it using the following command: `npm install -g newman`.
2. Execute the command `make dev` to set up the servers.
3. Run the command `make newman-test` to initiate the end-to-end (E2E) tests.


## Run in DigitalOcean

1. Install [Helm](https://helm.sh/docs/intro/install/) and [doctl](https://github.com/digitalocean/doctl#installing-doctl)

2. Create Kubernetes Clusters in DigitalOcean and install [metrics-server](https://github.com/kubernetes-sigs/metrics-server#installation) for monitoring CPU/MEM

3. Install [NGINX Ingress Controller](docs/ingress-nginx)

4. [Create secrets](docs/secrets.md) for each service

5. [Create PVCs](docs/pvcs.md) for each service

6. [Deploy services](docs/deploy-services.md)

7. [Installing and Configuring Cert-Manager](docs/cert-manager.md)

8. Try `index.murmurations.network/v2/ping` or `library.murmurations.network/v1/ping`

**Optional**

- [Logging, Monitoring and Alerting](docs/logging-monitoring-alerting/README.md)

- [Seed the data](docs/seed.md)
