package datastore

import (
	"context"
	"errors"

	"github.com/KauanCarvalho/fiap-sa-payment-service/internal/core/domain/entities"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var ErrDuplicateExternalReference = errors.New("payment with this external_reference already exists")

func (ds *datastore) CreatePayment(ctx context.Context, p *entities.Payment) error {
	ctx, cancel := context.WithTimeout(ctx, DefaultConnectionTimeout)
	defer cancel()

	if p.ID.IsZero() {
		p.ID = primitive.NewObjectID()
	}

	collection := ds.db.Database(ds.databaseName).Collection(PaymentCollectionName)

	_, err := collection.InsertOne(ctx, p)

	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return ErrDuplicateExternalReference
		}
		return err
	}

	return nil
}
