package providers

import (
	"fmt"
	"sort"

	"golang.org/x/mod/semver"
)

// UpgradeCandidate represents a provider that can be upgraded.
type UpgradeCandidate struct {
	Provider       Provider
	CurrentVersion string
	LatestVersion  string
}

// String returns a human-readable representation of the upgrade candidate.
func (u UpgradeCandidate) String() string {
	return fmt.Sprintf("%s: %s -> %s", u.Provider.NormalizedSource(), u.CurrentVersion, u.LatestVersion)
}

// UpgradePlanner collects upgrade candidates from version check results.
type UpgradePlanner struct {
	results []VersionCheckResult
}

// NewUpgradePlanner creates a new UpgradePlanner from a slice of check results.
func NewUpgradePlanner(results []VersionCheckResult) *UpgradePlanner {
	return &UpgradePlanner{results: results}
}

// Candidates returns all providers that have a newer version available.
func (p *UpgradePlanner) Candidates() []UpgradeCandidate {
	var candidates []UpgradeCandidate
	for _, r := range p.results {
		if r.Outdated {
			candidates = append(candidates, UpgradeCandidate{
				Provider:       r.Provider,
				CurrentVersion: r.CurrentVersion,
				LatestVersion:  r.LatestVersion,
			})
		}
	}
	sort.Slice(candidates, func(i, j int) bool {
		return candidates[i].Provider.NormalizedSource() < candidates[j].Provider.NormalizedSource()
	})
	return candidates
}

// HasMajorUpgrade returns true if any candidate involves a major version bump.
func (p *UpgradePlanner) HasMajorUpgrade() bool {
	for _, c := range p.Candidates() {
		cur := "v" + c.CurrentVersion
		latest := "v" + c.LatestVersion
		if semver.IsValid(cur) && semver.IsValid(latest) {
			if semver.Major(cur) != semver.Major(latest) {
				return true
			}
		}
	}
	return false
}
