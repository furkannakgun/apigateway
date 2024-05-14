package config

import (
	"errors"
	"os"

	"github.com/spf13/viper"
	_ "github.com/spf13/viper/remote"
	"github.vodafone.com/VFGroup-NetworkArchitecture-NAAP/structlog"
	"golang.org/x/exp/slices"
)

var Conf *viper.Viper
var ConfigDBHost string
var ServiceConfigKey string
var ConfigDBProvider string
var OpenAPIPath string = "/nebula-data/openAPIspec.yaml"
var logger = structlog.GetLogger("nebula")

// Mandatory config parameters
var mandatoryConfParams = [...]string{"targetURL"}

// Nebula config interface
var nc nebulaConfig

type nebulaConfig interface {
	loadConfig()
	loadOpenAPISpec()
	watchConfig(reloadConfig func())
}

func checkMandatoryConfig() error {
	// Check if any mandatory config parameters are missing
	for _, param := range mandatoryConfParams {
		if Conf.Get(param) == nil {
			return errors.New("Mandatory parameter missing: " + param)
		}
	}
	return nil
}

func LoadEnvVars() {
	if os.Getenv("ENV") == "dev" {
		logger.Warn("Local development environment configuration")
		logger.Warn("Mandatory config DB env vars not checked")
		// If running in a local environment, the openAPIspec.yaml
		// file shall be in the root directory (dev-only)
		logger.Warn("Loading local Open API spec from openAPIspec.yaml")
		OpenAPIPath = "openAPIspec.yaml"
		return
	}
	ConfigDBHost = os.Getenv("CONFIG_DB_URL")
	if ConfigDBHost == "" {
		panic("Config DB not set")
	}

	ServiceConfigKey = os.Getenv("CONFIG_DB_KEY")
	if ServiceConfigKey == "" {
		panic("Config DB microservice key not set")
	}

	ConfigDBProvider = os.Getenv("CONFIG_DB_PROVIDER")
	// Default to etcd3, options: etcd3, consul
	if ConfigDBProvider == "" {
		ConfigDBProvider = "etcd3"
	}

	dbProviders := []string{"etcd3", "consul"}
	if !slices.Contains(dbProviders, ConfigDBProvider) {
		panic("Config DB provider is not supported")
	}
}

func changeLoggingLevel() {
	loggingLevel := Conf.GetString("loggingLevel")
	if loggingLevel != "" {
		logger.ChangeLoggingLevel(loggingLevel)
		logger.Debugf("Logging level changed to %s", loggingLevel)
	}
}

func Load() {
	if os.Getenv("ENV") == "dev" {
		// If running in a local environment, we don't connect to a DB
		// The configuration file config.yaml shall exist in root dir (dev-only)
		nc = &localConfig{}
	} else {
		nc = &remoteConfig{}
	}
	// Order: load config, check config, load spec
	nc.loadConfig()
	err := checkMandatoryConfig()
	if err != nil {
		logger.Fatalv(err)
		panic(err)
	}
	changeLoggingLevel()
	nc.loadOpenAPISpec()
}

func WatchConfig(realoadConfig func()) {
	nc.watchConfig(realoadConfig)
}
