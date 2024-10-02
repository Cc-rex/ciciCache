package etcd

import (
	clientv3 "go.etcd.io/etcd/client/v3"
	"log"
	"time"
)

const defaultEtcdAddr = "localhost:2379"
const defaultDialTime = 5 // in seconds

func ConnectToEtcd(options ...interface{}) (*clientv3.Client, error) {
	addr := defaultEtcdAddr
	dialTime := defaultDialTime

	if len(options) > 0 {
		if a, ok := options[0].(string); ok && a != "" {
			addr = a
		}
		if len(options) > 1 {
			if dt, ok := options[1].(int); ok && dt != 0 {
				dialTime = dt
			}
		}
	}

	clientConfig := clientv3.Config{
		Endpoints:   []string{addr},
		DialTimeout: time.Duration(dialTime) * time.Second,
	}
	client, err := clientv3.New(clientConfig)
	if err != nil {
		log.Fatalf("Failed to connect to etcd: %v", err)
	}
	return client, err
}
