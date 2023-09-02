// config/config.go

package config

import (
	"log"
	"sync"

	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

// global instance of koanf and a mutex for thread-safety.
var (
	k    = koanf.New(".")
	once sync.Once
)

// LoadConfig loads the configuration from the specified path.
func LoadConfig(path string) {
	once.Do(func() {
		if err := k.Load(file.Provider(path), yaml.Parser()); err != nil {
			log.Fatalf("error loading config: %s", err)
		}
	})
}

// Get is a wrapper around k.Get for direct access to individual values.
func Get(key string) interface{} {
	return k.Get(key)
}
