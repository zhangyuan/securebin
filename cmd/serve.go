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
		if err := serve.Invoke(addr); err != nil {
			log.Fatalln(err)
		}
	},
}

var addr string

func init() {
	rootCmd.AddCommand(serveCmd)
	serveCmd.Flags().StringVarP(&addr, "listening addr", "l", ":8080", "listening addr")
}
