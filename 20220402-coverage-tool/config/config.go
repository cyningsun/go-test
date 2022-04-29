package config

import (
	"sync"
)

var (
	defaultConfig Config
	ignoreOnce    sync.Once
)

type Config struct {
	Ignore []string `yaml:"ignore"`
}

func IgnorePattern() []string {
	return defaultConfig.Ignore
}
