package main

import (
	"fmt"
	"github.com/astaxie/beego/logs"
	"kafka-logMgr/kafka"
	"kafka-logMgr/models"
	"kafka-logMgr/tailf"
)

var (
	AppConf models.Config
)

func main() {
	filename := "../conf/my.conf"
	adapterType := "ini"
	err := AppConf.LoadConfig(adapterType, filename)
	if err != nil {
		fmt.Println("LoadConfig failed", err)
		return
	}
	err = AppConf.InitLogger()
	if err != nil {
		fmt.Println("InitLogger failed")
		return
	}
	err = tailf.InitTailf(nil, AppConf.ChanSize, AppConf.LogCollect)
	for _, v := range AppConf.LogCollect {
		err = tailf.InitTailf(v, AppConf.ChanSize, AppConf.LogCollect)
		if err != nil {
			logs.Info("No file need to tail.")
		}
	}

	err = kafka.InitKafka(AppConf.KafkaIp)
	if err != nil {
		logs.Error("InitKafka failed")
		return
	}

	logs.Debug("init all success\n")

	err = ServerRun()
	if err != nil {
		logs.Error("Server start failed")
		return
	}

	logs.Info("programa running success")
}
