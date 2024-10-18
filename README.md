# Murmurations Services

> _This project is licensed under the terms of the GNU General Public License v3.0_

[![GitHub Workflow Status](https://img.shields.io/github/actions/workflow/status/MurmurationsNetwork/MurmurationsServices/main.yaml?branch=main&style=flat-square)](https://github.com/MurmurationsNetwork/MurmurationsServices/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/MurmurationsNetwork/MurmurationsServices?style=flat-square)](https://goreportcard.com/report/github.com/MurmurationsNetwork/MurmurationsServices)

## Run in Production

We are utilizing [Rancher](https://www.rancher.com/) to orchestrate the deployment of Murmurations services in Kubernetes clusters. For detailed instructions on setting up a Kubernetes cluster with Rancher and guidance on launching the index, library and other services that implement the Murmurations protocol, please refer to the [Rancher-managed Kubernetes documentation](docs/rancher/README.md).

## Troubleshooting

For troubleshooting, please refer to [Debugging Guide](./docs/debugging-guide/README.md).

## Run Locally

### Setup

1. Install [Docker Desktop](https://www.docker.com/products/docker-desktop) and enable Kubernetes

2. Install [Helm](https://helm.sh/docs/intro/install/)

3. Install [NGINX Ingress Controller](docs/ingress-nginx)

4. Install [Skaffold](https://skaffold.dev/docs/install/)

5. Download large docker files

    ```sh
    docker pull elasticsearch:7.17.5
    docker pull kibana:7.17.5
    ```

6. [Create secrets](docs/rancher/07-run-murmuration-services/secrets.md) for each service

7. Add the following to your host file `sudo vim /etc/hosts`

    ```sh
    127.0.0.1   index.murmurations.developers
    127.0.0.1   library.murmurations.developers
    127.0.0.1   data-proxy.murmurations.developers
    ```

### After Completing the Setup

1. Run `make dev` to start services

2. Try `index.murmurations.developers/v2/ping`, `library.murmurations.developers/v2/ping` and `data-proxy.murmurations.developers/v1/ping`

## Setting Up Pre-commit and Custom Git Hooks for Development

> Note: [Pre-commit](https://pre-commit.com) is a linter to ensure consistent style, etc. Please use it before submitting pull requests to this repository.

1. Install pre-commit on a Mac by running `brew install pre-commit`.

2. Add pre-commit file and change permission.

    ```sh
    touch .git/hooks/pre-commit
    chmod +x .git/hooks/pre-commit
    ```

3. Use `vim .git/hooks/pre-commit` to edit the pre-commit file.

   ```sh
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

## Testing

- [Testing Murmurations Services](docs/testing/README.md)

## Optional

- [Seed the data](docs/seed.md)
