package event

type Subject string

const (
	nodeCreated          Subject = "node:created"
	nodeValidated        Subject = "node:validated"
	nodeValidationFailed Subject = "node:validation_failed"
)
