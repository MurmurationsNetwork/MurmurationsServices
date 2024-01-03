package messaging

// Constants for NATS subjects in the messaging system. All subjects should
// follow the format "NODES.<event>", where <event> describes the specific
// nature of the event.

const (
	// NodeCreated is the subject for an event where a node has been created.
	NodeCreated = "NODES.created"

	// NodeValidated is the subject for an event where a node has been
	// successfully validated.
	NodeValidated = "NODES.validated"

	// NodeValidationFailed is the subject for an event where a node's validation
	// has failed.
	NodeValidationFailed = "NODES.validation_failed"
)
