package cmd

import (
	"log"

	"github.com/m87/ctx/bootstrap"
	"github.com/m87/ctx/core"
	"github.com/m87/ctx/server"
	"github.com/spf13/cobra"
)

func newServeCmd(manager *core.ContextManager) *cobra.Command {
	return &cobra.Command{
		Use: "serve",
		Short: "Start server",
		Long: `Start server to handle context management requests.
For example:
		ctx serve --port 8080
`,
		Run: func(cmd *cobra.Command, ars []string) {
			port, _ := cmd.Flags().GetString("port")
			srv := server.New(manager)
			log.Fatal(srv.Listen(":" + port))
		},
	}
}

func init() {
	serveCmd := newServeCmd(bootstrap.CreateManager())
	serveCmd.Flags().StringP("port", "p", "8080", "server port")
	rootCmd.AddCommand(serveCmd)
}
