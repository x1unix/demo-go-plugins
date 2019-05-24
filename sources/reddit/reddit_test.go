package main

import (
	"github.com/golang/mock/gomock"
	"github.com/jzelinskie/geddit"
	"github.com/stretchr/testify/assert"
	"github.com/x1unix/demo-go-plugins/server/feed"
	"github.com/x1unix/demo-go-plugins/sources/reddit/mocks"
	"testing"
)

func TestDataSource_Name(t *testing.T) {
	ds := DataSource{}
	assert.Equal(t, name, ds.Name())
}

func TestDataSource_Sections(t *testing.T) {
	exp := []string{"foo", "bar"}
	ds := DataSource{
		cfg: config{SubReddits: exp},
	}

	assert.Equal(t, exp, ds.Sections())
}

func TestDataSource_GetPosts(t *testing.T) {
	rdt := mocks.NewRedditMock(gomock.NewController(t))
	rdt.EXPECT().SubredditSubmissions("foo", gomock.Any(), geddit.ListingOptions{
		Limit: 10,
		After: "",
	}).Return([]*geddit.Submission{
		{
			ID:    "33",
			Title: "foo",
		},
	}, nil)

	expected := feed.Posts{
		{
			ID:         "33",
			Title:      "foo",
			SourceType: name,
		},
	}

	ds := DataSource{
		client: rdt,
	}

	got, err := ds.GetPosts("foo", feed.Selector{Count: 10})
	assert.NoError(t, err)
	assert.Equal(t, expected, got)
}
