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

type DayReport struct {
	Contexts  []*core.Context  `json:"contexts"`
	Intervals []*core.Interval `json:"intervals"`
}

type DayContextStats struct {
	ContextId     string  `json:"contextId"`
	Duration      int64   `json:"duration"`
	Percentage    float64 `json:"percentage"`
	IntervalCount int     `json:"intervalCount"`
}

type DayStats struct {
	Date         string                      `json:"date"`
	ContextStats []*DayContextStats          `json:"contextStats"`
	Contexts     []*core.Context             `json:"contexts"`
	Intervals    map[string][]*core.Interval `json:"intervals"`
	Distribution map[string]float64          `json:"distribution"`
}

func newHTTPClient(addr string, timeout time.Duration) *HttpClient {
	if timeout <= 0 {
		timeout = 15 * time.Second
	}
	return &HttpClient{
		baseURL: normalizeRemoteAddr(addr),
		client:  &http.Client{Timeout: timeout},
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

func (c *HttpClient) Request(method, path string, payload []byte) (int, []byte, error) {
	if strings.TrimSpace(c.baseURL) == "" {
		return 0, nil, fmt.Errorf("missing remote address")
	}

	var body io.Reader
	if len(payload) > 0 {
		body = bytes.NewReader(payload)
	}

	req, err := http.NewRequest(strings.ToUpper(strings.TrimSpace(method)), c.buildURL(path), body)
	if err != nil {
		return 0, nil, err
	}
	req.Header.Set("Accept", "application/json")
	if len(payload) > 0 {
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

func (c *HttpClient) requestJSON(method, path string, request any, response any) error {
	var payload []byte
	var err error
	if request != nil {
		payload, err = json.Marshal(request)
		if err != nil {
			return err
		}
	}

	status, body, err := c.Request(method, path, payload)
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

	if response != nil && len(body) > 0 {
		if err := json.Unmarshal(body, response); err != nil {
			return fmt.Errorf("failed to decode response: %w", err)
		}
	}

	return nil
}

func remoteCreateContext(context *core.Context) error {
	client := newHTTPClient(resolveRemoteAddr(), 15*time.Second)
	var created core.Context
	if err := client.requestJSON(http.MethodPost, "/context/", context, &created); err != nil {
		return err
	}
	if created.Id != "" {
		context.Id = created.Id
	}
	return nil
}

func remoteListContexts() ([]*core.Context, error) {
	client := newHTTPClient(resolveRemoteAddr(), 15*time.Second)
	var contexts []*core.Context
	if err := client.requestJSON(http.MethodGet, "/context/", nil, &contexts); err != nil {
		return nil, err
	}
	return contexts, nil
}

func remoteDeleteContext(id string) error {
	client := newHTTPClient(resolveRemoteAddr(), 15*time.Second)
	return client.requestJSON(http.MethodDelete, "/context/"+strings.TrimSpace(id), nil, nil)
}

func remoteUpdateContext(context *core.Context) error {
	client := newHTTPClient(resolveRemoteAddr(), 15*time.Second)
	return client.requestJSON(http.MethodPut, "/context/"+strings.TrimSpace(context.Id), context, context)
}

func remoteSwitchContext(id string, name string) error {
	client := newHTTPClient(resolveRemoteAddr(), 15*time.Second)
	payload := &core.Context{Id: strings.TrimSpace(id), Name: strings.TrimSpace(name)}
	return client.requestJSON(http.MethodPost, "/context/switch", payload, nil)
}

func remoteFreeContext() error {
	client := newHTTPClient(resolveRemoteAddr(), 15*time.Second)
	return client.requestJSON(http.MethodPost, "/context/free", nil, nil)
}

func remoteCreateInterval(interval *core.Interval) error {
	client := newHTTPClient(resolveRemoteAddr(), 15*time.Second)
	var created core.Interval
	if err := client.requestJSON(http.MethodPost, "/interval", interval, &created); err != nil {
		return err
	}
	if created.Id != "" {
		interval.Id = created.Id
	}
	return nil
}

func remoteUpdateInterval(interval *core.Interval) error {
	client := newHTTPClient(resolveRemoteAddr(), 15*time.Second)
	return client.requestJSON(http.MethodPut, "/interval/"+strings.TrimSpace(interval.Id), interval, interval)
}

func remoteDeleteInterval(id string) error {
	client := newHTTPClient(resolveRemoteAddr(), 15*time.Second)
	return client.requestJSON(http.MethodDelete, "/interval/"+strings.TrimSpace(id), nil, nil)
}

func remoteListIntervalsByDay(day string) (*DayReport, error) {
	client := newHTTPClient(resolveRemoteAddr(), 15*time.Second)
	var report DayReport
	if err := client.requestJSON(http.MethodGet, "/interval/day/"+day, nil, &report); err != nil {
		return nil, err
	}
	return &report, nil
}

func remoteSummaryDay(day string) (*DayStats, error) {
	client := newHTTPClient(resolveRemoteAddr(), 15*time.Second)
	var stats DayStats
	if err := client.requestJSON(http.MethodGet, "/interval/day/"+day+"/stats", nil, &stats); err != nil {
		return nil, err
	}
	return &stats, nil
}

func remoteListContextIntervals(contextID string) ([]*core.Interval, error) {
	client := newHTTPClient(resolveRemoteAddr(), 15*time.Second)
	var intervals []*core.Interval
	if err := client.requestJSON(http.MethodGet, "/context/"+strings.TrimSpace(contextID)+"/intervals", nil, &intervals); err != nil {
		return nil, err
	}
	return intervals, nil
}

func remoteMoveInterval(intervalID string, targetContextID string) error {
	client := newHTTPClient(resolveRemoteAddr(), 15*time.Second)
	path := "/interval/" + strings.TrimSpace(intervalID) + "/move/" + strings.TrimSpace(targetContextID)
	return client.requestJSON(http.MethodPatch, path, nil, nil)
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
			client := newHTTPClient(remoteAddr, timeout)
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
	requestCmd.Flags().DurationVar(&timeout, "timeout", 15*time.Second, "Request timeout")
	requestCmd.MarkFlagRequired("path")

	remoteCmd.AddCommand(requestCmd)
	return remoteCmd
}

func init() {
	rootCmd.AddCommand(NewRemoteCmd())
}
