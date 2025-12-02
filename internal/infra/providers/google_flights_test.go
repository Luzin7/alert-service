package providers

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGoogleFlightsGenerator_Generate(t *testing.T) {
	generator := GoogleFlightsGenerator{
		BaseURL: "https://www.google.com/travel/flights",
	}

	origin := "GRU"
	destination := "JFK"
	outbound := time.Date(2025, 12, 15, 0, 0, 0, 0, time.UTC)
	returnDate := time.Date(2025, 12, 20, 0, 0, 0, 0, time.UTC)

	link := generator.Generate(origin, destination, outbound, returnDate)

	assert.NotEmpty(t, link)
	assert.Contains(t, link, "https://www.google.com/travel/flights")
	assert.Contains(t, link, destination)
	assert.Contains(t, link, origin)
}

func TestGoogleFlightsGenerator_Generate_DifferentCities(t *testing.T) {
	generator := GoogleFlightsGenerator{
		BaseURL: "https://www.google.com/travel/flights",
	}

	testCases := []struct {
		name        string
		origin      string
		destination string
	}{
		{
			name:        "São Paulo to Miami",
			origin:      "GRU",
			destination: "MIA",
		},
		{
			name:        "Rio to Paris",
			origin:      "GIG",
			destination: "CDG",
		},
		{
			name:        "Brasília to London",
			origin:      "BSB",
			destination: "LHR",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			outbound := time.Date(2025, 10, 1, 0, 0, 0, 0, time.UTC)
			returnDate := time.Date(2025, 10, 10, 0, 0, 0, 0, time.UTC)

			link := generator.Generate(tc.origin, tc.destination, outbound, returnDate)

			assert.NotEmpty(t, link)
			assert.Contains(t, link, tc.destination)
			assert.Contains(t, link, tc.origin)
		})
	}
}

func TestGoogleFlightsGenerator_Generate_URLFormat(t *testing.T) {
	generator := GoogleFlightsGenerator{
		BaseURL: "https://www.google.com/travel/flights",
	}

	origin := "GRU"
	destination := "JFK"
	outbound := time.Date(2025, 12, 15, 0, 0, 0, 0, time.UTC)
	returnDate := time.Date(2025, 12, 20, 0, 0, 0, 0, time.UTC)

	link := generator.Generate(origin, destination, outbound, returnDate)

	assert.True(t, len(link) > 0, "Link should not be empty")
	assert.Contains(t, link, "https://", "Link should use HTTPS protocol")
}

func TestGoogleFlightsGenerator_Generate_EmptyBaseURL(t *testing.T) {
	generator := GoogleFlightsGenerator{
		BaseURL: "",
	}

	origin := "GRU"
	destination := "JFK"
	outbound := time.Date(2025, 12, 15, 0, 0, 0, 0, time.UTC)
	returnDate := time.Date(2025, 12, 20, 0, 0, 0, 0, time.UTC)

	link := generator.Generate(origin, destination, outbound, returnDate)

	assert.NotEmpty(t, link)
	assert.Contains(t, link, origin)
	assert.Contains(t, link, destination)
}
