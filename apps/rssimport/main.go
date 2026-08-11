package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"slices"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/jackc/pgx/v4"
	"github.com/mmcdole/gofeed"
	"github.com/samiam2013/cadencereader/config"
	"github.com/samiam2013/cadencereader/database"
)

const AuthorizationTitle = "Permission for CadenceReader Syndication"

func main() {
	slog.Info("starting importer binary")
	cfg, err := config.Load()
	if err != nil {
		slog.Error("Failed to load configuration", "error", err)
		os.Exit(1)
	}

	ctx := context.Background()
	conn, err := pgx.Connect(ctx, cfg.DatabaseURL.String())
	if err != nil {
		slog.Error("Failed to connect via pgx", "error", err)
		os.Exit(1)
	}
	db := database.New(conn)

	blogs, err := db.ListBlogs(ctx)
	if err != nil {
		slog.Error("Failed to list blogs for import", "error", err)
	}
	for _, blog := range blogs {
		slog.Info("starting import", "blog", blog)
		fp := gofeed.NewParser()
		feed, err := fp.ParseURL(blog.RssFeed)
		if err != nil {
			slog.Error("Failed to parse blog rss feed", "error", err)
			continue
		}
		slog.Info("feed parsed", "title", feed.Title, "type", feed.FeedType, "len", feed.Len())

		// get existing blog post titles
		existingPostTitles := []string{}
		posts, err := db.ListBlogPosts(ctx, sql.NullInt32{Valid: true, Int32: blog.ID})
		if err != nil {
			slog.Error("Failed to list existing blog posts", "error", err)
			continue
		}
		for _, post := range posts {
			existingPostTitles = append(existingPostTitles, post.Title)
		}

		authorized := false
		links := make(map[string]string)
		for _, feedItem := range feed.Items {
			slog.Info("feed item", "title", feedItem.Title, "link", feedItem.Link)
			fiTitle := strings.ToLower(strings.TrimSpace(feedItem.Title))
			if fiTitle == strings.ToLower(AuthorizationTitle) {
				authorized = true
			}
			// only add new titles
			if !slices.Contains(existingPostTitles, feedItem.Title) {
				links[feedItem.Title] = feedItem.Link
			}
		}

		if authorized {
			for title, link := range links {

				// get the content between <article> tags
				content, err := getArticle(link)
				if err != nil {
					slog.Error("Failed to get content from post link", "error", err)
					continue
				}

				newPost, err := db.CreateBlogPost(ctx,
					database.CreateBlogPostParams{
						BlogID:  sql.NullInt32{Valid: true, Int32: blog.ID},
						Title:   title,
						Content: content,
					})
				if err != nil {
					slog.Error("Failed creating new blog post", "error", err)
					continue
				}
				slog.Info("Created new blog post", "post", newPost)
			}
		}
	}
}

var ErrNoArticleTags = errors.New("no article tags found")

func getArticle(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("failed getting url for article: %w", err)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error parsing document: %w", err)
	}
	text := ""
	doc.Find("article").Each(func(i int, articleSel *goquery.Selection) {
		text = articleSel.Text()
	})
	if text != "" {
		return text, nil
	}
	return "", ErrNoArticleTags
}
