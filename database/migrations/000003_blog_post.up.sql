CREATE TABLE blog_post (
    id SERIAL PRIMARY KEY,
    blog_id INT REFERENCES blog(id),
    title TEXT NOT NULL UNIQUE,
    content TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);