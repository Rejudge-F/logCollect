package models

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/astaxie/beego/config"
	"github.com/astaxie/beego/logs"
)

func ConvertLogLevel(logLevel string) int {
	switch logLevel {
	case "Debug":
		return logs.LevelDebug
	case "Alert":
		return logs.LevelAlert
	case "Warning":
		return logs.LevelWarn
	case "Error":
		return logs.LevelError
	case "Notice":
		return logs.LevelNotice
	default:
		return 0
	}
}

func (appConf *Config) LoadConfig(adapterType, filename string) (err error) {
	conf, err := config.NewConfig(adapterType, filename)
	if err != nil {
		fmt.Println("LoadConfig faield", err)
		return err
	}
	appConf.LogPath = conf.String("LOG::LogPath")
	appConf.LogLevel = ConvertLogLevel(conf.String("LOG::LogLevel"))

	return nil
}

func (appConf *Config) LoadCollectConf(configer config.Configer) (err error) {
	cc := CollectCofig{}
	cc.LogPath = configer.String("COLLECTLOG::LogPath")
	if len(cc.LogPath) == 0 {
		return errors.New("invaild LogPath")
	}
	cc.LogLevel = ConvertLogLevel(configer.String("COLLECTLOG::LogLevel"))
	if cc.LogLevel == 0 {
		cc.LogLevel = ConvertLogLevel("Debug")
	}
	appConf.LogCollect = append(appConf.LogCollect, cc)
	return nil

}

func (appConf *Config) InitLogger() error {
	conf := make(map[string]interface{})
	conf["filename"] = appConf.LogPath
	conf["level"] = appConf.LogLevel

	confStr, err := json.Marshal(conf)
	if err != nil {
		return err
	}
	err = logs.SetLogger(logs.AdapterFile, string(confStr))
	if err != nil {
		return err
	}
	return nil
}
