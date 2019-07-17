package tailf

import (
	"errors"
	"github.com/astaxie/beego/logs"
	"github.com/hpcloud/tail"
	"kafka-logMgr/models"
)

var (
	tailObjMgr *TailObjMgr
)

type TailObj struct {
	Tail *tail.Tail
	Conf models.CollectConfig
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
	tailObjMgr = &TailObjMgr{
		Tails:   make([]*TailObj, 2),
		MsgChan: make(chan *TailMessage, ChanSize),
	}

	if len(config) == 0 {
		return errors.New("No Collect Config need to handle")
	}
	conf := tail.Config{}
	conf.MustExist = false
	conf.Poll = true
	conf.ReOpen = true
	conf.Follow = true

	for _, v := range config {
		tails, err := tail.TailFile(v.LogPath, conf)
		if err != nil {
			logs.Error("Read Collect failed: InitTailf")
		}
		tailObj := &TailObj{
			Tail: tails,
			Conf: v,
		}
		tailObjMgr.Tails = append(tailObjMgr.Tails, tailObj)

		go ReadFromTailObj(tailObj)

	}

	return nil
}

func ReadFromTailObj(tailObj *TailObj) {
	for true {
		line, ok := <-tailObj.Tail.Lines
		if !ok {
			logs.Error("Read failed from %s\n", tailObj.Conf.LogPath)
		}
		msg := &TailMessage{
			TextMsg: line.Text,
			Topic:   tailObj.Conf.Topic,
		}
		tailObjMgr.MsgChan <- msg
	}
}

func GetLine() *TailMessage {
	tailMsg := &TailMessage{}
	tailMsg = <-tailObjMgr.MsgChan
	return tailMsg
}
