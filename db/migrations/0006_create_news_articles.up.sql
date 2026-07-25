CREATE TABLE news_articles (
    id SERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    summary TEXT NOT NULL DEFAULT '',
    url TEXT NOT NULL,
    source VARCHAR(200) NOT NULL DEFAULT '',
    region VARCHAR(10) NOT NULL,
    currency_tags TEXT[] NOT NULL DEFAULT '{}',
    published_at TIMESTAMP NOT NULL DEFAULT NOW(),
    fetched_at TIMESTAMP NOT NULL DEFAULT NOW(),
    CONSTRAINT news_articles_url_unique UNIQUE (url)
);

CREATE INDEX news_articles_region_published_idx
    ON news_articles (region, published_at DESC);
