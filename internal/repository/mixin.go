package repository

import "time"

type PrivateID struct {
	ID       uint64 `bun:",pk,autoincrement,identity" json:"-"`
	PublicID string `bun:",notnull,nullzero,unique" json:"id,omitempty"`
}

type DateMixin struct {
	CreatedAt time.Time `bun:",notnull,nullzero,default:current_timestamp" json:"created_at,omitempty"`
	UpdatedAt time.Time `bun:",notnull,nullzero,default:current_timestamp" json:"updated_at,omitempty"`
}

type SoftDeleteMixin struct {
	DeletedAt *time.Time `bun:",soft_delete,nullzero" json:"deleted_at,omitempty"`
}
