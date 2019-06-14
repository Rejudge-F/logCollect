package models

type Config struct {
	LogPath    string
	LogLevel   int
	LogCollect []CollectCofig
	ChanSize   int
	KafkaIp    string
}

type CollectCofig struct {
	LogPath string
	Topic   string
}
