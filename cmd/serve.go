package cmd

import (
	"fmt"
	"net/http"

	"github.com/m87/ctx/server"
	"github.com/spf13/cobra"
)

func helloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, REST API!")
}

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Run as a server",

	Run: func(cmd *cobra.Command, args []string) {
		server.Serve()
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
}
