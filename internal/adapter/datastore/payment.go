package datastore

import (
	"context"
	"errors"

	"github.com/KauanCarvalho/fiap-sa-payment-service/internal/core/domain/entities"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var ErrDuplicateExternalReference = errors.New("payment with this external_reference already exists")
var ErrPaymentNotFound = errors.New("payment not found")

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

func (ds *datastore) FindPaymentByExternalReference(ctx context.Context, externalReference string) (*entities.Payment, error) {
	ctx, cancel := context.WithTimeout(ctx, DefaultConnectionTimeout)
	defer cancel()

	collection := ds.db.Database(ds.databaseName).Collection(PaymentCollectionName)

	var result entities.Payment
	err := collection.FindOne(ctx, bson.M{"external_reference": externalReference}).Decode(&result)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrPaymentNotFound
		}
		return nil, err
	}

	return &result, nil
}

func (ds *datastore) UpdatePayment(ctx context.Context, p *entities.Payment) error {
	ctx, cancel := context.WithTimeout(ctx, DefaultConnectionTimeout)
	defer cancel()

	if p.ID.IsZero() {
		return errors.New("cannot update payment with empty ID")
	}

	collection := ds.db.Database(ds.databaseName).Collection(PaymentCollectionName)

	filter := bson.M{"_id": p.ID}
	update := bson.M{"$set": p}

	result, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return ErrPaymentNotFound
	}

	return nil
}
