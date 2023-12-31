# NATS Client Package

The `natsclient` package package provides a simplified interaction with NATS, especially focusing on JetStream functionalities. It utilizes the Singleton pattern for efficient and consistent management of NATS connections across different services in the Murmurations Network.

## Features

- Singleton management of NATS client connections.
- Easy subscription to multiple subjects with JetStream.
- Automated stream and consumer management.
- Graceful disconnection with cleanup of resources.

## Getting Started

### Prerequisites

Before using the `natsclient` package, ensure that:
- You are using Go version 1.x or higher.
- The NATS server is accessible at the specified URL.

### Usage

#### Initializing the Client

Initialize the NATS client at the start of your application:

```go
err := natsclient.Initialize("nats://your-nats-server:4222")
if err != nil {
    log.Fatalf("Failed to initialize NATS client: %v", err)
}
```

#### Getting the Client Instance

Retrieve the singleton instance of the NATS client:

```go
client := natsclient.GetInstance()
```

#### Subscribing to Subjects

To subscribe to specific subjects:

```go
err := client.SubscribeToSubjects("subject1", "subject2")
if err != nil {
    log.Printf("Failed to subscribe to subjects: %v", err)
}
```

#### Disconnecting the Client

Properly disconnect the client when needed:

```go
err := client.Disconnect()
if err != nil {
    log.Printf("Error during client disconnection: %v", err)
}
```

## Advanced Usage

- **Custom Stream Configurations:** Modify `createStream` to tailor the stream configurations to your specific requirements.
- **Error Handling:** Comprehensive error handling is provided to ensure smooth operation and troubleshooting.

