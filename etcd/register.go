package etcd

import (
	"context"
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
	"log"
)

// register模块提供服务Service注册至etcd的能力

// Add 在租赁模式添加一对kv至etcd
func Add(c *clientv3.Client, leaseId clientv3.LeaseID, service string, addr string) error {
	// 创建服务路径
	key := fmt.Sprintf("/cicicache/%s", service)

	// 创建键值对
	val := addr // 假设addr是"IP:Port"格式
	// 获取KV API
	kv := clientv3.NewKV(c)

	// 设置键值对，附加租约ID
	_, err := kv.Put(c.Ctx(), key, val, clientv3.WithLease(leaseId))
	if err != nil {
		return fmt.Errorf("failed to put key-value in etcd: %v", err)
	}
	fmt.Printf("成功设置key：%s, value: %s\n", key, val)
	return nil
}

// Register 注册一个服务至etcd
func Register(service string, addr string, stop chan error) (*clientv3.Client, error) {
	// 创建一个etcd client
	cli, err := ConnectToEtcd()
	// 创建一个租约 配置5秒过期
	resp, err := cli.Grant(context.Background(), 5)
	if err != nil {
		cli.Close() // 确保出错时关闭客户端
		return nil, fmt.Errorf("create lease failed: %v", err)
	}
	leaseId := resp.ID

	// 注册服务
	err = Add(cli, leaseId, service, addr)
	if err != nil {
		cli.Close()
		return nil, fmt.Errorf("add etcd record failed: %v", err)
	}

	// 设置服务心跳检测
	ch, err := cli.KeepAlive(context.Background(), leaseId)
	if err != nil {
		cli.Close()
		return nil, fmt.Errorf("set keepalive failed: %v", err)
	}

	log.Printf("[%s] register service ok\n", addr)
	go func() {
		for {
			select {
			case err := <-stop:
				if err != nil {
					log.Println(err)
				}
				cli.Close()
				return
			case <-cli.Ctx().Done():
				log.Println("service closed")
				cli.Close()
				return
			case _, ok := <-ch:
				if !ok {
					log.Println("keep alive channel closed")
					cli.Revoke(context.Background(), leaseId)
					cli.Close()
					return
				}
			}
		}
	}()
	return cli, nil
}
