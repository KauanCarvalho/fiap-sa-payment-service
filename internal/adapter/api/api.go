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
	awsCfg  *config.AWSConfig
	mongoDB *mongo.Client
}

func NewServer(cfg *config.Config, mongoDB *mongo.Client, awsCfg *config.AWSConfig) *Server {
	return &Server{
		cfg:     cfg,
		awsCfg:  awsCfg,
		mongoDB: mongoDB,
	}
}

func (s *Server) Run() {
	// Stores.
	ds := newStores(s.mongoDB, s.cfg.MongoDatabaseName)

	// Clients.
	snsClient, err := config.NewSNSClient(s.awsCfg)
	if err != nil {
		zap.L().Fatal("Failed to create SNS client",
			zap.String("aws_region", s.awsCfg.Region),
			zap.String("aws_endpoint", s.awsCfg.SNSConfig.Endpoint),
			zap.Error(err),
		)
	}

	// Usecases.
	authorizePayment := usecase.NewAuthorizePaymentUseCase(ds)
	updatePaymentStatus := usecase.NewUpdatePaymentStatusUseCase(ds, snsClient)

	// Web server.
	r := GenerateRouter(s.cfg, ds, authorizePayment, updatePaymentStatus)

	err = r.Run(fmt.Sprintf(":%s", s.cfg.Port))
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

func GenerateRouter(
	cfg *config.Config,
	ds domain.Datastore,
	authorizePayment usecase.AuthorizePaymentUseCase,
	updatePaymentStatus usecase.UpdatePaymentStatusUseCase,
) *gin.Engine {
	r := gin.New()

	setupMiddlewares(r, cfg)
	registerRoutes(r, cfg, ds, authorizePayment, updatePaymentStatus)

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

func registerRoutes(
	r *gin.Engine,
	cfg *config.Config,
	ds domain.Datastore,
	authorizePayment usecase.AuthorizePaymentUseCase,
	updatePaymentStatus usecase.UpdatePaymentStatusUseCase,
) {
	healthCheckHandler := handler.NewHealthCheckHandler(ds)
	paymentHandler := handler.NewPaymentHandler(authorizePayment, updatePaymentStatus)

	r.GET("/healthcheck", healthCheckHandler.Ping)

	apiV1 := r.Group("/api/v1")
	{
		payments := apiV1.Group("/payments")
		{
			payments.POST("/authorize", paymentHandler.Authorize)
			payments.PATCH("/:external_reference/update-status", paymentHandler.UpdateStatus)
		}
	}

	if cfg.IsDevelopment() {
		docs.SwaggerInfo.BasePath = ""
		r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	}
}
