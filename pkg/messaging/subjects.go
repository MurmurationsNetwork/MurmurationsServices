package messaging

const (
	// NodeCreated is the subject for an event where a node has been created.
	NodeCreated = "node:created"

	// NodeValidated is the subject for an event where a node has been
	// successfully validated.
	NodeValidated = "node:validated"

	// NodeValidationFailed is the subject for an event where a node's validation
	// has failed.
	NodeValidationFailed = "node:validation_failed"
)
