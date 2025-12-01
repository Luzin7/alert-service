package consumer

import (
	"fmt"
	"time"

	"github.com/Luzin7/alert-service/internal/domain"
)

type PriceUpdatedPayload struct {
	MessageID    string    `json:"messageId"`
	AlertID      int64     `json:"alertId"`
	Origin       string    `json:"origin"`
	Destination  string    `json:"destination"`
	OutboundDate string    `json:"outboundDate"`
	ReturnDate   string    `json:"returnDate"`
	OldPrice     float64   `json:"oldPrice"`
	NewPrice     float64   `json:"newPrice"`
	Currency     string    `json:"currency"`
	TargetPrice  float64   `json:"targetPrice"`
	ToleranceUp  float64   `json:"toleranceUp"`
	CheckedAt    time.Time `json:"checkedAt"`
}

func (p *PriceUpdatedPayload) ToDomain() (*domain.Alert, error) {
	out, err := time.Parse("2006-01-02", p.OutboundDate)
	if err != nil {
		return nil, fmt.Errorf("data ida invalida: %w", err)
	}
	ret, err := time.Parse("2006-01-02", p.ReturnDate)
	if err != nil {
		return nil, fmt.Errorf("data volta invalida: %w", err)
	}

	return &domain.Alert{
		MessageID:    p.MessageID,
		ID:           p.AlertID,
		Origin:       p.Origin,
		Destination:  p.Destination,
		OutboundDate: out,
		ReturnDate:   ret,
		OldPrice:     p.OldPrice,
		NewPrice:     p.NewPrice,
		Currency:     p.Currency,
		TargetPrice:  p.TargetPrice,
		CheckedAt:    p.CheckedAt,
	}, nil
}
