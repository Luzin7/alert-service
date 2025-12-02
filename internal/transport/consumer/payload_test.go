package consumer

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPriceUpdatedPayload_ToDomain_Success(t *testing.T) {
	payload := &PriceUpdatedPayload{
		MessageID:    "msg-123",
		AlertID:      1,
		Origin:       "GRU",
		Destination:  "JFK",
		OutboundDate: "2025-12-15",
		ReturnDate:   "2025-12-20",
		OldPrice:     1500.00,
		NewPrice:     1200.00,
		Currency:     "BRL",
		TargetPrice:  1000.00,
		ToleranceUp:  100.00,
		CheckedAt:    time.Now(),
	}

	alert, err := payload.ToDomain()

	require.NoError(t, err)
	assert.Equal(t, "msg-123", alert.MessageID)
	assert.Equal(t, int64(1), alert.ID)
	assert.Equal(t, "GRU", alert.Origin)
	assert.Equal(t, "JFK", alert.Destination)
	assert.Equal(t, 1500.00, alert.OldPrice)
	assert.Equal(t, 1200.00, alert.NewPrice)
	assert.Equal(t, "BRL", alert.Currency)
	assert.Equal(t, 1000.00, alert.TargetPrice)

	expectedOut, _ := time.Parse("2006-01-02", "2025-12-15")
	expectedRet, _ := time.Parse("2006-01-02", "2025-12-20")
	assert.Equal(t, expectedOut, alert.OutboundDate)
	assert.Equal(t, expectedRet, alert.ReturnDate)
}

func TestPriceUpdatedPayload_ToDomain_InvalidOutboundDate(t *testing.T) {
	payload := &PriceUpdatedPayload{
		MessageID:    "msg-123",
		AlertID:      1,
		Origin:       "GRU",
		Destination:  "JFK",
		OutboundDate: "invalid-date",
		ReturnDate:   "2025-12-20",
		OldPrice:     1500.00,
		NewPrice:     1200.00,
		Currency:     "BRL",
		TargetPrice:  1000.00,
		CheckedAt:    time.Now(),
	}

	alert, err := payload.ToDomain()

	assert.Error(t, err)
	assert.Nil(t, alert)
	assert.Contains(t, err.Error(), "data ida invalida")
}

func TestPriceUpdatedPayload_ToDomain_InvalidReturnDate(t *testing.T) {
	payload := &PriceUpdatedPayload{
		MessageID:    "msg-123",
		AlertID:      1,
		Origin:       "GRU",
		Destination:  "JFK",
		OutboundDate: "2025-12-15",
		ReturnDate:   "invalid-date",
		OldPrice:     1500.00,
		NewPrice:     1200.00,
		Currency:     "BRL",
		TargetPrice:  1000.00,
		CheckedAt:    time.Now(),
	}

	alert, err := payload.ToDomain()

	assert.Error(t, err)
	assert.Nil(t, alert)
	assert.Contains(t, err.Error(), "data volta invalida")
}

func TestPriceUpdatedPayload_ToDomain_WrongDateFormat(t *testing.T) {
	payload := &PriceUpdatedPayload{
		MessageID:    "msg-123",
		AlertID:      1,
		Origin:       "GRU",
		Destination:  "JFK",
		OutboundDate: "15/12/2025",
		ReturnDate:   "20/12/2025",
		OldPrice:     1500.00,
		NewPrice:     1200.00,
		Currency:     "BRL",
		TargetPrice:  1000.00,
		CheckedAt:    time.Now(),
	}

	alert, err := payload.ToDomain()

	assert.Error(t, err)
	assert.Nil(t, alert)
}
