package event

type Subject string

const (
	// NodeCreated is the subject for an event where a node has been created.
	NodeCreated Subject = "node:created"

	// NodeValidated is the subject for an event where a node has been
	// successfully validated.
	NodeValidated Subject = "node:validated"

	// NodeValidationFailed is the subject for an event where a node's validation
	// has failed.
	NodeValidationFailed Subject = "node:validation_failed"
)
