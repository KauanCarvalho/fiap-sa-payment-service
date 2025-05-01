package usecase

import (
	"context"

	"github.com/KauanCarvalho/fiap-sa-payment-service/internal/core/domain"
	"github.com/KauanCarvalho/fiap-sa-payment-service/internal/core/domain/entities"
	"github.com/KauanCarvalho/fiap-sa-payment-service/internal/core/usecase/dto"
	"github.com/google/uuid"
)

type AuthorizePaymentUseCase interface {
	Run(ctx context.Context, input dto.AuthorizePaymentInput) (*entities.Payment, error)
}

type authorizePaymentUsecase struct {
	ds domain.Datastore
}

func NewAuthorizePaymentUseCase(ds domain.Datastore) AuthorizePaymentUseCase {
	return &authorizePaymentUsecase{ds: ds}
}

func (c *authorizePaymentUsecase) Run(ctx context.Context, input dto.AuthorizePaymentInput) (*entities.Payment, error) {
	payment := &entities.Payment{
		Amount:            input.Amount,
		Status:            entities.PaymentStatusPending,
		ExternalReference: input.ExternalReference,
		Provider:          selectProvider(),
		PaymentMethod:     input.PaymentMethod,
		QRCode:            uuid.New().String(),
	}

	err := c.ds.CreatePayment(ctx, payment)
	if err != nil {
		return nil, err
	}

	return payment, nil
}

func selectProvider() string {
	return "MercadoPago"
}
