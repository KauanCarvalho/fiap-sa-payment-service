package mock

import (
	"context"
	"errors"
)

type DatastoreMock struct {
	PingFn func(ctx context.Context) error
}

var ErrFunctionNotImplemented = errors.New("function not implemented")

func (m *DatastoreMock) Ping(ctx context.Context) error {
	if m.PingFn != nil {
		return m.PingFn(ctx)
	}

	return ErrFunctionNotImplemented
}
