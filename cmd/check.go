package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/sirrend/terrap-cli/internal/providers"
	"github.com/spf13/cobra"
)

var checkCmd = &cobra.Command{
	Use:   "check <provider> <version> <constraint>",
	Short: "Check whether a provider version satisfies a given constraint",
	Args:  cobra.ExactArgs(3),
	RunE: func(cmd *cobra.Command, args []string) error {
		return runCheck(cmd, args)
	},
}

type checkResult struct {
	Provider   string `json:"provider"`
	Version    string `json:"version"`
	Constraint string `json:"constraint"`
	Satisfies  bool   `json:"satisfies"`
}

func runCheck(cmd *cobra.Command, args []string) error {
	providerName := args[0]
	version := args[1]
	rawConstraint := args[2]

	c, err := providers.ParseConstraint(rawConstraint)
	if err != nil {
		return fmt.Errorf("invalid constraint %q: %w", rawConstraint, err)
	}

	result := checkResult{
		Provider:   providerName,
		Version:    version,
		Constraint: rawConstraint,
		Satisfies:  c.Satisfies(version),
	}

	jsonOut, _ := cmd.Flags().GetBool("json")
	if jsonOut {
		return printCheckJSON(result)
	}
	return printCheckText(result)
}

func printCheckJSON(r checkResult) error {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(r)
}

func printCheckText(r checkResult) error {
	status := "✗ does not satisfy"
	if r.Satisfies {
		status = "✓ satisfies"
	}
	fmt.Printf("Provider : %s\n", r.Provider)
	fmt.Printf("Version  : %s\n", r.Version)
	fmt.Printf("Constraint: %s\n", r.Constraint)
	fmt.Printf("Result   : %s\n", status)
	return nil
}

func init() {
	checkCmd.Flags().Bool("json", false, "Output result as JSON")
	rootCmd.AddCommand(checkCmd)
}
