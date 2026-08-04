package core

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const (
	DefaultRemoteTimeout = 15 * time.Second
	DefaultSyncLimit     = 100
)

type WorkspaceSyncClient interface {
	ListUnsyncedWorkspaces(limit int) ([]*Workspace, error)
	SyncWorkspace(workspace *Workspace) error
}

type RemoteClient struct {
	baseURL string
	client  *http.Client
}

type DayReport struct {
	Contexts  []*Context  `json:"contexts"`
	Intervals []*Interval `json:"intervals"`
}

type DayContextStats struct {
	ContextId     string  `json:"contextId"`
	Duration      int64   `json:"duration"`
	Percentage    float64 `json:"percentage"`
	IntervalCount int     `json:"intervalCount"`
}

type DayStats struct {
	Date         string                 `json:"date"`
	ContextStats []*DayContextStats     `json:"contextStats"`
	Contexts     []*Context             `json:"contexts"`
	Intervals    map[string][]*Interval `json:"intervals"`
	Distribution map[string]float64     `json:"distribution"`
}

func NewRemoteClient(addr string, timeout time.Duration) *RemoteClient {
	if timeout <= 0 {
		timeout = DefaultRemoteTimeout
	}
	return &RemoteClient{
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

func (c *RemoteClient) buildURL(path string) string {
	path = strings.TrimSpace(path)
	if path == "" {
		return c.baseURL
	}
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	return c.baseURL + path
}

func (c *RemoteClient) Request(method, path string, payload []byte) (int, []byte, error) {
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

func (c *RemoteClient) requestJSON(method, path string, request any, response any) error {
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

func (c *RemoteClient) ArchiveContext(contextID string) error {
	return c.requestJSON(http.MethodPost, "/context/"+url.PathEscape(strings.TrimSpace(contextID))+"/archive", nil, nil)
}

func (c *RemoteClient) RestoreContext(contextID string) error {
	return c.requestJSON(http.MethodPost, "/context/"+url.PathEscape(strings.TrimSpace(contextID))+"/restore", nil, nil)
}

func (c *RemoteClient) CreateWorkspace(workspace *Workspace) error {
	var created Workspace
	if err := c.requestJSON(http.MethodPost, "/workspace/", workspace, &created); err != nil {
		return err
	}
	if created.Id != "" {
		workspace.Id = created.Id
	}
	return nil
}

func (c *RemoteClient) ListWorkspaces() ([]*Workspace, error) {
	var workspaces []*Workspace
	if err := c.requestJSON(http.MethodGet, "/workspace/", nil, &workspaces); err != nil {
		return nil, err
	}
	return workspaces, nil
}

func (c *RemoteClient) ListUnsyncedWorkspaces(limit int) ([]*Workspace, error) {
	if limit <= 0 {
		limit = DefaultSyncLimit
	}
	path := "/api/sync/workspace/?limit=" + strconv.Itoa(limit)
	var workspaces []*Workspace
	if err := c.requestJSON(http.MethodGet, path, nil, &workspaces); err != nil {
		return nil, err
	}
	return workspaces, nil
}

func (c *RemoteClient) SyncWorkspace(workspace *Workspace) error {
	if workspace == nil {
		return fmt.Errorf("workspace is required")
	}
	var synced Workspace
	if err := c.requestJSON(http.MethodPost, "/api/sync/workspace/", workspace, &synced); err != nil {
		return err
	}
	*workspace = synced
	return nil
}

func (c *RemoteClient) GetWorkspace(id string) (*Workspace, error) {
	var workspace Workspace
	if err := c.requestJSON(http.MethodGet, "/workspace/"+url.PathEscape(strings.TrimSpace(id)), nil, &workspace); err != nil {
		return nil, err
	}
	return &workspace, nil
}

func (c *RemoteClient) DeleteWorkspace(id string) error {
	return c.requestJSON(http.MethodDelete, "/workspace/"+strings.TrimSpace(id), nil, nil)
}

func (c *RemoteClient) UpdateWorkspace(workspace *Workspace) error {
	return c.requestJSON(http.MethodPut, "/workspace/"+strings.TrimSpace(workspace.Id), workspace, workspace)
}

func (c *RemoteClient) CreateContext(context *Context) error {
	var created Context
	if err := c.requestJSON(http.MethodPost, "/context/", context, &created); err != nil {
		return err
	}
	if created.Id != "" {
		context.Id = created.Id
	}
	return nil
}

func (c *RemoteClient) ListContexts(workspaceID string) ([]*Context, error) {
	var contexts []*Context
	if err := c.requestJSON(http.MethodGet, remoteListContextsPath(workspaceID), nil, &contexts); err != nil {
		return nil, err
	}
	return contexts, nil
}

func (c *RemoteClient) DeleteContext(id string) error {
	return c.requestJSON(http.MethodDelete, "/context/"+strings.TrimSpace(id), nil, nil)
}

func (c *RemoteClient) UpdateContext(context *Context) error {
	return c.requestJSON(http.MethodPut, "/context/"+strings.TrimSpace(context.Id), context, context)
}

func (c *RemoteClient) SwitchContext(id string, name string, workspaceID string) error {
	payload := &Context{
		Id:          strings.TrimSpace(id),
		Name:        strings.TrimSpace(name),
		WorkspaceId: strings.TrimSpace(workspaceID),
	}
	return c.requestJSON(http.MethodPost, "/context/switch", payload, nil)
}

func (c *RemoteClient) FreeContext() error {
	return c.requestJSON(http.MethodPost, "/context/free", nil, nil)
}

func (c *RemoteClient) CreateInterval(interval *Interval) error {
	var created Interval
	if err := c.requestJSON(http.MethodPost, "/interval", interval, &created); err != nil {
		return err
	}
	if created.Id != "" {
		interval.Id = created.Id
	}
	return nil
}

func (c *RemoteClient) UpdateInterval(interval *Interval) error {
	return c.requestJSON(http.MethodPut, "/interval/"+strings.TrimSpace(interval.Id), interval, interval)
}

func (c *RemoteClient) DeleteInterval(id string) error {
	return c.requestJSON(http.MethodDelete, "/interval/"+strings.TrimSpace(id), nil, nil)
}

func (c *RemoteClient) ListIntervalsByDay(day string, workspaceID string) (*DayReport, error) {
	var report DayReport
	if err := c.requestJSON(http.MethodGet, remoteListIntervalsByDayPath(day, workspaceID), nil, &report); err != nil {
		return nil, err
	}
	return &report, nil
}

func (c *RemoteClient) SummaryDay(day string, workspaceID string) (*DayStats, error) {
	var stats DayStats
	if err := c.requestJSON(http.MethodGet, remoteSummaryDayPath(day, workspaceID), nil, &stats); err != nil {
		return nil, err
	}
	return &stats, nil
}

func remoteListContextsPath(workspaceID string) string {
	return "/context/?workspaceId=" + url.QueryEscape(strings.TrimSpace(workspaceID))
}

func remoteListIntervalsByDayPath(day string, workspaceID string) string {
	return "/interval/day/" + url.PathEscape(strings.TrimSpace(day)) +
		"?workspaceId=" + url.QueryEscape(strings.TrimSpace(workspaceID))
}

func remoteSummaryDayPath(day string, workspaceID string) string {
	return "/interval/day/" + url.PathEscape(strings.TrimSpace(day)) +
		"/stats?workspaceId=" + url.QueryEscape(strings.TrimSpace(workspaceID))
}

func (c *RemoteClient) ListContextIntervals(contextID string) ([]*Interval, error) {
	var intervals []*Interval
	if err := c.requestJSON(http.MethodGet, "/context/"+strings.TrimSpace(contextID)+"/intervals", nil, &intervals); err != nil {
		return nil, err
	}
	return intervals, nil
}

func (c *RemoteClient) MoveInterval(intervalID string, targetContextID string) error {
	path := "/interval/" + strings.TrimSpace(intervalID) + "/move/" + strings.TrimSpace(targetContextID)
	return c.requestJSON(http.MethodPatch, path, nil, nil)
}
