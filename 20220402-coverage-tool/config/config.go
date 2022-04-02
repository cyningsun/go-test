package config

import (
	"log"
	"regexp"
	"sync"
)

var (
	defaultConfig Config
	ignoreOnce    sync.Once
)

type Config struct {
	Ignore        []string `yaml:"ignore"`
	IgnorePattern []*regexp.Regexp
}

func IgnorePattern() []*regexp.Regexp {
	ignoreOnce.Do(func() {
		if len(defaultConfig.IgnorePattern) != len(defaultConfig.Ignore) {
			for _, pattern := range defaultConfig.Ignore {
				p, err := regexp.Compile(pattern)
				if err != nil {
					log.Fatalf("compile regex:%v", err)
				}
				defaultConfig.IgnorePattern = append(defaultConfig.IgnorePattern, p)
			}
		}
	})
	return defaultConfig.IgnorePattern
}
