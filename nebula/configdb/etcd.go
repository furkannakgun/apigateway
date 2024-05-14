package configdb

import (
	"context"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

type ETCDConfigDB struct{}

func (db *ETCDConfigDB) ReadKey(dbHost string, key string) []byte {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{dbHost},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		logger.Fatalf("Error reading config from DB (ETCD) : %s", err)
		panic(err)
	}

	defer cli.Close()

	kv := clientv3.NewKV(cli)
	resp, err := kv.Get(context.Background(), key)
	if err != nil {
		logger.Fatalf("Error reading key: %v", err)
		panic(err)
	}
	return resp.Kvs[0].Value
}

func (db *ETCDConfigDB) WatchKey(dbHost string, key string, reloadConfig func()) {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{dbHost},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		logger.Fatalf("Error creating config DB client: %s", err)
		panic(err)
	}

	watchChannel := cli.Watch(context.Background(), key, clientv3.WithPrefix())

	go func() {
		defer cli.Close()
		for {
			for wresp := range watchChannel {
				if wresp.Canceled {
					logger.Warn("ETCD channel closed")
					break
				}
				for range wresp.Events {
					logger.Info("Configuration reload triggered by event reception")
					reloadConfig()
				}
			}
		}
	}()
}
