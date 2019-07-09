package main

import (
	"github.com/astaxie/beego/logs"
	"kafka-logMgr/kafka"
	"kafka-logMgr/tailf"
)

func ServerRun() error {
	for true {
		msg := tailf.GetLine()

		err := kafka.SendToKafka(msg.TextMsg, msg.Topic)
		if err != nil {
			logs.Error("Send to kafka failed")
		}
	}
	return nil
}
