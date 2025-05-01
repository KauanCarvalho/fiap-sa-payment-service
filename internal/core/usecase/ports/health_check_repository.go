package ports

import "context"

type HealthCheckRepository interface {
	Ping(ctx context.Context) error
}
