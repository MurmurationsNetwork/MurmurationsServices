# This is a Makefile that defines a target called docker-build-library. This
# target specifies a command to build a Docker image for a library package. The
# command is executed by running docker build with the -f flag to specify the
# Dockerfile to use, and the -t flag to specify the name and tag for the
# resulting image. In this case, the Dockerfile is located at
# build/library/package/Dockerfile-prod and the resulting image is named
# murmurations/library.

# The purpose of this Makefile target is to automate the process of building a
# Docker image for the library package. By running make docker-build-library,
# the command specified in the target will be executed, resulting in a Docker
# image that can be used to deploy the library package.

docker-build-library:
	docker build -f build/library/package/Dockerfile-prod -t murmurations/library .
