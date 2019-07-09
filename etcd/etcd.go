package etcd

import (
	"context"
	"fmt"
	"github.com/astaxie/beego/logs"
	etcd_client "go.etcd.io/etcd/clientv3"
	"strings"
	"time"
)

type EtcdClient struct {
	client *etcd_client.Client
}

var (
	etcdClient *EtcdClient
)

func InitEtcd(addr []string, key string) error {
	cli, err := etcd_client.New(etcd_client.Config{
		Endpoints:   addr,
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		logs.Error(fmt.Sprintf("%s%s", "InitEtcd Failed", err))
		panic("Init Etcd Failed")
		return err
	}

	etcdClient = &EtcdClient{
		client: cli,
	}

	if strings.HasSuffix(key, "/") == false {
		key = fmt.Sprintf("%s%s", key, "/")
	}

	for _, ip := range localIpArray {
		etcdKey := fmt.Sprintf("%s%s", key, ip)
		logs.Info(etcdKey)
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		resp, err := etcdClient.client.Get(ctx, etcdKey)
		if err != nil {
			logs.Error(fmt.Sprintf("%s: %v", "Etcd Get Key Failed", err))
			continue
		}
		cancel()
		logs.Debug(fmt.Sprintf("%v", resp))
		for k, v := range resp.Kvs {
			logs.Info(fmt.Sprintf("etcd: [key: %v] [v: %v]", k, v.Value))
		}
	}
	return nil
}
