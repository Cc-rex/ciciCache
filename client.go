package ciciCache

import (
	"ciciCache/etcd"
	"ciciCache/protobufs"
	"context"
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
	"time"
)

// client 模块实现cicicache访问其他远程节点 从而获取缓存的能力

type client struct {
	name string // 服务名称 pcache/ip:addr
}

// Fetch 从远程节点获取缓存,实现Fetcher接口
func (c *client) Fetch(group string, key string) ([]byte, error) {
	// 创建一个etcd client
	cli, err := clientv3.New(defaultEtcdConfig)
	if err != nil {
		return nil, err
	}
	defer cli.Close()

	// 发现服务 取得与服务的连接
	conn, err := etcd.EtcdDial(cli, c.name)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	grpcClient := protobufs.NewCiciCacheClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	resp, err := grpcClient.Get(ctx, &protobufs.GetRequest{
		Group: group,
		Key:   key,
	})
	if err != nil {
		return nil, fmt.Errorf("could not get %s/%s from peer %s", group, key, c.name)
	}

	return resp.GetValue(), nil
}

func NewClient(service string) *client {
	return &client{name: service}
}
