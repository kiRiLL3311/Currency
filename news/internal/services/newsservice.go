package services

import (
	"log"
	"strings"

	"github.com/kiRiLL3311/Currency/news/internal/client"
	"github.com/kiRiLL3311/Currency/news/internal/models"
	"github.com/kiRiLL3311/Currency/news/internal/repository"
)

type NewsService struct {
	Repo   *repository.NewsRepository
	APIKey string
}

func NewNewsService(repo *repository.NewsRepository, apiKey string) *NewsService {
	return &NewsService{Repo: repo, APIKey: apiKey}
}

func (s *NewsService) List(userID int, region, keyword string) ([]models.Article, error) {
	region = strings.TrimSpace(strings.ToUpper(region))
	if region == "" {
		userRegion, err := s.Repo.GetUserRegion(userID)
		if err == nil {
			region = userRegion
		} else {
			region = "US"
		}
	}

	if _, ok := client.RegionQueries[region]; !ok {
		// Legacy Cyprus region rolls into Europe.
		if region == "CY" {
			region = "EU"
		} else {
			region = "US"
		}
	}

	return s.Repo.ListByRegion(region, strings.TrimSpace(keyword), 40)
}

func (s *NewsService) SyncRegion(region, keyword string) error {
	region = strings.TrimSpace(strings.ToUpper(region))
	if region == "" {
		region = "US"
	}
	if region == "CY" {
		region = "EU"
	}

	articles, err := client.FetchForRegion(region, s.APIKey, strings.TrimSpace(keyword))
	if err != nil {
		return err
	}

	var firstErr error
	for _, article := range articles {
		if err := s.Repo.Upsert(region, article); err != nil && firstErr == nil {
			firstErr = err
		}
	}
	return firstErr
}

func (s *NewsService) SyncAll() {
	for region := range client.RegionQueries {
		if err := s.SyncRegion(region, ""); err != nil {
			log.Printf("news sync %s failed: %v", region, err)
		} else {
			log.Printf("news sync %s completed", region)
		}
	}
}
