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

func (m *MockQueue) Consume(ctx context.Context, f func(context.Context, []byte) error) {
	m.Called(ctx, f)
}
