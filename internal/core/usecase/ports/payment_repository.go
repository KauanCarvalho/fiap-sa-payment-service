package ports

import (
	"context"

	"github.com/KauanCarvalho/fiap-sa-payment-service/internal/core/domain/entities"
)

type PaymentRepository interface {
	CreatePayment(ctx context.Context, payment *entities.Payment) error
}
