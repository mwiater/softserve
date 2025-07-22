// cmd/softserve/root.go
package cmd

import (
	"fmt"
	"os"

	"github.com/mwiater/softserve"
	"github.com/spf13/cobra"
)

// rootCmd is the base command for the softserve application.
// It provides a local static file server with optional API mocking and live-reload for development workflows.
// This command serves as the CLI entry point and may include additional subcommands for listing and diagnostics.
var rootCmd = &cobra.Command{
	Use:   "softserve",
	Short: "Local static server with live reload and API mocking",
	Long: `Softserve is a lightweight local static file server tailored for frontend development.
It supports automatic browser reloads on file changes, static API mocking, and optional HTTPS serving.`,
}

// Execute executes the root command along with any registered subcommands.
// If the command execution results in an error, the error is printed and the program exits
// with a non-zero status code.
func Execute() error {
	if err := rootCmd.Execute(); err != nil {
		return fmt.Errorf("failed to execute rootCmd: %w", err)
	}

	return nil
}

// init initializes the root command's configuration.
// This function is reserved for setting up additional top-level flags or configurations.
// Since subcommands are self-registered in their respective files, no further initialization
// is required here.
func init() {
	if err := softserve.LoadConfig(); err != nil {
		fmt.Printf("failed to load config: %w", err)
		os.Exit(1)
	}
	fmt.Println("✅ Config loaded successfully")
	fmt.Printf("📂 Web root: %s\n", softserve.AppConfig.WebRoot)
}
