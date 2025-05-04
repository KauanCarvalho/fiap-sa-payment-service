package datastore_test

import (
	"context"
	"testing"

	"github.com/KauanCarvalho/fiap-sa-payment-service/internal/adapter/datastore"
	"github.com/KauanCarvalho/fiap-sa-payment-service/internal/core/domain/entities"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestCreatePayment(t *testing.T) {
	t.Run("successfully insert payment", func(t *testing.T) {
		p := &entities.Payment{
			Amount:            1500,
			Status:            "authorized",
			ExternalReference: uuid.New().String(),
			Provider:          "pix",
			PaymentMethod:     "qr_code",
			QRCode:            "qr_code_data",
		}

		err := ds.CreatePayment(ctx, p)
		require.NoError(t, err, "CreatePayment should not return error")

		var result entities.Payment
		collection := mongoDB.Database(cfg.MongoDatabaseName).Collection("payments")
		err = collection.FindOne(ctx, bson.M{"_id": p.ID}).Decode(&result)
		require.NoError(t, err, "Should find inserted payment")

		assert.InEpsilon(t, p.Amount, result.Amount, 0.01)
		assert.Equal(t, p.Status, result.Status)
		assert.Equal(t, p.ExternalReference, result.ExternalReference)
		assert.Equal(t, p.Provider, result.Provider)
		assert.Equal(t, p.PaymentMethod, result.PaymentMethod)
		assert.Equal(t, p.QRCode, result.QRCode)
	})

	t.Run("fail to insert payment with context timeout", func(t *testing.T) {
		timeoutCtx, cancel := context.WithTimeout(ctx, 0)
		defer cancel()

		p := &entities.Payment{
			Amount:            1500,
			Status:            "pending",
			ExternalReference: uuid.New().String(),
			Provider:          "pix",
			PaymentMethod:     "qr_code",
			QRCode:            "timeout_qr",
		}

		err := ds.CreatePayment(timeoutCtx, p)
		assert.Error(t, err, "Should return an error due to context timeout")
	})

	t.Run("fail to insert payment due to duplicate external_reference", func(t *testing.T) {
		p1 := &entities.Payment{
			Amount:            1500,
			Status:            "authorized",
			ExternalReference: uuid.New().String(),
			Provider:          "pix",
			PaymentMethod:     "qr_code",
			QRCode:            "qr_code_1",
		}

		err := ds.CreatePayment(ctx, p1)
		require.NoError(t, err, "CreatePayment should not return error for first insertion")

		p2 := &entities.Payment{
			Amount:            2000,
			Status:            "authorized",
			ExternalReference: p1.ExternalReference,
			Provider:          "pix",
			PaymentMethod:     "qr_code",
			QRCode:            "qr_code_2",
		}

		err = ds.CreatePayment(ctx, p2)
		require.ErrorIs(t, datastore.ErrDuplicateExternalReference, err)
	})
}

func TestFindPaymentByExternalReference(t *testing.T) {
	t.Run("successfully find payment by external_reference", func(t *testing.T) {
		externalRef := uuid.New().String()
		p := &entities.Payment{
			Amount:            500,
			Status:            "pending",
			ExternalReference: externalRef,
			Provider:          "pix",
			PaymentMethod:     "qr_code",
			QRCode:            "find_qr",
		}
		err := ds.CreatePayment(ctx, p)
		require.NoError(t, err)

		result, err := ds.FindPaymentByExternalReference(ctx, externalRef)
		require.NoError(t, err)
		require.NotNil(t, result)

		assert.Equal(t, p.ID, result.ID)
		assert.Equal(t, p.ExternalReference, result.ExternalReference)
		assert.Equal(t, p.QRCode, result.QRCode)
	})

	t.Run("return error when payment not found", func(t *testing.T) {
		_, err := ds.FindPaymentByExternalReference(ctx, "nonexistent-ref")
		require.ErrorIs(t, err, datastore.ErrPaymentNotFound)
	})
}

func TestUpdatePayment(t *testing.T) {
	t.Run("successfully update payment", func(t *testing.T) {
		p := &entities.Payment{
			Amount:            800,
			Status:            "pending",
			ExternalReference: uuid.New().String(),
			Provider:          "pix",
			PaymentMethod:     "qr_code",
			QRCode:            "initial_qr",
		}
		err := ds.CreatePayment(ctx, p)
		require.NoError(t, err)

		// Modify and update
		p.Status = "authorized"
		p.QRCode = "updated_qr"
		err = ds.UpdatePayment(ctx, p)
		require.NoError(t, err)

		// Verify
		var result entities.Payment
		err = mongoDB.Database(cfg.MongoDatabaseName).
			Collection("payments").
			FindOne(ctx, bson.M{"_id": p.ID}).Decode(&result)
		require.NoError(t, err)

		assert.Equal(t, "authorized", result.Status)
		assert.Equal(t, "updated_qr", result.QRCode)
	})

	t.Run("fail to update with empty ID", func(t *testing.T) {
		p := &entities.Payment{
			Amount:            999,
			Status:            "pending",
			ExternalReference: uuid.New().String(),
			Provider:          "pix",
			PaymentMethod:     "qr_code",
			QRCode:            "invalid",
		}
		err := ds.UpdatePayment(ctx, p)
		require.Error(t, err)
	})

	t.Run("fail to update nonexistent payment", func(t *testing.T) {
		p := &entities.Payment{
			ID:                primitive.NewObjectID(),
			Amount:            123,
			Status:            "pending",
			ExternalReference: "not_exists",
			Provider:          "pix",
			PaymentMethod:     "qr_code",
			QRCode:            "invalid_qr",
		}
		err := ds.UpdatePayment(ctx, p)
		require.ErrorIs(t, err, datastore.ErrPaymentNotFound)
	})
}
