package sources

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/x1unix/demo-go-plugins/server/feed"
	"plugin"
)

const (
	pluginEntrypoint = "NewDataSource"
	pluginLoadErr    = "failed to load data source from library '%s', %s"
)

// Set is map of source library name to load and with params
type Set map[string]json.RawMessage

// Load loads data sources from source set
func Load(sources Set) error {
	for pluginPath, cfg := range sources {
		provider, err := loadProviderFromPlugin(pluginPath)
		if err != nil {
			return fmt.Errorf(pluginLoadErr, pluginPath, err)
		}

		if err = injectDataSource(provider, cfg); err != nil {
			return fmt.Errorf(pluginLoadErr, pluginPath, err)
		}

		logrus.Infof("'%s' loaded", pluginPath)
	}
	return nil
}

func injectDataSource(sourceProvider feed.SourceProvider, cfg json.RawMessage) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic: %s", r)
		}
	}()

	dataSource, err := sourceProvider(cfg)
	if err != nil {
		return fmt.Errorf("cannot construct data source: %s", err)
	}

	feed.AddSource(dataSource)
	return nil
}

func loadProviderFromPlugin(fileName string) (feed.SourceProvider, error) {
	p, err := plugin.Open(fileName)
	if err != nil {
		return nil, err
	}

	fnPtr, err := p.Lookup(pluginEntrypoint)
	if err != nil {
		return nil, fmt.Errorf("cannot find entrypoint function '%s' in library (%s)", pluginEntrypoint, err)
	}

	sourceProvider, ok := fnPtr.(feed.SourceProvider)
	if !ok {
		return nil, fmt.Errorf("invalid '%s' function signature, expected %T (got %T)", pluginEntrypoint, sourceProvider, fnPtr)
	}

	return sourceProvider, nil
}
