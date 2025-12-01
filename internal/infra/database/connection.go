package database

import (
	"context"

	"github.com/jackc/pgx/v5"
)

func DatabaseConnection(connectionString string, databaseName string) (*pgx.Conn, error) {
	connection, err := pgx.Connect(context.Background(), connectionString)
	if err != nil {
		return nil, err
	}

	if err := connection.Ping(context.Background()); err != nil {
		return nil, err
	}

	return connection, nil
}

func CloseDatabaseConnection(connection *pgx.Conn) error {
	return connection.Close(context.Background())
}
