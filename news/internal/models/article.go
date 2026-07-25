package models

import "time"

type Article struct {
	ID           int       `json:"id"`
	Title        string    `json:"title"`
	Summary      string    `json:"summary"`
	URL          string    `json:"url"`
	Source       string    `json:"source"`
	Region       string    `json:"region"`
	CurrencyTags []string  `json:"currency_tags"`
	PublishedAt  time.Time `json:"published_at"`
	FetchedAt    time.Time `json:"fetched_at"`
}

type ExternalArticle struct {
	Title       string
	Summary     string
	URL         string
	Source      string
	PublishedAt time.Time
	Tags        []string
}
