package helper

import (
	"context"
	"strings"
)

const ModeMonolith = "monolith"

type Endpoint[T any] func(context.Context, T) error

type ModeDispatcher[T any] struct {
	mode         string
	monolith     Endpoint[T]
	microservice Endpoint[T]
}

func NewModeDispatcher[T any](mode string, monolith Endpoint[T], microservice Endpoint[T]) *ModeDispatcher[T] {
	return &ModeDispatcher[T]{
		mode:         normalizeMode(mode),
		monolith:     monolith,
		microservice: microservice,
	}
}

func (d *ModeDispatcher[T]) Dispatch(ctx context.Context, req T) error {
	if d == nil {
		return nil
	}
	if d.mode == ModeMonolith {
		return callEndpoint(ctx, req, d.monolith)
	}
	return callEndpoint(ctx, req, d.microservice)
}

func callEndpoint[T any](ctx context.Context, req T, endpoint Endpoint[T]) error {
	if endpoint == nil {
		return nil
	}
	return endpoint(ctx, req)
}

func normalizeMode(mode string) string {
	return strings.ToLower(strings.TrimSpace(mode))
}
