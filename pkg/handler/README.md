# Handler Package

The Handler package provides common handlers across all services.

## Handlers

This package includes the following handlers:

### Deprecation Handler

`DeprecationHandler` manages requests made to deprecated API versions. When a request is made to a deprecated version of the API (v1), the handler returns a JSON error message instructing the client to use the updated version (v2).

### Ping Handler

`PingHandler` responds to ping requests with "pong!". It is used primarily for checking the availability and responsiveness of the service. A successful "pong!" response indicates that the service is up and running.
