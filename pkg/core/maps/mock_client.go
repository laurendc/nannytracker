package maps

import (
	"context"
)

// MockClient is a mock implementation of the maps client for testing
type MockClient struct {
	// MockDistance is the distance that will be returned by CalculateDistance
	MockDistance float64
}

// NewMockClient creates a new mock client for testing
func NewMockClient() *MockClient {
	return &MockClient{
		MockDistance: 10.0, // Default mock distance
	}
}

// CalculateDistance returns the mock distance
func (m *MockClient) CalculateDistance(ctx context.Context, origin, destination string) (float64, error) {
	return m.MockDistance, nil
}
