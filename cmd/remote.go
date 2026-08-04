package cmd

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/m87/ctx/core"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func resolveRemoteAddr() string {
	if strings.TrimSpace(RemoteAddr) != "" {
		return RemoteAddr
	}
	configured := strings.TrimSpace(viper.GetString("remote"))
	if configured != "" {
		RemoteAddr = configured
	}
	return RemoteAddr
}

func remoteClient() *core.RemoteClient {
	return core.NewRemoteClient(resolveRemoteAddr(), core.DefaultRemoteTimeout)
}

func NewRemoteCmd() *cobra.Command {
	var (
		method  string
		path    string
		data    string
		timeout time.Duration
	)

	remoteCmd := &cobra.Command{
		Use:   "remote",
		Short: "Remote REST utilities",
	}

	requestCmd := &cobra.Command{
		Use:   "request",
		Short: "Send a raw HTTP request to remote",
		RunE: func(cmd *cobra.Command, args []string) error {
			remoteAddr := resolveRemoteAddr()
			if strings.TrimSpace(remoteAddr) == "" {
				return fmt.Errorf("remote address is not set; use --remote or configure remote")
			}
			client := core.NewRemoteClient(remoteAddr, timeout)
			status, body, err := client.Request(method, path, []byte(data))
			if err != nil {
				return err
			}

			payload := map[string]any{
				"status": status,
				"body":   string(body),
			}
			return printOutput(cmd, payload, func() string {
				if len(body) == 0 {
					return fmt.Sprintf("HTTP %d", status)
				}
				return fmt.Sprintf("HTTP %d\n%s", status, string(body))
			}, nil)
		},
	}

	requestCmd.Flags().StringVarP(&method, "method", "X", http.MethodGet, "HTTP method")
	requestCmd.Flags().StringVarP(&path, "path", "p", "", "Request path, e.g. /context")
	requestCmd.Flags().StringVarP(&data, "data", "d", "", "Request body")
	requestCmd.Flags().DurationVar(&timeout, "timeout", core.DefaultRemoteTimeout, "Request timeout")
	requestCmd.MarkFlagRequired("path")

	remoteCmd.AddCommand(requestCmd)
	return remoteCmd
}

func init() {
	rootCmd.AddCommand(NewRemoteCmd())
}
