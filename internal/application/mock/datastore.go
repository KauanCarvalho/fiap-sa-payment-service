package mock

import (
	"context"
	"errors"

	"github.com/KauanCarvalho/fiap-sa-payment-service/internal/core/domain/entities"
)

type DatastoreMock struct {
	PingFn                           func(ctx context.Context) error
	CreatePaymentFn                  func(ctx context.Context, payment *entities.Payment) error
	FindPaymentByExternalReferenceFn func(ctx context.Context, externalReference string) (*entities.Payment, error)
	UpdatePaymentFn                  func(ctx context.Context, payment *entities.Payment) error
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

func (m *DatastoreMock) FindPaymentByExternalReference(ctx context.Context, externalReference string) (*entities.Payment, error) {
	if m.FindPaymentByExternalReferenceFn != nil {
		return m.FindPaymentByExternalReferenceFn(ctx, externalReference)
	}

	return nil, ErrFunctionNotImplemented
}

func (m *DatastoreMock) UpdatePayment(ctx context.Context, payment *entities.Payment) error {
	if m.UpdatePaymentFn != nil {
		return m.UpdatePaymentFn(ctx, payment)
	}

	return ErrFunctionNotImplemented
}
