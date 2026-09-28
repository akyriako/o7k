package pluginsdk

import (
	"context"
	"sync"
)

type clientState[T any] struct {
	generation     uint64
	client         T
	serviceClients sync.Map
}

type ClientProvider[T any] struct {
	host    Host
	connect func(context.Context, Context) (T, error)

	mu    sync.Mutex
	state *clientState[T]
}

func NewClientProvider[T any](host Host, connect func(context.Context, Context) (T, error)) *ClientProvider[T] {
	return &ClientProvider[T]{
		host:    host,
		connect: connect,
	}
}

func (p *ClientProvider[T]) stateFor(ctx context.Context) (*clientState[T], Context, error) {
	current, err := p.host.Context(ctx)
	if err != nil {
		return nil, Context{}, err
	}

	p.mu.Lock()
	defer p.mu.Unlock()

	if p.state != nil && p.state.generation == current.Generation {
		return p.state, current, nil
	}

	client, err := p.connect(ctx, current)
	if err != nil {
		return nil, Context{}, err
	}

	state := &clientState[T]{
		generation: current.Generation,
		client:     client,
	}

	p.state = state

	return state, current, nil
}

func (p *ClientProvider[T]) Client(ctx context.Context) (T, error) {
	state, _, err := p.stateFor(ctx)
	if err != nil {
		var zero T
		return zero, err
	}

	return state.client, nil
}

func GetServiceClient[T any, S any](ctx context.Context, provider *ClientProvider[T], key string, builder func(T, Context) (S, error)) (S, error) {
	state, current, err := provider.stateFor(ctx)
	if err != nil {
		var zero S
		return zero, err
	}

	if client, ok := state.serviceClients.Load(key); ok {
		return client.(S), nil
	}

	client, err := builder(state.client, current)
	if err != nil {
		var zero S
		return zero, err
	}

	actual, _ := state.serviceClients.LoadOrStore(key, client)
	return actual.(S), nil
}
