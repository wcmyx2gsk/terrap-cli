package providers

import (
	"testing"
)

func TestParseConstraint_Valid(t *testing.T) {
	cases := []struct {
		input    string
		wantOp  string
		wantVer string
	}{
		{">= 1.2.0", ">=", "1.2.0"},
		{"~> 2.0", "~>", "2.0"},
		{"= 3.1.4", "=", "3.1.4"},
		{"!= 0.9.0", "!=", "0.9.0"},
		{"1.0.0", "=", "1.0.0"},
		{"< 5.0.0", "<", "5.0.0"},
	}
	for _, tc := range cases {
		c, err := ParseConstraint(tc.input)
		if err != nil {
			t.Fatalf("ParseConstraint(%q) error: %v", tc.input, err)
		}
		if c.Operator != tc.wantOp {
			t.Errorf("op: got %q, want %q", c.Operator, tc.wantOp)
		}
		if c.Version != tc.wantVer {
			t.Errorf("ver: got %q, want %q", c.Version, tc.wantVer)
		}
	}
}

func TestParseConstraint_Invalid(t *testing.T) {
	_, err := ParseConstraint("")
	if err == nil {
		t.Error("expected error for empty constraint")
	}
}

func TestConstraint_Satisfies(t *testing.T) {
	cases := []struct {
		constraint string
		version    string
		want       bool
	}{
		{">= 1.0.0", "1.0.0", true},
		{">= 1.0.0", "2.0.0", true},
		{">= 1.0.0", "0.9.0", false},
		{"~> 2.0", "2.5.0", true},
		{"~> 2.0", "3.0.0", false},
		{"~> 2.1.0", "2.1.9", true},
		{"~> 2.1.0", "2.2.0", false},
		{"!= 1.0.0", "1.0.1", true},
		{"!= 1.0.0", "1.0.0", false},
		{"= 3.0.0", "3.0.0", true},
		{"= 3.0.0", "3.0.1", false},
		{"< 2.0.0", "1.9.9", true},
		{"< 2.0.0", "2.0.0", false},
	}
	for _, tc := range cases {
		c, err := ParseConstraint(tc.constraint)
		if err != nil {
			t.Fatalf("parse error: %v", err)
		}
		got := c.Satisfies(tc.version)
		if got != tc.want {
			t.Errorf("Satisfies(%q, %q) = %v, want %v", tc.constraint, tc.version, got, tc.want)
		}
	}
}
