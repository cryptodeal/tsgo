package config

import (
	"log"
	"os"

	"github.com/cryptodeal/tsgo/tsgo"
	"gopkg.in/yaml.v2"
)

const defaultFallbackType = "any"

func ReadFromFilepath(cfgFilepath string) tsgo.Config {
	b, err := os.ReadFile(cfgFilepath)
	if err != nil {
		log.Fatalf("Could not read config file from %s: %v", cfgFilepath, err)
	}
	conf := tsgo.Config{}
	err = yaml.Unmarshal(b, &conf)
	if err != nil {
		log.Fatalf("Could not parse config file froms: %v", err)
	}

	// apply defaults
	for _, packageConf := range conf.Packages {
		if packageConf.FallbackType == "" {
			packageConf.FallbackType = defaultFallbackType
		}
	}

	return conf
}
