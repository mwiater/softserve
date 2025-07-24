// cmd/softserve/root.go
package cmd

import (
	"fmt"

	"github.com/mwiater/softserve"
	"github.com/spf13/cobra"
)

var (
	webRoot   string
	ssl       bool
	httpPort  int
	httpsPort int
	logLevel  string
	api       bool
	apiPrefix string
)

// rootCmd is the base command for the softserve application.
// It provides a local static file server with optional API mocking and live-reload for development workflows.
// This command serves as the CLI entry point and may include additional subcommands for listing and diagnostics.
var rootCmd = &cobra.Command{
	Use:   "softserve",
	Short: "Local static server with live reload and API mocking",
	Long: `Softserve is a lightweight local static file server tailored for frontend development.
It supports automatic browser reloads on file changes, static API mocking, and optional HTTPS serving.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		cfg := softserve.Config{
			WebRoot:   webRoot,
			SSL:       ssl,
			HTTPPort:  httpPort,
			HTTPSPort: httpsPort,
			LogLevel:  logLevel,
			API:       api,
			APIPrefix: apiPrefix,
		}
		softserve.InitConfig(cfg)
		if cfg.API {
			if err := softserve.LoadAPIResponses(); err != nil {
				return fmt.Errorf("failed to load api.yaml: %w", err)
			}
		}
		fmt.Printf("📂 Web root: %s\n", cfg.WebRoot)
		return nil
	},
}

// Execute runs the root command and handles any errors that occur during execution.
func Execute() error {
	if err := rootCmd.Execute(); err != nil {
		return fmt.Errorf("failed to execute rootCmd: %w", err)
	}

	return nil
}

// init initializes the root command's configuration.
func init() {
	rootCmd.PersistentFlags().StringVar(&webRoot, "web-root", "examples/basic", "directory to serve")
	rootCmd.PersistentFlags().BoolVar(&ssl, "ssl", false, "enable HTTPS")
	rootCmd.PersistentFlags().IntVar(&httpPort, "http-port", 8080, "HTTP port")
	rootCmd.PersistentFlags().IntVar(&httpsPort, "https-port", 8443, "HTTPS port")
	rootCmd.PersistentFlags().StringVar(&logLevel, "log-level", "info", "log level")
	rootCmd.PersistentFlags().BoolVar(&api, "api", false, "enable API mocking")
	rootCmd.PersistentFlags().StringVar(&apiPrefix, "api-prefix", "/api/", "API prefix")
}
