-- name: CreateBlog :one
INSERT INTO blog (
  title, content_description, rss_feed
) VALUES (
  $1, $2, $3
)
RETURNING *;

-- name: ListBlogs :many
SELECT * FROM blog
ORDER BY created_at;

-- name: CreateBlogPost :one
INSERT INTO blog_post (
  blog_id, title, content
) VALUES (
  $1, $2, $3
)
RETURNING *;

-- name: ListBlogPosts :many
SELECT * FROM blog_post
WHERE blog_id = $1
ORDER BY created_at DESC;
