package repository

import (
	"context"
	"database/sql"
	"errors"
	"github.com/uptrace/bun"
)

type UserStatus uint8

const (
	UserStatusInactive UserStatus = iota + 1
	UserStatusActive
	UserStatusBlocked
)

type User struct {
	bun.BaseModel `bun:"table:users"`
	PrivateID

	Email    string     `bun:",notnull,unique,nullzero" json:"email,omitempty"`
	Password []byte     `bun:",notnull,nullzero" json:"-"`
	Status   UserStatus `bun:",notnull,nullzero" json:"-"`

	DateMixin
	SoftDeleteMixin
}

func (r *postgresRepo) SaveUser(user *User, ctx context.Context) error {
	r.sanitizeCreateModel(user)
	_, err := r.db.NewInsert().Model(user).Exec(ctx)
	return err
}

func (r *postgresRepo) FetchUserByID(id string, ctx context.Context) (u *User, err error) {
	u = new(User)
	err = r.db.NewSelect().Model(u).Where("public_id = ?", id).Scan(ctx)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	return u, err
}
