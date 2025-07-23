// config.go
package softserve

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

type Config struct {
	WebRoot       string `mapstructure:"web_root"`
	SSL           bool   `mapstructure:"ssl"`
	GenerateCerts bool   `mapstructure:"generate_certs"`
	CertsPath     string `mapstructure:"certs_path"`
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

	// === IMPORTANT: Add this line to convert the path after unmarshaling ===
	AppConfig.CertsPath = ConvertPath(AppConfig.CertsPath)

	if AppConfig.API {
		// Assuming LoadAPIResponses is defined elsewhere in the 'softserve' package
		if err := LoadAPIResponses(); err != nil {
			return fmt.Errorf("failed to load api.yaml: %w", err)
		}
	}

	if AppConfig.SSL {
		if AppConfig.CertsPath == "" {
			return fmt.Errorf("certs_path is required when ssl is enabled")
		}

		// EnsureAbsoluteAndExists will now receive an already-converted path on Windows
		if err := EnsureAbsoluteAndExists(AppConfig.CertsPath); err != nil {
			fmt.Printf("  Error: %v\n\n", err)
			os.Exit(1)
		} else {
			fmt.Printf("  Success: certs_path is an absolute, existing directory.\n")
		}
	}

	return nil
}
