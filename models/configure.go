package models

type Config struct {
	LogPath    string
	LogLevel   int
	LogCollect []CollectCofig
	ChanSize   int
}

type CollectCofig struct {
	LogPath string
	Topic   string
}
