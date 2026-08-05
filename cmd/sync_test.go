package cmd

import (
	"testing"

	"github.com/m87/ctx/core"
	"github.com/stretchr/testify/require"
)

func TestFormatSyncProgressUsesDirectionOfTransfer(t *testing.T) {
	tests := map[string]struct {
		direction core.SyncDirection
		current   int
		want      string
	}{
		"upload moves right": {
			direction: core.SyncDirectionUpload,
			current:   5,
			want:      "Uploading contexts   [==========>.........] 5/10",
		},
		"download moves left": {
			direction: core.SyncDirectionDownload,
			current:   5,
			want:      "Downloading contexts   [.........<==========] 5/10",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			got := formatSyncProgress(core.SyncProgress{
				Direction: test.direction,
				Resource:  "contexts",
				Current:   test.current,
				Total:     10,
			})

			require.Equal(t, test.want, got)
		})
	}
}

func TestFormatSyncProgressUsesFullBarWhenDownloadCompletes(t *testing.T) {
	got := formatSyncProgress(core.SyncProgress{
		Direction: core.SyncDirectionDownload,
		Resource:  "intervals",
		Current:   10,
		Total:     10,
	})

	require.Equal(t, "Downloading intervals  [====================] 10/10", got)
}
