package queue

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type MockQueue struct {
	mock.Mock
}

func NewMockQueue() *MockQueue {
	return &MockQueue{}
}

func (m *MockQueue) Publish(ctx context.Context, msg []byte) error {
	args := m.Called(ctx, msg)
	return args.Error(0)
}
