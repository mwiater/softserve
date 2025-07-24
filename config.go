// config.go
package softserve

import (
	"sync"
)

type Config struct {
	WebRoot   string
	SSL       bool
	HTTPPort  int
	HTTPSPort int
	LogLevel  string
	API       bool
	APIPrefix string
}

var (
	appConfig Config
	once      sync.Once
)

// InitConfig sets the configuration once. Subsequent calls are ignored.
func InitConfig(c Config) {
	once.Do(func() {
		appConfig = c
	})
}

func GetConfig() Config {
	return appConfig
}
