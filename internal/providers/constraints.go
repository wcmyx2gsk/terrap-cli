package providers

import (
	"fmt"
	"strings"

	"golang.org/x/mod/semver"
)

// Constraint represents a version constraint for a provider.
type Constraint struct {
	Operator string
	Version  string
}

// ParseConstraint parses a version constraint string like ">= 1.2.0" or "~> 2.0".
func ParseConstraint(raw string) (Constraint, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return Constraint{}, fmt.Errorf("empty constraint")
	}

	operators := []string{">=", "<=", "~>", "!=", ">", "<", "="}
	for _, op := range operators {
		if strings.HasPrefix(raw, op) {
			ver := strings.TrimSpace(raw[len(op):])
			if ver == "" {
				return Constraint{}, fmt.Errorf("missing version after operator %q", op)
			}
			return Constraint{Operator: op, Version: ver}, nil
		}
	}

	// No operator — treat as exact match
	return Constraint{Operator: "=", Version: raw}, nil
}

// Satisfies reports whether the given semver version string satisfies the constraint.
func (c Constraint) Satisfies(version string) bool {
	v := canonicalSemver(version)
	cv := canonicalSemver(c.Version)

	switch c.Operator {
	case "=":
		return semver.Compare(v, cv) == 0
	case ">=":
		return semver.Compare(v, cv) >= 0
	case "<=":
		return semver.Compare(v, cv) <= 0
	case ">":
		return semver.Compare(v, cv) > 0
	case "<":
		return semver.Compare(v, cv) < 0
	case "!=":
		return semver.Compare(v, cv) != 0
	case "~>":
		// Pessimistic constraint: >= c.Version and < next major (or minor)
		parts := strings.Split(c.Version, ".")
		if len(parts) == 2 {
			// ~> X.Y means >= X.Y.0 and < (X+1).0.0
			nextMajor := fmt.Sprintf("v%d.0.0", mustAtoi(parts[0])+1)
			return semver.Compare(v, cv) >= 0 && semver.Compare(v, nextMajor) < 0
		}
		// ~> X.Y.Z means >= X.Y.Z and < X.(Y+1).0
		if len(parts) >= 3 {
			nextMinor := fmt.Sprintf("v%s.%d.0", parts[0], mustAtoi(parts[1])+1)
			return semver.Compare(v, cv) >= 0 && semver.Compare(v, nextMinor) < 0
		}
		return false
	}
	return false
}

func canonicalSemver(v string) string {
	if !strings.HasPrefix(v, "v") {
		return "v" + v
	}
	return v
}

func mustAtoi(s string) int {
	n := 0
	fmt.Sscanf(s, "%d", &n)
	return n
}
