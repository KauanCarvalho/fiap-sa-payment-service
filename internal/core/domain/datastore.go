package domain

import (
	"github.com/KauanCarvalho/fiap-sa-payment-service/internal/core/usecase/ports"
)

type Datastore interface {
	ports.HealthCheckRepository
}
