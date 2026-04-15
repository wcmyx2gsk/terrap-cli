package providers

import (
	"fmt"

	"github.com/Masterminds/semver/v3"
)

// VersionCheckResult holds the outcome of a version check for a single provider.
type VersionCheckResult struct {
	Provider      Provider
	LatestVersion string
	IsOutdated    bool
	Error         error
}

// VersionChecker checks whether providers are up to date.
type VersionChecker struct {
	registry *RegistryClient
}

// NewVersionChecker creates a new VersionChecker.
func NewVersionChecker() *VersionChecker {
	return &VersionChecker{registry: NewRegistryClient()}
}

// CheckProvider fetches the latest version and compares it to the current constraint.
func (vc *VersionChecker) CheckProvider(p Provider) VersionCheckResult {
	result := VersionCheckResult{Provider: p}

	latest, err := vc.registry.GetLatestVersion(p.NormalizedSource())
	if err != nil {
		result.Error = err
		return result
	}
	result.LatestVersion = latest

	if p.Version == "" {
		result.IsOutdated = false
		return result
	}

	constraint, err := semver.NewConstraint(p.Version)
	if err != nil {
		result.Error = fmt.Errorf("invalid version constraint %q: %w", p.Version, err)
		return result
	}

	latestSemver, err := semver.NewVersion(latest)
	if err != nil {
		result.Error = fmt.Errorf("invalid latest version %q: %w", latest, err)
		return result
	}

	result.IsOutdated = !constraint.Check(latestSemver)
	return result
}

// CheckAll runs CheckProvider for every provider in the list.
func (vc *VersionChecker) CheckAll(providers ProviderList) []VersionCheckResult {
	results := make([]VersionCheckResult, 0, len(providers))
	for _, p := range providers {
		results = append(results, vc.CheckProvider(p))
	}
	return results
}
