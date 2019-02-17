package reddit

import (
	"encoding/json"
	"fmt"
	"github.com/x1unix/demo-go-plugins/server/feed"
)

type config struct {
	UserAgent string `json:"userAgent"`
}

// NewDataSource creates a new data source
//
// This method called by server on extension load
func NewDataSource(rawCfg json.RawMessage) (feed.Source, error) {
	cfg := new(config)
	if err := json.Unmarshal(rawCfg, cfg); err != nil {
		return nil, fmt.Errorf("invalid configuration format (%s)", err)
	}

	return newDataSource(*cfg), nil
}
