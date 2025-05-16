package cmd

import (
	"fmt"
	"log"
	"net/http"

	"github.com/spf13/cobra"
)

func helloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, REST API!")
}

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Run as a server",

	Run: func(cmd *cobra.Command, args []string) {
		http.HandleFunc("/api/hello", helloHandler)
		log.Fatal(http.ListenAndServe(":8080", nil))
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
}
