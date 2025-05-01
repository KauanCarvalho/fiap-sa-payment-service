package datastore_test

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/KauanCarvalho/fiap-sa-payment-service/internal/adapter/datastore"
	"github.com/KauanCarvalho/fiap-sa-payment-service/internal/config"
	"github.com/KauanCarvalho/fiap-sa-payment-service/internal/core/domain"
	"github.com/KauanCarvalho/fiap-sa-payment-service/internal/di"

	"go.mongodb.org/mongo-driver/mongo"
)

var (
	ctx     context.Context
	cfg     *config.Config
	mongoDB *mongo.Client
	ds      domain.Datastore
)

func TestMain(m *testing.M) {
	ctx = context.Background()
	cfg = config.Load()

	var err error
	mongoDB, err = di.NewMongoConnection(cfg)
	if err != nil {
		log.Fatalf("error when creating database connection pool: %v", err)
	}

	ds = datastore.NewDatastore(mongoDB, cfg.MongoDatabaseName)

	os.Exit(m.Run())
}
