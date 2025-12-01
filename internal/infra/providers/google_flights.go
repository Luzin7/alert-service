package providers

import (
	"fmt"
	"time"
)

type GoogleFlightsGenerator struct {
	BaseURL string
}

func (g GoogleFlightsGenerator) Generate(origin, dest string, out, ret time.Time) string {
	return fmt.Sprintf("https://www.google.com/travel/flights?q=Flights%%20to%%20%s%%20from%%20%s...", dest, origin)
}
