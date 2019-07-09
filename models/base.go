package models

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/config"
	"github.com/astaxie/beego/logs"
	"strconv"
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
	appConf.ChanSize, err = conf.Int("LOG::ChanSize")

	if err != nil {
		appConf.ChanSize = 100
	}
	appConf.KafkaIp = conf.String("KAFKA::ServerIp")

	err = appConf.LoadCollectConf(conf)
	if err != nil {
		panic("Load CollectConf faield")
		return
	}
	appConf.LoadEtcdConf(conf)
	jsonStr, err := json.MarshalIndent(appConf, "\t", "")
	logs.Info(fmt.Sprintf("%v", string(jsonStr)))
	return nil
}

func (appConf *Config) LoadEtcdConf(configure config.Configer) (err error) {
	etcdKey := configure.String("ETCD::EtcdKey")
	appConf.Etcd.Key = etcdKey
	for i := 1; ; i++ {
		addrkey := "EtcdAddr" + strconv.Itoa(i)
		addrValue := configure.String("ETCD::" + addrkey)
		if len(addrValue) == 0 {
			break
		}
		appConf.Etcd.Addr = append(appConf.Etcd.Addr, addrValue)
	}
	return nil
}

func (appConf *Config) LoadCollectConf(configure config.Configer) (err error) {
	for i := 1; ; i++ {
		cc := CollectConfig{}
		path := "COLLECTLOG::CollectLogPath" + strconv.Itoa(i)
		cc.LogPath = configure.String(path)
		if len(cc.LogPath) == 0 {
			break
		}
		topic := "COLLECTLOG::CollectLogTopic" + strconv.Itoa(i)
		cc.Topic = configure.String(topic)
		if len(cc.Topic) == 0 {
			cc.Topic = "default"
		}
		appConf.LogCollect = append(appConf.LogCollect, cc)
	}

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
