package models

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/config"
	"github.com/astaxie/beego/logs"
	"kafka-logMgr/etcd"
	"strconv"
	"strings"
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
	appConf.LoadEtcdConf(conf)
	err = etcd.InitEtcd(appConf.Etcd.Addr, appConf.Etcd.Key)
	if err != nil {
		logs.Error("Init Etcd Failed")
		return
	}
	err = appConf.LoadCollectConf(etcd.LocalIpArray)
	if err != nil {
		panic("Load CollectConf faield")
		return
	}

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

func (appConf *Config) LoadCollectConf(ip []string) (err error) {
	if strings.HasSuffix(appConf.Etcd.Key, "/") == false {
		appConf.Etcd.Key = fmt.Sprintf(appConf.Etcd.Key, "/")
	}
	for _, confIp := range ip {
		key := fmt.Sprintf("%s%s", appConf.Etcd.Key, confIp)
		etcd.EtcdCli.EtcdKeys = append(etcd.EtcdCli.EtcdKeys, key)
		collectStr, _ := etcd.GetKey(key)
		var logCollectConf []CollectConfig
		json.Unmarshal([]byte(collectStr), &logCollectConf)
		for _, collect := range logCollectConf {
			appConf.LogCollect[key] = append(appConf.LogCollect[key], collect)
			fmt.Println(collect)
		}
	}

	etcd.InitEtcdWatch()

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
