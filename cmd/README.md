# `/cmd/`

## `<service_name>/main.go`

This is a Go program that defines the `main` package for a service. In Go, the
`main` package is a special package that defines the entry point of an
executable program. The `main` package must have a `main` function defined,
which is the first function that gets executed when the program is run. The
`main` function can call other functions and packages to perform the desired
functionality of the program.

The import statements import two packages: `logger` and the named service. The
`logger` package is used for logging messages, while the other contains the
implementation of the named service.

The `main` function by logging a message to indicate that the service is
starting up.

It then creates a new instance of the `<service_name>.Service` struct using the
`<service_name>.NewService` function. This function initializes the service and
returns a pointer to the `Service` struct.

It then starts a new goroutine that waits for the service to start up using the
`WaitUntilUp` method of the `Service` struct. This method returns a channel that
is closed when the service is up and running. The goroutine waits for the
channel to be closed and then logs a message to indicate that the service has
started.

Finally, it calls the `Run` method of the `Service` struct to start the service.
This method blocks until the service is shut down.
