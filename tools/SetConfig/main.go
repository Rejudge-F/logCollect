package main

import (
	"context"
	"encoding/json"
	"fmt"
	etcd_client "go.etcd.io/etcd/clientv3"
	"kafka-logMgr/models"
	"time"
)

var (
	etcdKey = "EtcdKey/192.168.137.233"
)

func SetLogConfToEtcd() error {
	cli, err := etcd_client.New(etcd_client.Config{
		Endpoints:   []string{"localhost:2379"},
		DialTimeout: 2 * time.Second,
	})
	if err != nil {
		fmt.Println("Etcd Create New Client Failed: ", err)
		return err
	}
	fmt.Println("Connect Etcd Success!")
	defer cli.Close()
	var collectConf []models.CollectConfig
	collectConf = append(collectConf, models.CollectConfig{
		LogPath: "../logs/collect.log",
		Topic:   "test",
	})
	collectConf = append(collectConf, models.CollectConfig{
		LogPath: "../logs/collect2.log",
		Topic:   "collect2",
	})
	confStr, err := json.Marshal(collectConf)
	if err != nil {
		fmt.Println("json Faild: ", err)
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	_, err = cli.Put(ctx, etcdKey, string(confStr))
	if err != nil {
		fmt.Println(err)
		return err
	}
	cancel()
	return nil
}

func main() {
	SetLogConfToEtcd()
}
