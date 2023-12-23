package nats

import "github.com/nats-io/nats.go"

// mockClient can be used for testing purposes.
type mockClient struct{}

func (m *mockClient) JetStream() nats.JetStreamContext {
	return nil
}

func (m *mockClient) Disconnect() {
	// Mock implementation
}

func (m *mockClient) setClient(_ *nats.Conn, _ nats.JetStreamContext) {
	// Mock implementation
}
