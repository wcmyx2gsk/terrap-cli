package providers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func makeCheckResult(source, current, latest string, outdated bool) VersionCheckResult {
	p := Provider{Source: source, Version: current}
	return VersionCheckResult{
		Provider:       p,
		CurrentVersion: current,
		LatestVersion:  latest,
		Outdated:       outdated,
	}
}

func TestUpgradePlanner_NoCandidates(t *testing.T) {
	results := []VersionCheckResult{
		makeCheckResult("hashicorp/aws", "4.0.0", "4.0.0", false),
		makeCheckResult("hashicorp/google", "5.1.0", "5.1.0", false),
	}
	planner := NewUpgradePlanner(results)
	assert.Empty(t, planner.Candidates())
	assert.False(t, planner.HasMajorUpgrade())
}

func TestUpgradePlanner_WithCandidates(t *testing.T) {
	results := []VersionCheckResult{
		makeCheckResult("hashicorp/aws", "4.0.0", "4.1.0", true),
		makeCheckResult("hashicorp/google", "5.1.0", "5.1.0", false),
		makeCheckResult("hashicorp/azurerm", "3.0.0", "3.2.1", true),
	}
	planner := NewUpgradePlanner(results)
	candidates := planner.Candidates()
	assert.Len(t, candidates, 2)
	// sorted by source
	assert.Equal(t, "registry.terraform.io/hashicorp/aws", candidates[0].Provider.NormalizedSource())
	assert.Equal(t, "registry.terraform.io/hashicorp/azurerm", candidates[1].Provider.NormalizedSource())
}

func TestUpgradePlanner_HasMajorUpgrade(t *testing.T) {
	results := []VersionCheckResult{
		makeCheckResult("hashicorp/aws", "3.0.0", "4.0.0", true),
	}
	planner := NewUpgradePlanner(results)
	assert.True(t, planner.HasMajorUpgrade())
}

func TestUpgradePlanner_NoMajorUpgrade(t *testing.T) {
	results := []VersionCheckResult{
		makeCheckResult("hashicorp/aws", "4.0.0", "4.5.0", true),
	}
	planner := NewUpgradePlanner(results)
	assert.False(t, planner.HasMajorUpgrade())
}

func TestUpgradeCandidate_String(t *testing.T) {
	c := UpgradeCandidate{
		Provider:       Provider{Source: "hashicorp/aws", Version: "4.0.0"},
		CurrentVersion: "4.0.0",
		LatestVersion:  "4.1.0",
	}
	assert.Equal(t, "registry.terraform.io/hashicorp/aws: 4.0.0 -> 4.1.0", c.String())
}
