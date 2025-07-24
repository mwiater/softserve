// cmd/softserve/serve.go

package cmd

import (
	"github.com/mwiater/softserve"
	"github.com/spf13/cobra"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "serve",
	Long:  `serve`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return softserve.StartServer()
	},
}

// init adds the list command to the root command.
func init() {
	rootCmd.AddCommand(serveCmd)
}
