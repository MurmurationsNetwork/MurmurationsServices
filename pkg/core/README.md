# Core Package

## Overview

The `core` package serves as the foundation for setting up and managing server applications in Go.

## Key Components

### `InstallShutdownHandler`

`InstallShutdownHandler` handles system signals and executes shutdown tasks in an orderly manner. This function takes another function as an argument, which is executed when the application receives an interrupt (`SIGINT`) or terminate (`SIGTERM`) signal.

Usage:

```go
core.InstallShutdownHandler(func() {
	// Your shutdown logic here
})
```

By using `InstallShutdownHandler`, you ensure that all necessary shutdown tasks are taken care of before your service is terminated. Note that the shutdown function provided will run asynchronously in a new goroutine, allowing the main function to continue its operations until a shutdown signal is received.
