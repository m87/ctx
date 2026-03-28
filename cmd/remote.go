package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/m87/ctx/core"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type HttpClient struct {
	baseURL string
	client  *http.Client
}

func newHTTPClient(addr string, timeout time.Duration) *HttpClient {
	if timeout <= 0 {
		timeout = 15 * time.Second
	}

	return &HttpClient{
		baseURL: normalizeRemoteAddr(addr),
		client: &http.Client{
			Timeout: timeout,
		},
	}
}

func normalizeRemoteAddr(addr string) string {
	trimmed := strings.TrimSpace(addr)
	if trimmed == "" {
		return ""
	}

	if !strings.HasPrefix(trimmed, "http://") && !strings.HasPrefix(trimmed, "https://") {
		trimmed = "http://" + trimmed
	}

	return strings.TrimRight(trimmed, "/")
}

func (c *HttpClient) buildURL(path string) string {
	path = strings.TrimSpace(path)
	if path == "" {
		return c.baseURL
	}
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	return c.baseURL + path
}

func (c *HttpClient) Request(method, path string, data []byte) (int, []byte, error) {
	if strings.TrimSpace(c.baseURL) == "" {
		return 0, nil, fmt.Errorf("missing remote address")
	}

	var body io.Reader
	if len(data) > 0 {
		body = bytes.NewReader(data)
	}

	method = strings.ToUpper(strings.TrimSpace(method))
	if method == "" {
		method = http.MethodGet
	}

	req, err := http.NewRequest(method, c.buildURL(path), body)
	if err != nil {
		return 0, nil, err
	}
	req.Header.Set("Accept", "application/json")
	if len(data) > 0 {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return 0, nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return resp.StatusCode, nil, err
	}

	return resp.StatusCode, respBody, nil
}

func (c *HttpClient) CreateContext(context *core.Context) error {
	payload, err := json.Marshal(context)
	if err != nil {
		return err
	}

	status, body, err := c.Request(http.MethodPost, "/context", payload)
	if err != nil {
		return err
	}
	if status >= http.StatusBadRequest {
		bodyText := strings.TrimSpace(string(body))
		if bodyText == "" {
			bodyText = http.StatusText(status)
		}
		return fmt.Errorf("remote request failed (%d): %s", status, bodyText)
	}

	var created core.Context
	if len(body) > 0 {
		if err := json.Unmarshal(body, &created); err == nil && created.Id != "" {
			context.Id = created.Id
		}
	}

	return nil
}

func remoteCreateContext(context *core.Context) error {
	httpClient := newHTTPClient(resolveRemoteAddr(), 15*time.Second)
	return httpClient.CreateContext(context)
}

func remoteListContexts(cmd *cobra.Command) ([]*core.Context, error) {
	httpClient := newHTTPClient(resolveRemoteAddr(), 15*time.Second)
	return httpClient.ListContexts()
}

