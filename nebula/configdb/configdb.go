package configdb

import (
	"github.vodafone.com/VFGroup-NetworkArchitecture-NAAP/structlog"
)

type ConfigDB interface {
	ReadKey(dbHost string, key string) []byte
	WatchKey(dbHost string, key string, reloadConfig func())
}

var logger = structlog.GetLogger("nebula")
