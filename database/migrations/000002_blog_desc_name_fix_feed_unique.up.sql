ALTER TABLE blog RENAME COLUMN content_desription TO content_description;
ALTER TABLE blog ADD CONSTRAINT rss_feed_unique UNIQUE (rss_feed);
