package extension

import (
	"fmt"

	"github.com/jzelinskie/geddit"
	"github.com/x1unix/demo-go-plugins/server/feed"
)

const name = "reddit"

type DataSource struct {
	cfg    Config
	client *geddit.Session
}

func (s *DataSource) Name() string {
	return name
}

func (s *DataSource) Sections() []string {
	return s.cfg.SubReddits
}

func (s *DataSource) GetPosts(sectionName string, selector feed.Selector) (feed.Posts, error) {
	subs, err := s.client.SubredditSubmissions(sectionName, geddit.HotSubmissions, geddit.ListingOptions{
		Limit: int(selector.Count),
		After: selector.AfterID,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to get subreddit '%s' submissions: %s", sectionName, err)
	}

	return submissionsToPosts(subs), nil
}

func (s *DataSource) Dispose() error {
	return nil
}

func NewDataSource(cfg Config) *DataSource {
	return &DataSource{
		client: geddit.NewSession(cfg.UserAgent),
		cfg:    cfg,
	}
}
