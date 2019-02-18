package main

import (
	"encoding/json"
	"fmt"

	"github.com/x1unix/demo-go-plugins/server/feed"
	"github.com/x1unix/demo-go-plugins/sources/reddit/extension"
)

// NewDataSource creates a new data source
//
// This method called by server on extension load
func NewDataSource(rawCfg json.RawMessage) (feed.Source, error) {
	cfg := new(extension.Config)
	if err := json.Unmarshal(rawCfg, cfg); err != nil {
		return nil, fmt.Errorf("invalid configuration format (%s)", err)
	}

	return extension.NewDataSource(*cfg), nil
}
