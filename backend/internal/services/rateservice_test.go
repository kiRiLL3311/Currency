package services

import (
	"fmt"
	"testing"

	"github.com/kiRiLL3311/Currency/backend/internal/models"
)

type MockRepo struct{}

func (m *MockRepo) GetRate(base, target string) (float64, error) {
	rates := map[string]float64{
		"EUR": 0.85,
		"GBP": 0.75,
		"JPY": 145,
	}

	rate, ok := rates[target]
	if !ok {
		return 0, fmt.Errorf("rate not found")
	}

	return rate, nil
}

func (m *MockRepo) GetAll() ([]models.Rate, error) {
	return []models.Rate{}, nil
}

func (m *MockRepo) SaveRate(from, to string, rate float64, previousRate *float64) error {
	return nil
}

func TestConvert(t *testing.T) {
	repo := &MockRepo{}
	service := NewRateService(repo)

	tests := []struct {
		name          string
		from          string
		to            string
		amount        float64
		expectedRate  float64
		expectedValue float64
		wantErr       bool
	}{
		{
			name:          "USD to EUR",
			from:          "USD",
			to:            "EUR",
			amount:        100,
			expectedRate:  0.85,
			expectedValue: 85,
		},
		{
			name:          "EUR to USD",
			from:          "EUR",
			to:            "USD",
			amount:        85,
			expectedRate:  1 / 0.85,
			expectedValue: 100,
		},
		{
			name:          "EUR to GBP",
			from:          "EUR",
			to:            "GBP",
			amount:        100,
			expectedRate:  0.75 / 0.85,
			expectedValue: 100 * (0.75 / 0.85),
		},
		{
			name:          "Same currency",
			from:          "USD",
			to:            "USD",
			amount:        50,
			expectedRate:  1,
			expectedValue: 50,
		},

		{
			name:    "Unknown currency",
			from:    "USD",
			to:      "ABC",
			amount:  100,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			value, rate, err := service.Convert(tt.from, tt.to, tt.amount)

			// Handle cases where an error is expected.
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected an error, got nil")
				}
				return
			}

			// For all other cases, any error is unexpected.
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			const epsilon = 0.000001

			if abs(rate-tt.expectedRate) > epsilon {
				t.Errorf("expected rate %f, got %f", tt.expectedRate, rate)
			}

			if abs(value-tt.expectedValue) > epsilon {
				t.Errorf("expected value %f, got %f", tt.expectedValue, value)
			}
		})
	}
}

func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}
