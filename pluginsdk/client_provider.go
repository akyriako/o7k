package pluginsdk

import (
	"context"
	"sync"
)

type ClientProvider[T any] struct {
	host    Host
	connect func(context.Context, Context) (T, error)

	mu         sync.Mutex
	generation uint64
	client     T
	valid      bool
}

func NewClientProvider[T any](host Host, connect func(context.Context, Context) (T, error)) *ClientProvider[T] {
	return &ClientProvider[T]{
		host:    host,
		connect: connect,
	}
}

func (p *ClientProvider[T]) Client(ctx context.Context) (T, error) {
	current, err := p.host.Context(ctx)
	if err != nil {
		var zero T
		return zero, err
	}

	p.mu.Lock()
	defer p.mu.Unlock()

	if p.valid && p.generation == current.Generation {
		return p.client, nil
	}

	client, err := p.connect(ctx, current)
	if err != nil {
		var zero T
		return zero, err
	}

	p.client = client
	p.generation = current.Generation
	p.valid = true

	return p.client, nil
}
