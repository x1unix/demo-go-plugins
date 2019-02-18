package main

import (
	"github.com/jzelinskie/geddit"
	"github.com/x1unix/demo-go-plugins/server/feed"
	"html"
	"strings"
)

func submissionsToPosts(subs []*geddit.Submission) feed.Posts {
	result := make(feed.Posts, 0, len(subs))
	for _, s := range subs {
		post := feed.Post{
			ID:         s.ID,
			URL:        s.URL,
			Title:      s.Title,
			Text:       unescapeContent(s.SelftextHTML),
			SourceType: name,
			CreatedAt:  s.DateCreated,
		}

		if s.ThumbnailURL != "self" && s.ThumbnailURL != "default" && s.ThumbnailURL != "nsfw" {
			post.ImageURL = s.ThumbnailURL
		}

		result = append(result, post)
	}

	return result
}

func unescapeContent(content string) string {
	r := html.UnescapeString(content)
	r = strings.Replace(r, `\"`, `"`, -1)
	return r
}
