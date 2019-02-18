package main

import (
	"github.com/jzelinskie/geddit"
	"github.com/x1unix/demo-go-plugins/server/feed"
)

func submissionsToPosts(subs []*geddit.Submission) feed.Posts {
	result := make(feed.Posts, 0, len(subs))
	for _, s := range subs {
		result = append(result, feed.Post{
			ID:         s.ID,
			URL:        s.URL,
			Title:      s.Title,
			Text:       s.SelftextHTML,
			ImageURL:   s.ThumbnailURL,
			SourceType: name,
			CreatedAt:  s.DateCreated,
		})
	}

	return result
}
