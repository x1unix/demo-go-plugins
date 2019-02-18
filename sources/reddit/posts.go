package main

import (
	"github.com/jzelinskie/geddit"
	"github.com/x1unix/demo-go-plugins/server/feed"
)

func submissionsToPosts(subs []*geddit.Submission) feed.Posts {
	result := make(feed.Posts, 0, len(subs))
	for _, s := range subs {
		post := feed.Post{
			ID:         s.ID,
			URL:        s.URL,
			Title:      s.Title,
			Text:       s.SelftextHTML,
			SourceType: name,
			CreatedAt:  s.DateCreated,
		}

		if s.ThumbnailURL != "self" && s.ThumbnailURL != "default" {
			post.ImageURL = s.ThumbnailURL
		}

		result = append(result, post)
	}

	return result
}
