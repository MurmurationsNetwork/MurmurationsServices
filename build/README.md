# `/build/`

## `<service_name>/mk/Makefile`

This `Makefile` file defines a target called `docker-build-<service_name>` that
specifies a command to build a Docker image for a service. The command is
executed by running `docker build` with the `-f` flag to specify the Dockerfile
to use, and the `-t` flag to specify the name and tag for the resulting image.

The purpose of a Makefile target is to automate the process of building a Docker
image for a service. By running `make docker-build-<service_name>`, the command
specified in the target will be executed, resulting in a Docker image that can
be used to deploy the service.

## `<service_name>/package/Dockerfile-dev`

This is a Dockerfile that defines a development environment for a service. 

The `FROM` command specifies the base image to use, which in this case is
`golang:1.19.5-alpine`. Alpine is a lightweight Linux distribution, which makes
it a good choice for Docker images. The `golang:1.19.5-alpine` image is an image
that contains the Go programming language.

The `RUN` command updates the package index inside the container using the `apk`
package manager, which is similar to `apt-get` on Ubuntu.

The `WORKDIR` command sets the working directory inside the container to
`/src/<service_name>`. This is where the application code will be copied to.

The `COPY` command copies the `go.mod` and `go.sum` files to the working
directory. The `go.mod` file contains the dependencies for the application and
the `go.sum` file contains the checksums for the dependencies.

The `RUN` command then downloads the dependencies specified in the `go.mod` file
using `go mod download`. This command is run before copying the application code
to take advantage of Docker's caching mechanism. If the `go.mod` file hasn't
changed, then the dependencies don't need to be downloaded again, speeding up
the build process.

The next `RUN` command installs the latest version of `reflex`, which is a tool
for hot reloading Go applications. Hot reloading allows developers to see
changes to their code immediately without having to manually restart the
application. The `go install` command installs reflex from the
[github.com/cespare/reflex](https://github.com/cespare/reflex) repository.

The `CMD` command specifies the command to run when the container starts. In
this case, it runs `reflex` with the configuration file located at
`/build/<service_name>/reflex.conf`. This will start the application using
`reflex` for hot reloading when the Docker container is run.

## `<service_name>/package/Dockerfile-prod`

This is a Dockerfile that defines a production environment for a service. 

The build stage compiles the Go app and creates a fully static binary, which is
then copied to the runtime stage. The runtime stage sets up the environment for
running the app and copies the static binary and files to the appropriate
locations. By using a Docker image, the app can be easily deployed to different
environments without worrying about dependencies or configuration.

## `<service_name>/reflex.conf`

The `reflex.conf` file is used with the reflex tool, which is a command-line
tool for running commands when files change. The `-r` flag specifies a regular
expression to match files, and the `-R` flag specifies a regular expression to
exclude files. In this case, the regular expression `-R '_test\.go$'` excludes
files with a `_test.go` suffix.

When a file with a `.go` extension is modified, the `go run ./cmd/library`
command is executed. This command compiles and runs the library service in the
cmd directory of the project.
