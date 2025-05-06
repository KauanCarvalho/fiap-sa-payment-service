package ports

import (
	"context"

	"github.com/KauanCarvalho/fiap-sa-payment-service/internal/core/domain/entities"
)

type PaymentRepository interface {
	FindPaymentByExternalReference(ctx context.Context, externalReference string) (*entities.Payment, error)
	CreatePayment(ctx context.Context, payment *entities.Payment) error
	UpdatePayment(ctx context.Context, payment *entities.Payment) error
}
