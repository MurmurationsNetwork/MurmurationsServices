package event

type Subject string

const (
	NodeCreated          Subject = "node:created"
	NodeValidated        Subject = "node:validated"
	NodeValidationFailed Subject = "node:validation_failed"
)

var SubjectsList = []string{
	string(NodeCreated),
	string(NodeValidated),
	string(NodeValidationFailed),
}
