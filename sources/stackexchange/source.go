package main

import (
	"net/http"
	"strconv"

	"github.com/x1unix/demo-go-plugins/server/feed"
)

type DataSource struct {
	cfg    config
	client apiClient
}

func (s *DataSource) Name() string {
	return s.cfg.Site
}

func (s *DataSource) Sections() []string {
	return s.cfg.Tags
}

func (s *DataSource) GetPosts(sectionName string, selector feed.Selector) (feed.Posts, error) {
	questions, err := s.client.getQuestions(selector.Count, sectionName)
	if err != nil {
		return nil, err
	}

	return s.questionsToPosts(questions), nil
}

func (s *DataSource) Dispose() error {
	return nil
}

func (s *DataSource) questionsToPosts(questions []question) feed.Posts {
	r := make(feed.Posts, 0, len(questions))
	for _, q := range questions {
		post := feed.Post{
			ID:         strconv.FormatInt(q.QuestionID, 10),
			URL:        q.Link,
			Title:      q.Title,
			SourceType: s.cfg.Site,
			CreatedAt:  q.CreationDate,
		}

		r = append(r, post)
	}

	return r
}

func newDataSource(cfg config) *DataSource {
	return &DataSource{
		cfg: cfg,
		client: apiClient{
			httpClient: &http.Client{},
			site:       cfg.Site,
		},
	}
}
