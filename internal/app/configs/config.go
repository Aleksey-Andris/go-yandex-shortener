// The configs package is designed to configure the application.
package configs

// AppConfig - golabal varible containing launch parameters.
var AppConfig *Config

// AppConfig - structure containing launch parameters.
//
// ServAddr — server launch address.
//
// BaseShortURL — server address for shortened URLs.
type Config struct {
	ServAddr string

	BaseShortURL string
}

// InitConfig - constructor for Config.
func InitConfig(servAddr, shortURL string) {
	AppConfig = &Config{
		ServAddr:     servAddr,
		BaseShortURL: shortURL,
	}
}
