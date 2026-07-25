package repository

import (
	"database/sql"
	"strings"
	"time"

	"github.com/kiRiLL3311/Currency/news/internal/models"
	"github.com/lib/pq"
)

type NewsRepository struct {
	DB *sql.DB
}

func NewNewsRepository(db *sql.DB) *NewsRepository {
	return &NewsRepository{DB: db}
}

func (r *NewsRepository) ListByRegion(region, keyword string, limit int) ([]models.Article, error) {
	if limit <= 0 {
		limit = 40
	}

	keyword = strings.TrimSpace(keyword)

	var (
		rows *sql.Rows
		err  error
	)

	if keyword == "" {
		rows, err = r.DB.Query(`
			SELECT id, title, summary, url, source, region, currency_tags, published_at, fetched_at
			FROM news_articles
			WHERE region = $1
			ORDER BY published_at DESC
			LIMIT $2
		`, region, limit)
	} else {
		pattern := "%" + keyword + "%"
		rows, err = r.DB.Query(`
			SELECT id, title, summary, url, source, region, currency_tags, published_at, fetched_at
			FROM news_articles
			WHERE region = $1
			  AND (
				title ILIKE $2
				OR summary ILIKE $2
				OR source ILIKE $2
				OR array_to_string(currency_tags, ' ') ILIKE $2
			  )
			ORDER BY published_at DESC
			LIMIT $3
		`, region, pattern, limit)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	articles := make([]models.Article, 0)
	for rows.Next() {
		var a models.Article
		var tags pq.StringArray
		if err := rows.Scan(
			&a.ID,
			&a.Title,
			&a.Summary,
			&a.URL,
			&a.Source,
			&a.Region,
			&tags,
			&a.PublishedAt,
			&a.FetchedAt,
		); err != nil {
			return nil, err
		}
		a.CurrencyTags = []string(tags)
		articles = append(articles, a)
	}

	return articles, nil
}

func (r *NewsRepository) Upsert(region string, article models.ExternalArticle) error {
	published := article.PublishedAt
	if published.IsZero() {
		published = time.Now().UTC()
	}

	tags := article.Tags
	if tags == nil {
		tags = []string{}
	}

	_, err := r.DB.Exec(`
		INSERT INTO news_articles(title, summary, url, source, region, currency_tags, published_at, fetched_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, NOW())
		ON CONFLICT (url) DO UPDATE SET
			title = EXCLUDED.title,
			summary = EXCLUDED.summary,
			source = EXCLUDED.source,
			region = EXCLUDED.region,
			currency_tags = EXCLUDED.currency_tags,
			published_at = EXCLUDED.published_at,
			fetched_at = NOW()
	`,
		article.Title,
		article.Summary,
		article.URL,
		article.Source,
		region,
		pq.Array(tags),
		published,
	)

	return err
}

func (r *NewsRepository) GetUserRegion(userID int) (string, error) {
	var region string
	err := r.DB.QueryRow(`
		SELECT region FROM users WHERE id = $1
	`, userID).Scan(&region)
	if err != nil {
		return "US", err
	}
	if region == "" {
		return "US", nil
	}
	return region, nil
}
