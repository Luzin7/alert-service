package domain

import "time"

type Alert struct {
	ID           int64
	MessageID    string
	Origin       string
	Destination  string
	OutboundDate time.Time
	ReturnDate   time.Time
	NewPrice     float64
	OldPrice     float64
	TargetPrice  float64
	Currency     string
	CheckedAt    time.Time
	Link         string
}
