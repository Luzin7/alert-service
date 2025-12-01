package database

import (
	"context"

	"github.com/jackc/pgx/v5"
)

type Repository struct {
	database *pgx.Conn
}

func NewRepository(db *pgx.Conn) *Repository {
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
