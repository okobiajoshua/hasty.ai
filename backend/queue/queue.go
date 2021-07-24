package queue

import "context"

type Queue interface {
	// Consume(ctx context.Context, f func([]byte) error)
	Publish(ctx context.Context, msg []byte) error
}
