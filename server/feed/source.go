package feed

import "encoding/json"

// SourceProvider is source provider function
type SourceProvider func(rawCfg json.RawMessage) (Source, error)

// Source represents external data source for posts
type Source interface {
	Name() string
	GetPosts(sectionName string, selector Selector) (Posts, error)
	Sections() []string
	Dispose() error
}
