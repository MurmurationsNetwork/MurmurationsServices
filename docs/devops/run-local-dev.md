# Running a Local Development Environment

## Setup

1. Install [Docker Desktop](https://www.docker.com/products/docker-desktop) and enable Kubernetes

2. Install [Helm](https://helm.sh/docs/intro/install/)

3. Install [NGINX Ingress Controller](/docs/devops/install-ingress-nginx.md)

4. Install [Skaffold](https://skaffold.dev/docs/install/)

5. Download large docker files

    ```sh
    docker pull elasticsearch:7.17.27
    docker pull kibana:7.17.27
    ```

6. [Create secrets](/docs/rancher/07-run-murmuration-services/secrets.md) for each service

7. Add the following to your host file `sudo vim /etc/hosts`

    ```sh
    127.0.0.1   index.murmurations.developers
    127.0.0.1   library.murmurations.developers
    127.0.0.1   data-proxy.murmurations.developers
    ```

## After Completing the Setup

1. Run `make dev` to start services

2. Try `GET`ting the following API endpoints to check the services are running:

- <http://index.murmurations.developers/v2/ping>
- <http://library.murmurations.developers/v2/ping>
- <http://data-proxy.murmurations.developers/v1/ping>

## Setting Up Pre-commit and Custom Git Hooks for Development (Optional)

> Note: [Pre-commit](https://pre-commit.com) is a linter to ensure consistent style, etc. **Please use it before submitting pull requests to this repository.** There is no need to install if you are not planning to submit pull requests.

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
