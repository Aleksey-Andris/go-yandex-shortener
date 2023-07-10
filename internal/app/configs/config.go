package configs

var AppConfig *Config

type Config struct {
	ServAddr     string
	BaseShortUrl string
}

func InitConfig(servAddr, shortURL string) {
	AppConfig = &Config{
		ServAddr:     servAddr,
		BaseShortUrl: shortURL,
	}
}
