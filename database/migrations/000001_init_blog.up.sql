CREATE TABLE blog (
    id SERIAL PRIMARY KEY,
    title TEXT NOT NULL UNIQUE,
    content_desription TEXT NOT NULL,
    rss_feed TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
