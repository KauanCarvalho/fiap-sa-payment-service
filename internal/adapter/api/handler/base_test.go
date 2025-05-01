package handler_test

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/KauanCarvalho/fiap-sa-payment-service/internal/adapter/api"
	"github.com/KauanCarvalho/fiap-sa-payment-service/internal/adapter/datastore"
	"github.com/KauanCarvalho/fiap-sa-payment-service/internal/config"
	"github.com/KauanCarvalho/fiap-sa-payment-service/internal/core/domain"
	"github.com/KauanCarvalho/fiap-sa-payment-service/internal/di"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

var (
	ctx       context.Context
	cfg       *config.Config
	mongoDB   *mongo.Client
	ds        domain.Datastore
	ginEngine *gin.Engine
)

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)

	logger := zap.NewNop()
	zap.ReplaceGlobals(logger)

	ctx = context.Background()
	cfg = config.Load()

	var err error
	mongoDB, err = di.NewMongoConnection(cfg)
	if err != nil {
		log.Fatalf("error when initializing database connection: %v", err)
	}

	ds = datastore.NewDatastore(mongoDB)

	ginEngine = api.GenerateRouter(cfg, ds)

	os.Exit(m.Run())
}
