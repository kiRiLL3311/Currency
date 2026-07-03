package services

import (
	"github.com/kiRiLL3311/Currency/backend/internal/models"

	"github.com/kiRiLL3311/Currency/backend/internal/client"
)

type RateRepository interface {
	GetRate(base, target string) (float64, error)
	GetAll() ([]models.Rate, error)
	SaveRate(from, to string, rate float64) error
}

type RateService struct {
	Repo RateRepository
}

func NewRateService(repo RateRepository) *RateService {
	return &RateService{Repo: repo}
}

func (s *RateService) GetRates() ([]models.Rate, error) {
	return s.Repo.GetAll()
}

// func (s *RateService) Convert(from, to string, amount float64) (float64, float64, error) {
// 	rate, err := s.Repo.GetRate(from, to)
// 	if err != nil {
// 		return 0, 0, err
// 	}

//		return rate, amount * rate, nil
//	}
func (s *RateService) Convert(from, to string, amount float64) (float64, float64, error) {
	// API stores all rates relative to USD
	const base = "USD"

	if from == to {
		return amount, 1, nil
	}

	// Base -> X
	if from == base {
		rate, err := s.Repo.GetRate(base, to)
		if err != nil {
			return 0, 0, err
		}

		return amount * rate, rate, nil
	}

	// X -> Base
	if to == base {
		rate, err := s.Repo.GetRate(base, from)
		if err != nil {
			return 0, 0, err
		}

		return amount / rate, 1 / rate, nil
	}

	// X -> Y
	fromRate, err := s.Repo.GetRate(base, from)
	if err != nil {
		return 0, 0, err
	}

	toRate, err := s.Repo.GetRate(base, to)
	if err != nil {
		return 0, 0, err
	}

	crossRate := toRate / fromRate

	return amount * crossRate, crossRate, nil
}
func (s *RateService) SyncRates(base string) error {
	data, err := client.GetRates(base)
	if err != nil {
		return err
	}

	for target, rate := range data.Rates {
		err := s.Repo.SaveRate(base, target, rate)
		if err != nil {
			return err
		}
	}

	return nil
}
