package main

import (
	"database/sql"
	"log/slog"
	"net/http"
	"net/url"
	"strings"
)

func addBlog(db *sql.DB) http.HandlerFunc {
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

		if _, err := db.Exec(
			`INSERT INTO blog (title, content_description, rss_feed) 
			VALUES ($1, $2, $3)`,
			title, description, rssURL.String(),
		); err != nil {
			if strings.Contains(err.Error(), "violates unique") {
				http.Error(w, "blog title or rss feed not unique, already added?", http.StatusBadRequest)
				slog.Error("Failed uniqueness check", "error", err)
				return
			}
			http.Error(w, "failed", http.StatusInternalServerError)
			slog.Error("Failed inserting blog", "error", err)
			return
		}

		// TODO: redirect to home page with success message and prompt
		// user to check back later (mention import cycle time)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("success"))
	}
}
