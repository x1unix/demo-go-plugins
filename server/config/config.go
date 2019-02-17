package config

import "github.com/x1unix/demo-go-plugins/server/feed/sources"

type Config struct {
	Listen  string      `json:"listen"`
	Debug   bool        `json:"debug"`
	Sources sources.Set `json:"sources"`
}
