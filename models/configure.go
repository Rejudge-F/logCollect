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
	LogPath string
	Topic   string
}

type EtcdConfig struct {
	Addr []string
	Key  string
}
