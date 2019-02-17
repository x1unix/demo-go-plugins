package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/x1unix/demo-go-plugins/server/feed/sources"
	"io/ioutil"
)

// Main is main configuration
var Main Config

type Config struct {
	Listen  string      `json:"listen"`
	Debug   bool        `json:"debug"`
	Sources sources.Set `json:"sources"`
}

// Load loads configuration from file
func Load(configPath string) (err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("failed to read config from file '%s': %s", configPath, err)
		}
	}()

	raw, err := ioutil.ReadFile(configPath)
	if err != nil {
		return err
	}

	if !json.Valid(raw) {
		return errors.New("found JSON syntax error")
	}

	if err := json.Unmarshal(raw, &Main); err != nil {
		return fmt.Errorf("failed to parse JSON file (%s)", err)
	}

	return err
}
