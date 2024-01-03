# Messaging Package

## Overview
The `messaging` package provides an abstraction layer over NATS (a distributed messaging system) for Go-based microservices. It simplifies publishing and subscribing to various events within the Murmurations Network. This package handles the intricacies of NATS client setup, JSON marshaling, and message handling, allowing developers to focus on business logic.

## Key Features
- **Simplified Event Publishing**: Easily publish events without worrying about the underlying NATS client setup.
- **Event Subscriptions**: Subscribe to specific subjects with queue support for load-balanced message handling.
- **Structured Event Data**: Define and use structured data for events, enhancing code readability and maintainability.

## Usage

### Installation
To use the `messaging` package in your project, import it as follows:
```go
import "github.com/MurmurationsNetwork/MurmurationsServices/messaging"
```

### Publishing Events
You can publish events using either `Publish` or `PublishSync` functions. `Publish` is asynchronous, while `PublishSync` waits for an acknowledgment from the NATS server.

#### Example: Publishing an Event
```go
err := messaging.Publish(messaging.NodeCreated, eventData)
if err != nil {
    // handle error
}
```

#### Example: Synchronous Publishing
```go
err := messaging.PublishSync(messaging.NodeValidated, eventData)
if err != nil {
    // handle error
}
```

### Subscribing to Events
Use the `QueueSubscribe` function to subscribe to a specific subject. This function ensures load balancing across multiple instances of your service.

#### Example: Subscribing to an Event
```go
err := messaging.QueueSubscribe("subject", "queue", func(msg *nats.Msg) {
    // handle the message

    // acknowledge the message after successful processing.
    if ackErr := msg.Ack(); ackErr != nil {
    }
})
if err != nil {
    // handle error
}
```
