package database

import (
	"context"

	"github.com/jackc/pgx/v5"
)

type DBConnection interface {
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
}

type Repository struct {
	database DBConnection
}

func NewRepository(db DBConnection) *Repository {
	return &Repository{
		database: db,
	}
}

func (r *Repository) GetUserEmail(ctx context.Context, alertID int64) (string, error) {
	var email string
	err := r.database.QueryRow(ctx, "SELECT email FROM users WHERE alert_id=$1", alertID).Scan(&email)
	if err != nil {
		return "", err
	}
	return email, nil
}
