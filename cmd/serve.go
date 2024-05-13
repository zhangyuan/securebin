package cmd

import (
	"log"
	"securebin/pkg/serve"

	"github.com/spf13/cobra"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Serve the API requests",
	Run: func(cmd *cobra.Command, args []string) {
		if err := serve.Invoke(); err != nil {
			log.Fatalln(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
}
