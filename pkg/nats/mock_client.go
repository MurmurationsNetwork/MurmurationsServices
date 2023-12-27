package nats

import (
	"github.com/nats-io/nats.go"

	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/event"
)

// mockClient can be used for testing purposes.
type mockClient struct{}

func (m *mockClient) JetStream() nats.JetStreamContext {
	return nil
}

func (m *mockClient) Disconnect() {
	// Mock implementation
}

func (m *mockClient) SubscribeToSubjects(_ ...event.Subject) error {
	return nil
}
