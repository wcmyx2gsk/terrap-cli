package cmd

import (
	"bytes"
	"io"
	"os"
	"testing"

	"github.com/sirrend/terrap-cli/internal/providers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func captureStdout(fn func()) string {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	fn()
	w.Close()
	os.Stdout = old
	var buf bytes.Buffer
	io.Copy(&buf, r)
	return buf.String()
}

func TestPrintUpgradeText_Empty(t *testing.T) {
	out := captureStdout(func() {
		err := printUpgradeText(nil, false)
		require.NoError(t, err)
	})
	assert.Contains(t, out, "All providers are up to date.")
}

func TestPrintUpgradeText_WithCandidates(t *testing.T) {
	candidates := []providers.UpgradeCandidate{
		{
			Provider:       providers.Provider{Source: "hashicorp/aws", Version: "4.0.0"},
			CurrentVersion: "4.0.0",
			LatestVersion:  "4.1.0",
		},
	}
	out := captureStdout(func() {
		err := printUpgradeText(candidates, false)
		require.NoError(t, err)
	})
	assert.Contains(t, out, "1 provider(s)")
	assert.Contains(t, out, "4.0.0 -> 4.1.0")
	assert.NotContains(t, out, "major version")
}

func TestPrintUpgradeText_MajorWarning(t *testing.T) {
	candidates := []providers.UpgradeCandidate{
		{
			Provider:       providers.Provider{Source: "hashicorp/aws", Version: "3.0.0"},
			CurrentVersion: "3.0.0",
			LatestVersion:  "4.0.0",
		},
	}
	out := captureStdout(func() {
		err := printUpgradeText(candidates, true)
		require.NoError(t, err)
	})
	assert.Contains(t, out, "major version bump")
}

func TestPrintUpgradeJSON(t *testing.T) {
	candidates := []providers.UpgradeCandidate{
		{
			Provider:       providers.Provider{Source: "hashicorp/aws", Version: "4.0.0"},
			CurrentVersion: "4.0.0",
			LatestVersion:  "4.2.0",
		},
	}
	out := captureStdout(func() {
		err := printUpgradeJSON(candidates)
		require.NoError(t, err)
	})
	assert.Contains(t, out, "4.2.0")
	assert.Contains(t, out, "{")
}
