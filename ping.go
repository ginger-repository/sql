package sql

import (
	"context"

	"github.com/ginger-core/errors"
)

func (repo *repo) Ping(ctx context.Context) errors.Error {
	sdb, err := repo.DB.DB()
	if err != nil {
		return errors.New(err).WithTrace("repo.DB.DB")
	}
	if err := sdb.PingContext(ctx); err != nil {
		return errors.New(err).WithTrace("sdb.PingContext")
	}
	return nil
}
