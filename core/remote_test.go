package core

import (
	"io"
	"net/http"
	"strings"
	"testing"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(request *http.Request) (*http.Response, error) {
	return f(request)
}

func TestRemoteWorkspaceScopedRequestsIncludeWorkspaceID(t *testing.T) {
	tests := map[string]struct {
		got  string
		want string
	}{
		"contexts": {
			got:  remoteListContextsPath("workspace with spaces"),
			want: "/context/?workspaceId=workspace+with+spaces",
		},
		"intervals": {
			got:  remoteListIntervalsByDayPath("2026-06-15", "workspace with spaces"),
			want: "/interval/day/2026-06-15?workspaceId=workspace+with+spaces",
		},
		"summary": {
			got:  remoteSummaryDayPath("2026-06-15", "workspace with spaces"),
			want: "/interval/day/2026-06-15/stats?workspaceId=workspace+with+spaces",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			if test.got != test.want {
				t.Fatalf("expected %q, got %q", test.want, test.got)
			}
		})
	}
}

func TestRemoteClientListContexts(t *testing.T) {
	client := &RemoteClient{
		baseURL: "http://remote",
		client: &http.Client{Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
			if r.Method != http.MethodGet {
				t.Errorf("expected method %q, got %q", http.MethodGet, r.Method)
			}
			if r.URL.Path != "/context/" {
				t.Errorf("expected path %q, got %q", "/context/", r.URL.Path)
			}
			if got := r.URL.Query().Get("workspaceId"); got != "workspace with spaces" {
				t.Errorf("expected workspace ID %q, got %q", "workspace with spaces", got)
			}
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader(`[{"id":"context-id","name":"Context"}]`)),
				Header:     make(http.Header),
			}, nil
		})},
	}

	contexts, err := client.ListContexts("workspace with spaces")
	if err != nil {
		t.Fatalf("list contexts: %v", err)
	}
	if len(contexts) != 1 || contexts[0].Id != "context-id" {
		t.Fatalf("unexpected contexts: %#v", contexts)
	}
}

func TestRemoteClientReturnsResponseError(t *testing.T) {
	client := &RemoteClient{
		baseURL: "http://remote",
		client: &http.Client{Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusNotFound,
				Body:       io.NopCloser(strings.NewReader("workspace not found\n")),
				Header:     make(http.Header),
			}, nil
		})},
	}
	_, err := client.GetWorkspace("missing")
	if err == nil || !strings.Contains(err.Error(), "remote request failed (404): workspace not found") {
		t.Fatalf("expected remote response error, got %v", err)
	}
}
