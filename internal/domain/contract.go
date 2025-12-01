package domain

import (
	"context"
	"time"
)

type LinkGenerator interface {
	Generate(origin, dest string, out, ret time.Time) string
}

type TempEmailSender interface {
	Send(to, subject, body string) error
}

type AlertRepository interface {
	GetUserEmail(ctx context.Context, alertID int64) (string, error)
}
