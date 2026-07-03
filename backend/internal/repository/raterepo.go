package repository

import (
	"Currency/backend/internal/models"
	"database/sql"
)

type RateRepository struct {
	DB *sql.DB
}

func NewRateRepository(db *sql.DB) *RateRepository {
	return &RateRepository{DB: db}
}

func (r *RateRepository) GetAll() ([]models.Rate, error) {
	rows, err := r.DB.Query(`
		SELECT id, base_currency, target_currency, rate
		FROM rates
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rates []models.Rate

	for rows.Next() {
		var rate models.Rate

		err := rows.Scan(
			&rate.ID,
			&rate.BaseCurrency,
			&rate.TargetCurrency,
			&rate.Rate,
		)
		if err != nil {
			return nil, err
		}

		rates = append(rates, rate)
	}

	return rates, nil
}

func (r *RateRepository) GetRate(from, to string) (float64, error) {
	var rate float64

	err := r.DB.QueryRow(
		`SELECT rate
		 FROM rates
		 WHERE base_currency = $1
		   AND target_currency = $2`,
		from,
		to,
	).Scan(&rate)

	return rate, err
}

func (r *RateRepository) SaveRate(from, to string, rate float64) error {
	_, err := r.DB.Exec(`
		INSERT INTO rates(base_currency, target_currency, rate)
		VALUES ($1, $2, $3)
		ON CONFLICT (base_currency, target_currency)
		DO UPDATE SET rate = EXCLUDED.rate
	`, from, to, rate)

	return err
}
