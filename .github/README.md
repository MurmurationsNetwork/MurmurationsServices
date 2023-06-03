# `/.github`

## GitHub Workflows

GitHub Workflows is a feature of GitHub Actions that automates tasks and
workflows in a repository. A workflow is defined in a YAML file in the
[.github/workflows directory](/.github/workflows) of a repository and can be
triggered by events such as pushes, pull requests or scheduled intervals.
Workflows run on GitHub-hosted virtual machines or on self-hosted runners.
MurmurationsServices uses Workflows to build, test and deploy code.

### `pull_request.yaml`

This workflow defines two jobs that run tests and linting for the project when a
pull request is opened or updated on the main branch. Automatically running
tests and linting on every pull request helps ensure that the code is high
quality and free of errors.

### `main.yaml`

This workflow is triggered by a push event on the `main` branch and automates
the build, test and deployment process for the project's components. The
workflow uses `make` commands to build and publish Docker images for each
component and deploys the components to a Kubernetes cluster on DigitalOcean.

The `test` job executes unit tests for the project and consists of three steps.
The first step uses the `actions/checkout` action to check out the code for the
project. The second step uses the `actions/setup-go` action to set up the Go
environment with version 1.19.5. The third step runs the `make test` command,
which is defined in the project's `Makefile` and runs the tests for the project.

This `build-*` jobs build and publish Docker images for different components of
the project. Each job runs on an `ubuntu-latest` runner and has three steps. The
first step uses the `actions/checkout` action to check out the code for the
project. The second step uses the `docker/login-action` action to log in to
DockerHub using the `secrets.DOCKERHUB_USERNAME` and `secrets.DOCKERHUB_TOKEN`
secrets. The third step runs a `make` command to build and publish the Docker
image for the component.

The `deploy` job automates the deployment process for the project's components.
The job is triggered by the completion of several other jobs, which are
specified using the `needs` key. The job runs on an `ubuntu-latest` runner and
deploys various components of the project to a Kubernetes cluster on
DigitalOcean.

The job consists of several steps that use the `make` command to restart the
deployments for each component. The first step uses the `actions/checkout`
action to check out the code for the project. The second step installs `doctl`,
a command-line tool for managing DigitalOcean resources, and saves the
`kubeconfig` for the Kubernetes cluster with short-lived credentials. The
remaining steps restart the deployments for each component using the make
`deploy-<component>` command.

The `e2e-test` job is disabled because the end-to-end tests reference test
schemas not contained in the production library.

### `test.yaml`

This workflow is exactly the same as the one above for the `test`, `build-*` and
`deploy` jobs except that the deployment is made to a test environment, not the
production one.

The `e2e-test` job automates the end-to-end testing process for the project's
components. The job is triggered by the completion of the `deploy` job, and uses
the `newman` tool to execute the end-to-end tests and a custom shell script to
check the availability of the endpoints.


## Dependabot

This is a `dependabot.yml` file that specifies the update schedule for two
package ecosystems: `gomod` and `github-actions`. The `updates` key contains a
list of dictionaries, where each dictionary specifies the package ecosystem, the
directory to update, and the update schedule. In this case, both ecosystems are
set to update daily, and the root directory is specified with the forward slash.

The `gomod` package ecosystem is used for Go modules, which are a way to manage
dependencies in Go projects. Dependabot is a tool that automates dependency
updates, so this file is telling Dependabot to check for updates to Go modules
and GitHub Actions on a daily basis.

Overall, this file is a configuration file that helps automate the process of
keeping dependencies up to date. By specifying the package ecosystems and update
schedule, Dependabot can automatically check for updates and create pull
requests to update the dependencies in the project.
