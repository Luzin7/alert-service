package consumer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHandler_Handle_InvalidJSON(t *testing.T) {
	handler := &Handler{}

	invalidJSON := []byte(`{"invalid json}`)

	err := handler.Handle(invalidJSON)

	assert.Error(t, err)
}

func TestHandler_Handle_InvalidDate(t *testing.T) {
	handler := &Handler{}

	invalidDateJSON := []byte(`{
		"messageId": "msg-123",
		"alertId": 1,
		"origin": "GRU",
		"destination": "JFK",
		"outboundDate": "invalid-date",
		"returnDate": "2025-12-20",
		"oldPrice": 1500.00,
		"newPrice": 1200.00,
		"currency": "BRL",
		"targetPrice": 1000.00,
		"checkedAt": "2025-12-02T10:00:00Z"
	}`)

	err := handler.Handle(invalidDateJSON)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "data ida invalida")
}
