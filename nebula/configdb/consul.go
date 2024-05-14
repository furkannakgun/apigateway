package configdb

import (
	"context"
	"time"

	capi "github.com/hashicorp/consul/api"
	capiwatcher "github.com/pteich/consul-kv-watcher"
)

type ConsulConfigDB struct{}

func (db *ConsulConfigDB) ReadKey(dbHost string, key string) []byte {
	capiConfig := capi.DefaultNonPooledConfig()
	capiConfig.Address = dbHost
	cli, err := capi.NewClient(capiConfig)
	if err != nil {
		logger.Fatalf("Error opening config DB client: %v", err)
		panic(err)
	}
	// KV handle
	kv := cli.KV()
	pair, _, err := kv.Get(key, nil)
	if err != nil {
		logger.Fatalf("Error reading key: %v", err)
		panic(err)
	}
	return pair.Value
}

func (db *ConsulConfigDB) WatchKey(dbHost string, key string, reloadConfig func()) {
	capiConfig := capi.DefaultConfig()
	capiConfig.Address = dbHost
	client, err := capi.NewClient(capiConfig)
	if err != nil {
		logger.Fatalf("Error opening client to watch config: %v", err)
		panic(err)
	}
	watcher := capiwatcher.New(client, 10*time.Second, 2*time.Second)
	watchChannel, err := watcher.WatchTree(context.Background(), key)
	if err != nil {
		logger.Fatalf("Error opening config watcher channel: %v", err)
		panic(err)
	}

	go func() {
		for {
			for range watchChannel {
				logger.Infof("Configuration reload triggered by event reception")
				reloadConfig()
			}
		}
	}()
}
