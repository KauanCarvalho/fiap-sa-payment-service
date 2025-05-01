package api

import (
	"fmt"

	"github.com/KauanCarvalho/fiap-sa-payment-service/internal/adapter/api/handler"
	"github.com/KauanCarvalho/fiap-sa-payment-service/internal/adapter/api/middleware"
	"github.com/KauanCarvalho/fiap-sa-payment-service/internal/adapter/datastore"
	"github.com/KauanCarvalho/fiap-sa-payment-service/internal/config"
	"github.com/KauanCarvalho/fiap-sa-payment-service/internal/core/domain"
	"github.com/KauanCarvalho/fiap-sa-payment-service/internal/core/usecase"
	docs "github.com/KauanCarvalho/fiap-sa-payment-service/swagger"

	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

type Server struct {
	cfg     *config.Config
	mongoDB *mongo.Client
}

func NewServer(cfg *config.Config, mongoDB *mongo.Client) *Server {
	return &Server{
		cfg:     cfg,
		mongoDB: mongoDB,
	}
}

func (s *Server) Run() {
	// Stores.
	ds := newStores(s.mongoDB, s.cfg.MongoDatabaseName)

	// Usecases.
	authorizePayment := usecase.NewAuthorizePaymentUseCase(ds)

	// Web server.
	r := GenerateRouter(s.cfg, ds, authorizePayment)

	err := r.Run(fmt.Sprintf(":%s", s.cfg.Port))
	if err != nil {
		zap.L().Fatal("Failed to start server",
			zap.String("port", s.cfg.Port),
			zap.Error(err),
		)
	}
}

func newStores(mongoDB *mongo.Client, databaseName string) domain.Datastore {
	return datastore.NewDatastore(mongoDB, databaseName)
}

func GenerateRouter(cfg *config.Config, ds domain.Datastore, authorizePayment usecase.AuthorizePaymentUseCase) *gin.Engine {
	r := gin.New()

	setupMiddlewares(r, cfg)
	registerRoutes(r, cfg, ds, authorizePayment)

	return r
}

func setupMiddlewares(r *gin.Engine, cfg *config.Config) {
	r.Use(
		middleware.Logger(),
		middleware.Recovery(),
	)

	r.Use(requestid.New(
		requestid.WithGenerator(func() string {
			return cfg.AppName + "-" + uuid.New().String()
		}),
	))
}

func registerRoutes(r *gin.Engine, cfg *config.Config, ds domain.Datastore, authorizePayment usecase.AuthorizePaymentUseCase) {
	healthCheckHandler := handler.NewHealthCheckHandler(ds)
	paymentHandler := handler.NewPaymentHandler(authorizePayment)

	r.GET("/healthcheck", healthCheckHandler.Ping)

	apiV1 := r.Group("/api/v1")
	{
		payments := apiV1.Group("/payments")
		{
			payments.POST("/authorize", paymentHandler.Authorize)
		}
	}

	if cfg.IsDevelopment() {
		docs.SwaggerInfo.BasePath = ""
		r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	}
}
