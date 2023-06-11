# Validation Service

## Overview

This service is designed to validate nodes created by users. It ensures that the nodes adhere to the specified schema.

## Key Components

- `internal/controller/event/listeners.go`: This component's responsibility is to listen for node creation events. When a node is created, it triggers the validation process.

- `internal/model/node.go`: This component represents a node stored in the index.

- `internal/service/validation_service.go`: This component contains the core logic for validating a node.

- `pkg/validation/service.go`: This component contains the setup for the validation service.

## Implementation Details

The service operates by first listening for node-created events. When a node is created, the event triggers the validation process. The node's profile is read from the provided URL and validated against a default schema. It is then validated against any schemas linked in the profile data. If the node's profile passes all validation requirements, a `NodeValidated` event is published.
