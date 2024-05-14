package config

import (
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

type localConfig struct{}

func (lc *localConfig) loadConfig() {
	// Load config to test in local PC
	// File shall be config.yaml in root dir (of the repo)
	Conf = viper.New()
	logger.Warn("Loading local configuration from config.yaml")
	Conf.SetConfigFile("config.yaml")
	err := Conf.ReadInConfig()
	if err != nil {
		panic(err)
	}
}

func (lc *localConfig) loadOpenAPISpec() {
	// Nothing to do but fulfill the interface
	return
}

func (lc *localConfig) watchConfig(realoadConfig func()) {
	Conf.OnConfigChange(func(e fsnotify.Event) {
		logger.Infof("Config file changed: %s", e.Name)
		realoadConfig()
	})
	Conf.WatchConfig()
}
