package datastore

import (
	"time"

	"github.com/KauanCarvalho/fiap-sa-payment-service/internal/core/domain"
	"go.mongodb.org/mongo-driver/mongo"
)

const DefaultConnectionTimeout = 5 * time.Second

type datastore struct {
	db *mongo.Client
}

func NewDatastore(db *mongo.Client) domain.Datastore {
	return &datastore{db: db}
}
