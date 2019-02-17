package feed

import (
	"fmt"
	"github.com/sirupsen/logrus"
)

var loadedSources = make(map[string]Source)

// AddSource registers a new external feed source
func AddSource(src Source) {
	loadedSources[src.Name()] = src
}

// GetSource gets source by name
func GetSource(name string) (Source, error) {
	src, ok := loadedSources[name]
	if !ok {
		return nil, fmt.Errorf("unknown feed source '%s'", name)
	}

	return src, nil
}

// SourceNames returns names of known sources
func SourceNames() []string {
	names := make([]string, 0, len(loadedSources))
	for k := range loadedSources {
		names = append(names, k)
	}

	return names
}

// FlushSources unregisters all sources
func FlushSources() {
	for k := range loadedSources {
		if err := loadedSources[k].Dispose(); err != nil {
			logrus.Warnf("source '%s' returned error on dispose (%s)", k, err)
		}

		delete(loadedSources, k)
		logrus.Debugf("source '%s' unregistered", k)
	}
}
