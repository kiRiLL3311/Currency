package repository

import (
	"database/sql"

	"github.com/kiRiLL3311/Currency/backend/internal/models"
)

type RateRepository struct {
	DB *sql.DB
}

func NewRateRepository(db *sql.DB) *RateRepository {
	return &RateRepository{DB: db}
}

func (r *RateRepository) GetAll() ([]models.Rate, error) {
	rows, err := r.DB.Query(`
		SELECT id, base_currency, target_currency, rate, previous_rate
		FROM rates
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Non-nil empty slice encodes as [] in JSON (nil encodes as null).
	rates := make([]models.Rate, 0)

	for rows.Next() {
		var rate models.Rate
		var previous sql.NullFloat64

		err := rows.Scan(
			&rate.ID,
			&rate.BaseCurrency,
			&rate.TargetCurrency,
			&rate.Rate,
			&previous,
		)
		if err != nil {
			return nil, err
		}

		if previous.Valid {
			prev := previous.Float64
			rate.PreviousRate = &prev
			if prev != 0 {
				change := ((rate.Rate - prev) / prev) * 100
				rate.ChangePercent = &change
			}
		}

		rates = append(rates, rate)
	}

	return rates, nil
}

func (r *RateRepository) GetRate(base, target string) (float64, error) {
	var rate float64

	err := r.DB.QueryRow(`
		SELECT rate
		FROM rates
		WHERE base_currency = $1
		  AND target_currency = $2
	`, base, target).Scan(&rate)

	return rate, err
}

func (r *RateRepository) SaveRate(from, to string, rate float64, previousRate *float64) error {
	var previous interface{}
	if previousRate != nil {
		previous = *previousRate
	}

	_, err := r.DB.Exec(`
		INSERT INTO rates(base_currency, target_currency, rate, previous_rate)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (base_currency, target_currency)
		DO UPDATE SET
			rate = EXCLUDED.rate,
			previous_rate = COALESCE(EXCLUDED.previous_rate, rates.previous_rate)
	`, from, to, rate, previous)

	return err
}
