package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/kiRiLL3311/Currency/news/internal/models"
	"github.com/mmcdole/gofeed"
)

type RegionQuery struct {
	Query string
	Tags  []string
	Feeds []string
}

var RegionQueries = map[string]RegionQuery{
	"US": {
		Query: "forex OR FX OR dollar OR Federal Reserve OR currency markets",
		Tags:  []string{"USD", "FX"},
		Feeds: []string{
			"https://feeds.bbci.co.uk/news/business/rss.xml",
			"https://www.cnbc.com/id/10000664/device/rss/rss.html",
		},
	},
	"EU": {
		Query: "forex OR euro OR ECB OR European Central Bank OR currency markets",
		Tags:  []string{"EUR", "FX"},
		Feeds: []string{
			"https://feeds.bbci.co.uk/news/business/rss.xml",
			"https://www.ecb.europa.eu/rss/press.html",
		},
	},
	"GB": {
		Query: "forex OR pound OR sterling OR Bank of England OR currency markets",
		Tags:  []string{"GBP", "FX"},
		Feeds: []string{
			"https://feeds.bbci.co.uk/news/business/rss.xml",
		},
	},
	"JP": {
		Query: "forex OR yen OR Bank of Japan OR currency markets",
		Tags:  []string{"JPY", "FX"},
		Feeds: []string{
			"https://feeds.bbci.co.uk/news/business/rss.xml",
		},
	},
	"AU": {
		Query: "forex OR Australian dollar OR RBA OR currency markets",
		Tags:  []string{"AUD", "FX"},
		Feeds: []string{
			"https://feeds.bbci.co.uk/news/business/rss.xml",
		},
	},
}

type newsAPIResponse struct {
	Status   string `json:"status"`
	Articles []struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		URL         string `json:"url"`
		PublishedAt string `json:"publishedAt"`
		Source      struct {
			Name string `json:"name"`
		} `json:"source"`
	} `json:"articles"`
}

func FetchForRegion(region, apiKey, keyword string) ([]models.ExternalArticle, error) {
	cfg, ok := RegionQueries[region]
	if !ok {
		cfg = RegionQueries["US"]
	}

	keyword = strings.TrimSpace(keyword)
	cfg.Query = buildQuery(cfg.Query, keyword)

	if strings.TrimSpace(apiKey) != "" {
		articles, err := fetchNewsAPI(cfg, apiKey)
		if err == nil && len(articles) > 0 {
			return filterByKeyword(articles, keyword), nil
		}
	}

	articles, err := fetchRSS(cfg)
	if err != nil {
		return nil, err
	}
	return filterByKeyword(articles, keyword), nil
}

func buildQuery(baseQuery, keyword string) string {
	keyword = strings.TrimSpace(keyword)
	if keyword == "" {
		return baseQuery
	}
	// Keyword leads the pull; region FX terms keep results market-relevant.
	return fmt.Sprintf("(%s) AND (%s)", keyword, baseQuery)
}

func filterByKeyword(articles []models.ExternalArticle, keyword string) []models.ExternalArticle {
	keyword = strings.TrimSpace(strings.ToLower(keyword))
	if keyword == "" {
		return articles
	}

	out := make([]models.ExternalArticle, 0, len(articles))
	for _, a := range articles {
		haystack := strings.ToLower(strings.Join([]string{
			a.Title,
			a.Summary,
			a.Source,
			strings.Join(a.Tags, " "),
		}, " "))
		if strings.Contains(haystack, keyword) {
			out = append(out, a)
		}
	}
	return out
}

func fetchNewsAPI(cfg RegionQuery, apiKey string) ([]models.ExternalArticle, error) {
	endpoint := fmt.Sprintf(
		"https://newsapi.org/v2/everything?q=%s&language=en&sortBy=publishedAt&pageSize=30&apiKey=%s",
		url.QueryEscape(cfg.Query),
		url.QueryEscape(apiKey),
	)

	resp, err := http.Get(endpoint)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("newsapi status %d", resp.StatusCode)
	}

	var payload newsAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return nil, err
	}

	out := make([]models.ExternalArticle, 0, len(payload.Articles))
	for _, a := range payload.Articles {
		if a.Title == "" || a.URL == "" {
			continue
		}
		published, _ := time.Parse(time.RFC3339, a.PublishedAt)
		out = append(out, models.ExternalArticle{
			Title:       a.Title,
			Summary:     a.Description,
			URL:         a.URL,
			Source:      a.Source.Name,
			PublishedAt: published,
			Tags:        cfg.Tags,
		})
	}
	return out, nil
}

func fetchRSS(cfg RegionQuery) ([]models.ExternalArticle, error) {
	fp := gofeed.NewParser()
	out := make([]models.ExternalArticle, 0)

	for _, feedURL := range cfg.Feeds {
		feed, err := fp.ParseURL(feedURL)
		if err != nil {
			continue
		}
		for _, item := range feed.Items {
			if item.Title == "" || item.Link == "" {
				continue
			}
			published := time.Now().UTC()
			if item.PublishedParsed != nil {
				published = *item.PublishedParsed
			}
			summary := item.Description
			if summary == "" {
				summary = item.Content
			}
			source := feed.Title
			out = append(out, models.ExternalArticle{
				Title:       item.Title,
				Summary:     stripHTML(summary),
				URL:         item.Link,
				Source:      source,
				PublishedAt: published,
				Tags:        cfg.Tags,
			})
		}
	}

	if len(out) == 0 {
		return nil, fmt.Errorf("no articles fetched for feeds")
	}
	return out, nil
}

func stripHTML(s string) string {
	var b strings.Builder
	inTag := false
	for _, r := range s {
		switch {
		case r == '<':
			inTag = true
		case r == '>':
			inTag = false
		case !inTag:
			b.WriteRune(r)
		}
	}
	return strings.TrimSpace(b.String())
}
