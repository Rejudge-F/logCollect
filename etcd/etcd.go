package etcd

import (
	"context"
	"fmt"
	"github.com/astaxie/beego/logs"
	etcd_client "go.etcd.io/etcd/clientv3"
	"time"
)

var (
	ChKeyMessage chan *etcd_client.Event
)

type EtcdClient struct {
	client   *etcd_client.Client
	EtcdKeys []string
	EtcdAddr []string
}

var (
	EtcdCli *EtcdClient
)

func InitEtcd(addr []string, key string) error {
	ChKeyMessage = make(chan *etcd_client.Event, 100)
	cli, err := etcd_client.New(etcd_client.Config{
		Endpoints:   addr,
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		logs.Error(fmt.Sprintf("%s%s", "InitEtcd Failed", err))
		//panic("Init Etcd Failed")
		return err
	}

	EtcdCli = &EtcdClient{
		client:   cli,
		EtcdAddr: addr,
	}
	return nil
}

func GetKey(key string) (s string, err error) {
	s = ""
	err = nil
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	resp, err := EtcdCli.client.Get(ctx, key)
	if err != nil {
		logs.Error(fmt.Sprint("GetKey Failed: %v", err))
		return
	}
	kvs := resp.Kvs
	cancel()
	for _, v := range kvs {
		if string(v.Key) == key {
			return string(v.Value), nil
		}
	}
	return
}

func InitEtcdWatch() {
	logs.Info("Start Watch Etcd.")
	for _, key := range EtcdCli.EtcdKeys {
		go watchKey(key)
	}
}

func watchKey(key string) error {
	cli, err := etcd_client.New(etcd_client.Config{
		Endpoints:   EtcdCli.EtcdAddr,
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		logs.Error(err)
		return err
	}
	for {
		resp := cli.Watch(context.Background(), key)
		for msg := range resp {
			for _, ev := range msg.Events {
				fmt.Println(ev.Type, " ", string(ev.Kv.Key), " ", string(ev.Kv.Value))
				//if ev.Type == etcd_client.EventTypeDelete {
				//	DeleteCollect(string(ev.Kv.Key), string(ev.Kv.Value))
				//} else if ev.Type == etcd_client.EventTypePut {
				//	UpdateCollect(string(ev.Kv.Key), string(ev.Kv.Value))
				//}
				ChKeyMessage <- ev
			}
		}
	}
}
