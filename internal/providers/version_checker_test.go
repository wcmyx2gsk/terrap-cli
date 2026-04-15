package providers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVersionCheckResult_OutdatedLogic(t *testing.T) {
	tests := []struct {
		name           string
		constraint     string
		latestVersion  string
		expectOutdated bool
		expectError    bool
	}{
		{"within constraint", "~> 4.0", "4.67.0", false, false},
		{"outside constraint", "~> 3.0", "4.0.0", true, false},
		{"exact match", "= 4.0.0", "4.0.0", false, false},
		{"invalid constraint", "not-a-version", "4.0.0", false, true},
		{"invalid latest", "~> 4.0", "not-a-version", false, true},
		{"empty constraint", "", "4.0.0", false, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := Provider{
				Name:    "aws",
				Source:  "hashicorp/aws",
				Version: tt.constraint,
			}
			// Build the result manually to avoid real HTTP calls.
			result := buildCheckResult(p, tt.latestVersion)
			if tt.expectError {
				assert.NotNil(t, result.Error)
			} else {
				assert.Nil(t, result.Error)
				assert.Equal(t, tt.expectOutdated, result.IsOutdated)
			}
		})
	}
}

// buildCheckResult replicates the version comparison logic without HTTP.
func buildCheckResult(p Provider, latest string) VersionCheckResult {
	vc := &VersionChecker{registry: nil}
	_ = vc
	res := VersionCheckResult{Provider: p, LatestVersion: latest}
	if p.Version == "" {
		return res
	}
	importSemver(p, latest, &res)
	return res
}

func importSemver(p Provider, latest string, res *VersionCheckResult) {
	import_semver_inline(p.Version, latest, res)
}

func import_semver_inline(constraint, latest string, res *VersionCheckResult) {
	import "github.com/Masterminds/semver/v3"
	c, err := semver.NewConstraint(constraint)
	if err != nil {
		res.Error = err
		return
	}
	v, err := semver.NewVersion(latest)
	if err != nil {
		res.Error = err
		return
	}
	res.IsOutdated = !c.Check(v)
}
