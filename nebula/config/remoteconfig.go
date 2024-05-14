package config

import (
	"os"

	"github.com/spf13/viper"
	"github.vodafone.com/VFGroup-NetworkArchitecture-NAAP/nebula/configdb"
)

var dbProvider configdb.ConfigDB

type remoteConfig struct{}

func (rc *remoteConfig) loadConfig() {
	Conf = viper.New()
	Conf.AddRemoteProvider(ConfigDBProvider, ConfigDBHost, ServiceConfigKey+"/config")
	Conf.SetConfigType("yaml")
	err := Conf.ReadRemoteConfig()
	if err != nil {
		logger.Fatalf("Error reading config: %v", err)
		panic(err)
	}
}

func (rc *remoteConfig) loadOpenAPISpec() {
	// Load Open API spec from the configured provider and store it to be used by the validator
	var spec []byte
	// When loading the spec we set up the db provider type on the module's var
	if ConfigDBProvider == "etcd3" {
		dbProvider = &configdb.ETCDConfigDB{}
	} else if ConfigDBProvider == "consul" {
		dbProvider = &configdb.ConsulConfigDB{}
	}
	// Read spec
	spec = dbProvider.ReadKey(ConfigDBHost, ServiceConfigKey+"/spec")

	err := os.WriteFile(OpenAPIPath, spec, 0644)
	if err != nil {
		logger.Fatalf("Error reading spec: %v", err)
		panic(err)
	}
}

func (rc *remoteConfig) watchConfig(reloadConfig func()) {
	dbProvider.WatchKey(ConfigDBHost, ServiceConfigKey, reloadConfig)
}
