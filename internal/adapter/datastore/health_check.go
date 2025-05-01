package datastore

import (
	"context"
)

func (ds *datastore) Ping(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, DefaultConnectionTimeout)
	defer cancel()

	return ds.db.Ping(ctx, nil)
}
