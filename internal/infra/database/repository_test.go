package database

import (
	"context"
	"errors"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRepository_GetUserEmail_Success(t *testing.T) {
	mock, err := pgxmock.NewConn()
	require.NoError(t, err)
	defer mock.Close(context.Background())

	repo := NewRepository(mock)

	alertID := int64(1)
	expectedEmail := "user@example.com"

	rows := pgxmock.NewRows([]string{"email"}).AddRow(expectedEmail)
	mock.ExpectQuery("SELECT email FROM users WHERE alert_id=\\$1").
		WithArgs(alertID).
		WillReturnRows(rows)

	email, err := repo.GetUserEmail(context.Background(), alertID)

	assert.NoError(t, err)
	assert.Equal(t, expectedEmail, email)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_GetUserEmail_NotFound(t *testing.T) {
	mock, err := pgxmock.NewConn()
	require.NoError(t, err)
	defer mock.Close(context.Background())

	repo := NewRepository(mock)

	alertID := int64(999)

	mock.ExpectQuery("SELECT email FROM users WHERE alert_id=\\$1").
		WithArgs(alertID).
		WillReturnError(pgx.ErrNoRows)

	email, err := repo.GetUserEmail(context.Background(), alertID)

	assert.Error(t, err)
	assert.Equal(t, pgx.ErrNoRows, err)
	assert.Empty(t, email)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_GetUserEmail_DatabaseError(t *testing.T) {
	mock, err := pgxmock.NewConn()
	require.NoError(t, err)
	defer mock.Close(context.Background())

	repo := NewRepository(mock)

	alertID := int64(1)
	expectedError := errors.New("database connection error")

	mock.ExpectQuery("SELECT email FROM users WHERE alert_id=\\$1").
		WithArgs(alertID).
		WillReturnError(expectedError)

	email, err := repo.GetUserEmail(context.Background(), alertID)

	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	assert.Empty(t, email)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_GetUserEmail_DifferentAlertIDs(t *testing.T) {
	testCases := []struct {
		name          string
		alertID       int64
		expectedEmail string
	}{
		{
			name:          "Alert ID 1",
			alertID:       1,
			expectedEmail: "user1@example.com",
		},
		{
			name:          "Alert ID 100",
			alertID:       100,
			expectedEmail: "user100@example.com",
		},
		{
			name:          "Alert ID 9999",
			alertID:       9999,
			expectedEmail: "user9999@example.com",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mock, err := pgxmock.NewConn()
			require.NoError(t, err)
			defer mock.Close(context.Background())

			repo := NewRepository(mock)

			rows := pgxmock.NewRows([]string{"email"}).AddRow(tc.expectedEmail)
			mock.ExpectQuery("SELECT email FROM users WHERE alert_id=\\$1").
				WithArgs(tc.alertID).
				WillReturnRows(rows)

			email, err := repo.GetUserEmail(context.Background(), tc.alertID)

			assert.NoError(t, err)
			assert.Equal(t, tc.expectedEmail, email)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestRepository_NewRepository(t *testing.T) {
	mock, err := pgxmock.NewConn()
	require.NoError(t, err)
	defer mock.Close(context.Background())

	repo := NewRepository(mock)

	assert.NotNil(t, repo)
	assert.NotNil(t, repo.database)
}

func TestRepository_GetUserEmail_ContextCancellation(t *testing.T) {
	mock, err := pgxmock.NewConn()
	require.NoError(t, err)
	defer mock.Close(context.Background())

	repo := NewRepository(mock)

	alertID := int64(1)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	mock.ExpectQuery("SELECT email FROM users WHERE alert_id=\\$1").
		WithArgs(alertID).
		WillReturnError(context.Canceled)

	email, err := repo.GetUserEmail(ctx, alertID)

	assert.Error(t, err)
	assert.Empty(t, email)
	assert.NoError(t, mock.ExpectationsWereMet())
}
