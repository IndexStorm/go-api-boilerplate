package repository

import (
	"context"
	"github.com/rs/zerolog"
	"github.com/uptrace/bun"
	"time"
)

type postgresRepo struct {
	db     *bun.DB
	logger *zerolog.Logger
}

func (r *postgresRepo) Bootstrap() {
	var err error
	db := r.db
	logger := r.logger
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	_, err = db.NewCreateTable().IfNotExists().
		Model((*User)(nil)).
		Exec(ctx)
	if err != nil {
		logger.Fatal().Err(err).Msg("db:table:create:User")
	}
}
