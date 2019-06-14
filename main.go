package main

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/logs"
	"kafka-logMgr/kafka"
	"kafka-logMgr/models"
	"kafka-logMgr/tailf"
)

var (
	appConf models.Config
)

func main() {
	filename := "./conf/my.conf"
	adapterType := "ini"
	err := appConf.LoadConfig(adapterType, filename)
	if err != nil {
		fmt.Println("LoadConfig failed", err)
		return
	}
	err = appConf.InitLogger()
	if err != nil {
		fmt.Println("InitLogger failed")
		return
	}
	confStr, err := json.MarshalIndent(appConf, " ", "")
	if err != nil {
		fmt.Println("json failed")
	}

	err = tailf.InitTailf(appConf.LogCollect, appConf.ChanSize)
	if err != nil {
		logs.Error("InitTailf failed", err)
		return
	}

	err = kafka.InitKafka(appConf.KafkaIp)
	if err != nil {
		logs.Error("InitKafka failed")
		return
	}
	logs.Debug("init all success\n")
	logs.Debug("%v\n", string(confStr))

	err = ServerRun()
	if err != nil {
		logs.Error("Server start failed")
		return
	}

	logs.Info("programa running success")
}
