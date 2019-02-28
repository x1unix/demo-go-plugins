package feed

import "encoding/json"

// SourceProvider is source provider function
type SourceProvider func(rawCfg json.RawMessage) (SourceReader, error)

// SourceReader represents external data source for posts
type SourceReader interface {
	Name() string
	GetPosts(sectionName string, selector Selector) (Posts, error)
	Sections() []string
	Dispose() error
}
