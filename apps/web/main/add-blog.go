package main

import (
	"context"
	"log/slog"
	"net/http"
	"net/url"
	"strings"

	"github.com/samiam2013/cadencereader/database"
)

func addBlog(db *database.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO: is this necessary?
		if err := r.ParseForm(); err != nil {
			http.Error(w, "failed to parse form", http.StatusInternalServerError)
			slog.Error("Failed to parse form", "error", err)
			return
		}

		title := r.FormValue("title")
		if strings.TrimSpace(title) == "" {
			http.Error(w, "title empty", http.StatusBadRequest)
			slog.Error("Title empty when adding blog")
			return
		}

		description := r.FormValue("description")
		if strings.TrimSpace(description) == "" {
			http.Error(w, "description empty", http.StatusBadRequest)
			slog.Error("Description empty when adding blog")
			return
		}

		rssURLstr := r.FormValue("rss_url")
		rssURL, err := url.Parse(rssURLstr)
		if err != nil {
			http.Error(w, "invalid RSS URL", http.StatusBadRequest)
			slog.Error("Failed to parse rss url", "error", err)
			return
		}

		ctx := context.Background()
		blog, err := db.CreateBlog(ctx,
			database.CreateBlogParams{
				Title:              title,
				ContentDescription: description,
				RssFeed:            rssURL.String(),
			})
		if err != nil {
			http.Error(w, "failed", http.StatusInternalServerError)
			slog.Error("Failed to create blog row", "error", err)
			return
		}
		slog.Info("Created new blog row", "blog", blog)

		// TODO: redirect to home page with success message and prompt
		// user to check back later (mention import cycle time)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("success"))
	}
}
