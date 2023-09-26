package repository

import (
	"context"
	"github.com/rs/zerolog"
	"github.com/uptrace/bun"
)

type Repository interface {
	Bootstrap()

	SaveUser(user *User, ctx context.Context) error
	FetchUserByID(id string, ctx context.Context) (*User, error)
}

func NewPostgresRepo(db *bun.DB, logger zerolog.Logger) Repository {
	return &postgresRepo{db: db, logger: logger.With().Str("cat", "repo").Caller().Logger()}
}
