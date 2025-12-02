package usecases

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Luzin7/alert-service/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type MockLinkGenerator struct {
	mock.Mock
}

func (m *MockLinkGenerator) Generate(origin, dest string, out, ret time.Time) string {
	args := m.Called(origin, dest, out, ret)
	return args.String(0)
}

type MockAlertRepository struct {
	mock.Mock
}

func (m *MockAlertRepository) GetUserEmail(ctx context.Context, alertID int64) (string, error) {
	args := m.Called(ctx, alertID)
	return args.String(0), args.Error(1)
}

type MockEmailSender struct {
	mock.Mock
}

func (m *MockEmailSender) Send(to, subject, body string) error {
	args := m.Called(to, subject, body)
	return args.Error(0)
}

func TestProcessAlert_Execute_LinkGeneration(t *testing.T) {
	mockLinkGen := new(MockLinkGenerator)

	alert := &domain.Alert{
		ID:           1,
		Origin:       "GRU",
		Destination:  "JFK",
		OutboundDate: time.Date(2025, 12, 15, 0, 0, 0, 0, time.UTC),
		ReturnDate:   time.Date(2025, 12, 20, 0, 0, 0, 0, time.UTC),
		NewPrice:     1200.00,
		OldPrice:     1500.00,
		Currency:     "BRL",
	}

	expectedLink := "https://www.google.com/travel/flights?q=Flights%20to%20JFK%20from%20GRU..."
	mockLinkGen.On("Generate", "GRU", "JFK", alert.OutboundDate, alert.ReturnDate).Return(expectedLink)

	link := mockLinkGen.Generate(alert.Origin, alert.Destination, alert.OutboundDate, alert.ReturnDate)

	assert.Equal(t, expectedLink, link)
	assert.NotEmpty(t, link)
	mockLinkGen.AssertExpectations(t)
}

func TestProcessAlert_Execute_RepositoryError(t *testing.T) {
	mockLinkGen := new(MockLinkGenerator)
	mockRepo := new(MockAlertRepository)

	useCase := NewProcessAlert(mockLinkGen, mockRepo, nil)

	alert := &domain.Alert{
		ID:           1,
		Origin:       "GRU",
		Destination:  "JFK",
		OutboundDate: time.Date(2025, 12, 15, 0, 0, 0, 0, time.UTC),
		ReturnDate:   time.Date(2025, 12, 20, 0, 0, 0, 0, time.UTC),
		NewPrice:     1200.00,
		Currency:     "BRL",
	}

	expectedLink := "https://www.google.com/travel/flights?q=Flights%20to%20JFK%20from%20GRU..."
	mockLinkGen.On("Generate", "GRU", "JFK", alert.OutboundDate, alert.ReturnDate).Return(expectedLink)

	expectedError := errors.New("database error")
	mockRepo.On("GetUserEmail", mock.Anything, int64(1)).Return("", expectedError)

	ctx := context.Background()
	err := useCase.Execute(ctx, alert)

	require.Error(t, err)
	assert.Equal(t, expectedError, err)
	mockLinkGen.AssertExpectations(t)
	mockRepo.AssertExpectations(t)
}

func TestProcessAlert_NewProcessAlert(t *testing.T) {
	mockLinkGen := new(MockLinkGenerator)
	mockRepo := new(MockAlertRepository)

	useCase := NewProcessAlert(mockLinkGen, mockRepo, nil)

	assert.NotNil(t, useCase)
	assert.Equal(t, mockLinkGen, useCase.linkGen)
	assert.Equal(t, mockRepo, useCase.repo)
}
