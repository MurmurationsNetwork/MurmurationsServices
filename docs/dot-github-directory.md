# `/.github`

## GitHub Workflows

GitHub Workflows is a feature of GitHub Actions that automates tasks and
workflows in a repository. A workflow is defined in a YAML file in the
[.github/workflows directory](/.github/workflows) of a repository and can be
triggered by events such as pushes, pull requests or scheduled intervals.
Workflows run on GitHub-hosted virtual machines or on self-hosted runners.
MurmurationsServices uses Workflows to build, test and deploy code.

### `main.yaml`

This workflow is triggered by a push event on the `main` branch and automates
the test, lint, build and deployment processes for the project's components. The
workflow uses `make` commands to build and publish Docker images for each
component and deploys the components to a Kubernetes cluster.

### `pr.yaml`

This workflow defines the jobs that run tests, lint, build and deploy the
project when a pull request is opened or updated on the main branch.
Importantly, it also checks file changes and determine which services should be
excluded from the build and deployment processes, which speeds up deployment.
Automatically running tests and linting on every pull request helps ensure that
the code is high quality and free of errors.

### `redeploy.yaml`

This workflow is initiated by specific review comments from authorized users in
pull request reviews. It automates the build and deployment processes when a
comment containing `/rebuild` is detected. This enables the deployment of new
code to the pretest environment, where code reviewers can test out changes
committed to the pull request.

### `test.yaml`

This workflow is exactly the same as the `main.yaml` workflow above except that
the deployment is made to a test environment, not the live one.

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
