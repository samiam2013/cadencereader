-- name: CreateBlog :one
INSERT INTO blog (
  title, content_description, rss_feed
) VALUES (
  $1, $2, $3
)
RETURNING *;