package kafka

import (
	"github.com/Shopify/sarama"
	"github.com/astaxie/beego/logs"
)

var (
	producer sarama.SyncProducer
)

func InitKafka(serverIp string) (err error) {
	conf := sarama.NewConfig()
	conf.Producer.RequiredAcks = sarama.WaitForAll
	conf.Producer.Partitioner = sarama.NewRandomPartitioner
	conf.Producer.Return.Successes = true

	producer, err = sarama.NewSyncProducer([]string{serverIp}, conf)

	if err != nil {
		logs.Error("Init AsyncProducer failed")
		return err
	}

	return nil
}

func SendToKafka(msg, topic string) error {
	produceMsg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(msg),
	}
	partition, offset, err := producer.SendMessage(produceMsg)
	if err != nil {
		logs.Error("Send to kafka failed")
		return err
	}

	logs.Info("partition: %v, offset: %v\n", partition, offset)
	return nil
}
