CREATE TABLE rss_feed (
    id SERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    feed_url TEXT NOT NULL UNIQUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
