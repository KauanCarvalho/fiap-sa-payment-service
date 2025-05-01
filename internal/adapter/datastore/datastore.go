package datastore

import (
	"context"
	"fmt"
	"time"

	"github.com/KauanCarvalho/fiap-sa-payment-service/internal/core/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

const DefaultConnectionTimeout = 5 * time.Second
const PaymentCollectionName = "payments"

type datastore struct {
	db           *mongo.Client
	databaseName string
}

func NewDatastore(db *mongo.Client, databaseName string) domain.Datastore {
	datastore := &datastore{
		db:           db,
		databaseName: databaseName,
	}

	err := datastore.ensureIndexes()
	if err != nil {
		zap.L().Info("Failed to ensure indexes")
	}

	return datastore
}

func (ds *datastore) ensureIndexes() error {
	collection := ds.db.Database(ds.databaseName).Collection(PaymentCollectionName)

	indexModel := mongo.IndexModel{
		Keys:    bson.M{"external_reference": 1},
		Options: options.Index().SetUnique(true),
	}

	_, err := collection.Indexes().CreateOne(context.Background(), indexModel)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return nil
		}
		return fmt.Errorf("failed to create index: %w", err)
	}

	return nil
}
