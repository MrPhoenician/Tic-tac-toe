package model

import "github.com/jackc/pgx/v5/pgtype"

type User struct {
	ID           string           `db:"id"`
	Login        string           `db:"login"`
	PasswordHash string           `db:"password_hash"`
	CreatedAt    pgtype.Timestamp `db:"created_at,omitempty"`
}
