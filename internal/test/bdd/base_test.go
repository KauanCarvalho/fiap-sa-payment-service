package bdd_test

import (
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/KauanCarvalho/fiap-sa-payment-service/internal/adapter/api"
	"github.com/KauanCarvalho/fiap-sa-payment-service/internal/adapter/datastore"
	"github.com/KauanCarvalho/fiap-sa-payment-service/internal/config"
	"github.com/KauanCarvalho/fiap-sa-payment-service/internal/core/domain"
	"github.com/KauanCarvalho/fiap-sa-payment-service/internal/core/usecase"
	"github.com/KauanCarvalho/fiap-sa-payment-service/internal/di"
	"github.com/cucumber/godog"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"

	"go.uber.org/zap"
)

var (
	engine   *gin.Engine
	cfg      *config.Config
	recorder *httptest.ResponseRecorder
	request  *http.Request
	mongoDB  *mongo.Client
	ds       domain.Datastore
	ap       usecase.AuthorizePaymentUseCase
	bodyData map[string]any
)

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)

	logger := zap.NewNop()
	zap.ReplaceGlobals(logger)

	cfg = config.Load()

	var err error
	mongoDB, err = di.NewMongoConnection(cfg)
	if err != nil {
		log.Fatalf("error when initializing database connection: %v", err)
	}

	ds = datastore.NewDatastore(mongoDB, cfg.MongoDatabaseName)
	ap = usecase.NewAuthorizePaymentUseCase(ds)
	engine = api.GenerateRouter(cfg, ds, ap)

	code := m.Run()

	os.Exit(code)
}

func resetState() {
	recorder = httptest.NewRecorder()
	bodyData = make(map[string]any)
}

func ResetAndLoadFixtures() {
	resetState()
}

func TestFeatures(t *testing.T) {
	suite := godog.TestSuite{
		Name: "Features",
		ScenarioInitializer: func(sc *godog.ScenarioContext) {
			InitializeScenarioPaymentAPI(sc)
		},
		Options: &godog.Options{
			Format: "pretty",
			Paths:  []string{"./features"},
		},
	}

	if suite.Run() != 0 {
		t.Fatal("tests failed")
	}
}
