package tailf

import (
	"encoding/json"
	"errors"
	"github.com/astaxie/beego/logs"
	"github.com/hpcloud/tail"
	"kafka-logMgr/models"
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
	Tails   []*TailObj
	MsgChan chan *TailMessage
}

func InitTailf(config []models.CollectConfig, ChanSize int) error {
	TailObjManager = &TailObjMgr{
		Tails:   make([]*TailObj, 2),
		MsgChan: make(chan *TailMessage, ChanSize),
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
	var configs []models.CollectConfig
	json.Unmarshal([]byte(value), &configs)
	for _, addConfig := range configs {
		exist := false
		for _, config := range TailObjManager.Tails {
			if addConfig.LogPath == config.Conf.LogPath {
				exist = true
				break
			}
		}
		if !exist {
			TailFile(addConfig)
		}
	}
}

func DeleteCollect(key, value string) {
	var configs []models.CollectConfig
	json.Unmarshal([]byte(value), &configs)
	for _, delConfig := range configs {
		for _, config := range TailObjManager.Tails {
			if delConfig.LogPath == config.Conf.LogPath {
				config.Shut = true
			}
		}
	}
}
