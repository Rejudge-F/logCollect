package main

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/logs"
	"kafka-logMgr/models"
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
	logs.Debug("init all success\n")
	logs.Debug("%v\n", string(confStr))

}
