package mock

import (
	"context"
	"errors"

	"github.com/KauanCarvalho/fiap-sa-payment-service/internal/core/domain/entities"
)

type DatastoreMock struct {
	PingFn          func(ctx context.Context) error
	CreatePaymentFn func(ctx context.Context, payment *entities.Payment) error
}

var ErrFunctionNotImplemented = errors.New("function not implemented")

func (m *DatastoreMock) Ping(ctx context.Context) error {
	if m.PingFn != nil {
		return m.PingFn(ctx)
	}

	return ErrFunctionNotImplemented
}

func (m *DatastoreMock) CreatePayment(ctx context.Context, payment *entities.Payment) error {
	if m.CreatePaymentFn != nil {
		return m.CreatePaymentFn(ctx, payment)
	}

	return ErrFunctionNotImplemented
}
