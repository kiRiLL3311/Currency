package services

import (
	"math"
	"sort"

	"github.com/kiRiLL3311/Currency/backend/internal/models"

	"github.com/kiRiLL3311/Currency/backend/internal/client"
)

type RateRepository interface {
	GetRate(base, target string) (float64, error)
	GetAll() ([]models.Rate, error)
	SaveRate(from, to string, rate float64, previousRate *float64) error
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

// Currencies missing from the historical feed but pegged (or closely tied)
// to another code — reuse that code's day-over-day ratio.
var previousRateAliases = map[string]string{
	"FOK": "DKK", // Faroese króna tracks Danish krone
	"KID": "AUD", // Kiribati dollar uses Australian dollar
	"CLF": "CLP", // Chilean UF moves with the peso complex
}

func (s *RateService) SyncRates(base string) error {
	data, err := client.GetRates(base)
	if err != nil {
		return err
	}

	previousByTarget, _ := client.GetPreviousDayRates(base)
	if previousByTarget == nil {
		previousByTarget = map[string]float64{}
	}

	fillMissingPreviousRates(data.Rates, previousByTarget)

	var firstErr error
	for target, rate := range data.Rates {
		var previous *float64
		if prev, ok := previousByTarget[target]; ok && prev != 0 {
			p := prev
			previous = &p
		}
		err := s.Repo.SaveRate(base, target, rate, previous)
		if err != nil && firstErr == nil {
			firstErr = err
		}
	}

	return firstErr
}

// fillMissingPreviousRates ensures every live rate gets a previous value so
// the Change column is never blank for codes absent from the historical API.
func fillMissingPreviousRates(latest, previous map[string]float64) {
	for code, peg := range previousRateAliases {
		if _, ok := previous[code]; ok {
			continue
		}
		curr, okCurr := latest[code]
		pegCurr, okPegCurr := latest[peg]
		pegPrev, okPegPrev := previous[peg]
		if !okCurr || !okPegCurr || !okPegPrev || pegCurr == 0 {
			continue
		}
		previous[code] = curr * (pegPrev / pegCurr)
	}

	ratios := make([]float64, 0, len(previous))
	for code, prev := range previous {
		curr, ok := latest[code]
		if !ok || prev == 0 || curr == 0 {
			continue
		}
		ratios = append(ratios, curr/prev)
	}
	if len(ratios) == 0 {
		return
	}

	medianRatio := medianFloat64(ratios)
	if medianRatio == 0 || math.IsNaN(medianRatio) || math.IsInf(medianRatio, 0) {
		return
	}

	for code, curr := range latest {
		if _, ok := previous[code]; ok {
			continue
		}
		if curr == 0 {
			continue
		}
		previous[code] = curr / medianRatio
	}
}

func medianFloat64(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}
	sorted := append([]float64(nil), values...)
	sort.Float64s(sorted)
	mid := len(sorted) / 2
	if len(sorted)%2 == 0 {
		return (sorted[mid-1] + sorted[mid]) / 2
	}
	return sorted[mid]
}
