package cmd

import (
	"testing"
)

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
