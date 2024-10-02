package etcd

import (
	"context"
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/naming/resolver"
	"google.golang.org/grpc"
	"log"
	"time"
)

// EtcdDial 向grpc请求一个服务
func EtcdDial(c *clientv3.Client, service string) (*grpc.ClientConn, error) {
	etcdResolver, err := resolver.NewBuilder(c)
	if err != nil {
		return nil, err
	}
	return grpc.Dial(
		"etcd:///"+service,
		grpc.WithResolvers(etcdResolver),
		grpc.WithInsecure(),
		grpc.WithBlock(),
	)
}

// Discover 根据一致性哈希计算的节点名发现节点真实地址
func Discover(c *clientv3.Client, peerName string) string {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	keyToFind := "/cicicache/" + peerName
	resp, err := c.Get(ctx, keyToFind)
	if err != nil {
		log.Fatalf("Failed to get servers from etcd: %v", err)
	}
	// fmt.Println("成功从etcd获得kv！")
	if len(resp.Kvs) > 0 {
		fmt.Printf("Found key: %s, Value: %s\n From Etcd!", resp.Kvs[0].Key, resp.Kvs[0].Value)
		serverAddr := resp.Kvs[0].Value
		return string(serverAddr)
	} else {
		fmt.Println("Key not found:", keyToFind)
		return ""
	}
}
