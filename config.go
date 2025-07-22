// config.go
package softserve

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	WebRoot       string `mapstructure:"web_root"`
	SSL           bool   `mapstructure:"ssl"`
	GenerateCerts bool   `mapstructure:"generate_certs"`
	HTTPPort      int    `mapstructure:"http_port"`
	HTTPSPort     int    `mapstructure:"https_port"`
	LogLevel      string `mapstructure:"log_level"`
	API           bool   `mapstructure:"api"`
	APIPrefix     string `mapstructure:"api_prefix"`
}

var AppConfig Config

func LoadConfig() error {
	viper.SetConfigName("softserve")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".") // look in current dir

	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("error reading config file: %w", err)
	}

	if err := viper.Unmarshal(&AppConfig); err != nil {
		return fmt.Errorf("unable to decode into struct: %w", err)
	}

	return nil
}
