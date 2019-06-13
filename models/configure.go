package models

type Config struct {
	LogPath    string
	LogLevel   int
	LogCollect []CollectCofig
}

type CollectCofig struct {
	LogPath  string
	LogLevel int
}
