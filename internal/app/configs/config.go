package configs

var AppConfig *Config

type Config struct {
	ServAddr     string
	BaseShortURL string
}

func InitConfig(servAddr, shortURL string) {
	AppConfig = &Config{
		ServAddr:     servAddr,
		BaseShortURL: shortURL,
	}
}
