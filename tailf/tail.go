package tailf

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/hpcloud/tail"
	etcd_client "go.etcd.io/etcd/clientv3"
	"kafka-logMgr/etcd"
	"kafka-logMgr/models"
	"sync"
)

var (
	TailObjManager *TailObjMgr
)

type TailObj struct {
	Tail *tail.Tail
	Conf models.CollectConfig
	Shut bool
}

type TailMessage struct {
	TextMsg string
	Topic   string
}

type TailObjMgr struct {
	mutex   sync.Mutex
	Tails   []*TailObj
	MsgChan chan *TailMessage
}

func InitTailf(config []models.CollectConfig, ChanSize int, collects map[string][]models.CollectConfig) error {
	if TailObjManager == nil {
		TailObjManager = &TailObjMgr{
			Tails:   make([]*TailObj, 2),
			MsgChan: make(chan *TailMessage, ChanSize),
		}
		go WatchEtcd(collects)
	}

	if len(config) == 0 {
		return errors.New("No Collect Config need to handle")
	}
	for _, v := range config {
		TailFile(v)
	}

	return nil
}

func ReadFromTailObj(tailObj *TailObj) {
	for true {
		if tailObj.Shut {
			break
		}
		fmt.Println(tailObj)
		line, ok := <-tailObj.Tail.Lines
		if !ok {
			logs.Error("Read failed from %s\n", tailObj.Conf.LogPath)
		}
		msg := &TailMessage{
			TextMsg: line.Text,
			Topic:   tailObj.Conf.Topic,
		}
		TailObjManager.MsgChan <- msg
	}
}

func GetLine() *TailMessage {
	tailMsg := &TailMessage{}
	tailMsg = <-TailObjManager.MsgChan
	return tailMsg
}

func TailFile(v models.CollectConfig) {
	conf := tail.Config{}
	conf.MustExist = false
	conf.Poll = true
	conf.ReOpen = true
	conf.Follow = true
	tails, err := tail.TailFile(v.LogPath, conf)
	if err != nil {
		logs.Error("Read Collect failed: InitTailf")
	}
	tailObj := &TailObj{
		Tail: tails,
		Conf: v,
		Shut: false,
	}
	TailObjManager.Tails = append(TailObjManager.Tails, tailObj)
	go ReadFromTailObj(tailObj)
}

func UpdateCollect(key, value string) {
	configs := make([]models.CollectConfig, 100)
	json.Unmarshal([]byte(value), &configs)
	for _, addConfig := range configs {
		exist := false
		if TailObjManager.Tails != nil {
			fmt.Println(TailObjManager.Tails)
			for _, config := range TailObjManager.Tails {
				if config == nil {
					continue
				}
				if addConfig.LogPath == config.Conf.LogPath {
					exist = true
					break
				}
			}
		}
		if !exist {
			TailFile(addConfig)
		}
	}
}

func DeleteCollect(key string, collects map[string][]models.CollectConfig) {
	configs := collects[key]
	for _, delConfig := range configs {
		for _, config := range TailObjManager.Tails {
			if config == nil {
				continue
			}
			if delConfig.LogPath == config.Conf.LogPath {
				config.Shut = true
				fmt.Println("delconfig", config)
				break
			}
		}
	}
}

func WatchEtcd(collects map[string][]models.CollectConfig) {
	for true {
		if ev, ok := <-etcd.ChKeyMessage; ok {
			TailObjManager.mutex.Lock()
			if ev.Type == etcd_client.EventTypeDelete {
				DeleteCollect(string(ev.Kv.Key), collects)
			} else if ev.Type == etcd_client.EventTypePut {
				UpdateCollect(string(ev.Kv.Key), string(ev.Kv.Value))
			}
			TailObjManager.mutex.Unlock()
		}
	}
}
