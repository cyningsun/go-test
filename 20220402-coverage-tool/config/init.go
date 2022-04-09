package config

import (
	"log"

	"github.com/cyningsun/go-test/20220402-coverage-tool/walk"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
)

func init() {
	v := viper.New()

	v.SetConfigName(".xcover") // name of config file (without extension)
	v.SetConfigType("yml")     // REQUIRED if the config file does not have the extension in the name
	v.AddConfigPath("$HOME/")  // call multiple times to add many search paths
	v.AddConfigPath(".")       // optionally look for config in the working directory

	ancestor, _ := walk.Ancestor(".")
	for _, each := range ancestor {
		viper.AddConfigPath(each)
	}

	if err := v.ReadInConfig(); err != nil {
		log.Fatalf("Fatal error config file: %w \n", err)
	}

	if err := v.Unmarshal(&defaultConfig, viper.DecoderConfigOption(func(decoderConfig *mapstructure.DecoderConfig) {
		decoderConfig.TagName = "yaml"
	})); err != nil { // Handle errors reading the config file
		log.Fatalf("Fatal error config file: %w \n", err)
	}
}
