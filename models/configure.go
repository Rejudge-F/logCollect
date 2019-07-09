package models

type Config struct {
	LogPath    string
	LogLevel   int
	LogCollect []CollectConfig
	ChanSize   int
	KafkaIp    string
	Etcd       EtcdConfig
}

type CollectConfig struct {
	LogPath string `json:"logpath"`
	Topic   string `json:"topic"`
}

type EtcdConfig struct {
	Addr []string
	Key  string
}