func (c *HttpClient) ListContexts() ([]*core.Context, error) {
	status, body, err := c.Request(http.MethodGet, "/context/", nil)
	if err != nil {
		return nil, err
	}
	if status >= http.StatusBadRequest {
		bodyText := strings.TrimSpace(string(body))
		if bodyText == "" {
			bodyText = http.StatusText(status)
		}
		return nil, fmt.Errorf("remote request failed (%d): %s", status, bodyText)
	}

	var contexts []*core.Context
	if len(body) > 0 {
		if err := json.Unmarshal(body, &contexts); err != nil {
			return nil, fmt.Errorf("failed to parse response: %w", err)
		}
	}

	return contexts, nil
}

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
		Short: "Send an HTTP request to the configured remote",
		RunE: func(cmd *cobra.Command, args []string) error {
			remoteAddr := resolveRemoteAddr()
			if strings.TrimSpace(remoteAddr) == "" {
				return fmt.Errorf("remote address is not set; use --remote or set remote in config")
			}

			httpClient := newHTTPClient(remoteAddr, timeout)
			status, body, err := httpClient.Request(method, path, []byte(data))
			if err != nil {
				return err
			}

			cmd.Printf("HTTP %d\n", status)
			if len(body) > 0 {
				cmd.Println(string(body))
			}

			if status >= http.StatusBadRequest {
				return fmt.Errorf("remote request failed with status %d", status)
			}

			return nil
		},
	}

	listContextsCmd := &cobra.Command{
		Use:   "list-contexts",
		Short: "List contexts from remote",
		RunE: func(cmd *cobra.Command, args []string) error {
			remoteAddr := resolveRemoteAddr()
			if strings.TrimSpace(remoteAddr) == "" {
				return fmt.Errorf("remote address is not set; use --remote or set remote in config")
			}

			httpClient := newHTTPClient(remoteAddr, timeout)
			status, body, err := httpClient.Request(http.MethodGet, "/context", nil)
			if err != nil {
				return err
			}

			cmd.Printf("HTTP %d\n", status)
			if len(body) > 0 {
				cmd.Println(string(body))
			}

			if status >= http.StatusBadRequest {
				return fmt.Errorf("remote request failed with status %d", status)
			}

			return nil
		},
	}

	var (
		switchID   string
		switchName string
	)
	switchCmd := &cobra.Command{
		Use:   "switch",
		Short: "Switch active context on remote",
		RunE: func(cmd *cobra.Command, args []string) error {
			remoteAddr := resolveRemoteAddr()
			if strings.TrimSpace(remoteAddr) == "" {
				return fmt.Errorf("remote address is not set; use --remote or set remote in config")
			}
			if strings.TrimSpace(switchID) == "" && strings.TrimSpace(switchName) == "" {
				return fmt.Errorf("provide --id or --name")
			}

			payloadMap := map[string]string{}
			if strings.TrimSpace(switchID) != "" {
				payloadMap["id"] = strings.TrimSpace(switchID)
			}
			if strings.TrimSpace(switchName) != "" {
				payloadMap["name"] = strings.TrimSpace(switchName)
			}

			payload, err := json.Marshal(payloadMap)
			if err != nil {
				return err
			}

			httpClient := newHTTPClient(remoteAddr, timeout)
			status, body, err := httpClient.Request(http.MethodPost, "/context/switch", payload)
			if err != nil {
				return err
			}

			cmd.Printf("HTTP %d\n", status)
			if len(body) > 0 {
				cmd.Println(string(body))
			}

			if status >= http.StatusBadRequest {
				return fmt.Errorf("remote request failed with status %d", status)
			}

			return nil
		},
	}

	freeCmd := &cobra.Command{
		Use:   "free",
		Short: "Free active context on remote",
		RunE: func(cmd *cobra.Command, args []string) error {
			remoteAddr := resolveRemoteAddr()
			if strings.TrimSpace(remoteAddr) == "" {
				return fmt.Errorf("remote address is not set; use --remote or set remote in config")
			}

			httpClient := newHTTPClient(remoteAddr, timeout)
			status, body, err := httpClient.Request(http.MethodPost, "/context/free", nil)
			if err != nil {
				return err
			}

			cmd.Printf("HTTP %d\n", status)
			if len(body) > 0 {
				cmd.Println(string(body))
			}

			if status >= http.StatusBadRequest {
				return fmt.Errorf("remote request failed with status %d", status)
			}

			return nil
		},
	}

	requestCmd.Flags().StringVarP(&method, "method", "X", http.MethodGet, "HTTP method")
	requestCmd.Flags().StringVarP(&path, "path", "p", "", "Request path, e.g. /context")
	requestCmd.Flags().StringVarP(&data, "data", "d", "", "JSON request body")
	requestCmd.Flags().DurationVar(&timeout, "timeout", 15*time.Second, "Request timeout, e.g. 10s")
	requestCmd.MarkFlagRequired("path")
	switchCmd.Flags().StringVar(&switchID, "id", "", "Context ID")
	switchCmd.Flags().StringVarP(&switchName, "name", "n", "", "Context name (used when creating/switching)")

	remoteCmd.AddCommand(requestCmd)
	remoteCmd.AddCommand(listContextsCmd)
	remoteCmd.AddCommand(switchCmd)
	remoteCmd.AddCommand(freeCmd)
	return remoteCmd
}

func init() {
	rootCmd.AddCommand(NewRemoteCmd())
}
