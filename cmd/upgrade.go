package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/sirrend/terrap-cli/internal/providers"
	"github.com/spf13/cobra"
)

var upgradeJSONOutput bool

var upgradeCmd = &cobra.Command{
	Use:   "upgrade",
	Short: "Show available provider upgrades",
	Long:  `Checks all configured providers and lists those with newer versions available.`,
	RunE:  runUpgrade,
}

func runUpgrade(cmd *cobra.Command, args []string) error {
	initData, err := loadInitData()
	if err != nil {
		return fmt.Errorf("failed to load init data: %w", err)
	}

	regClient := providers.NewRegistryClient()
	checker := providers.NewVersionChecker(regClient)

	var results []providers.VersionCheckResult
	for _, p := range initData.Providers {
		r, err := checker.Check(p)
		if err != nil {
			// Print warnings to stderr so they don't pollute JSON output
			fmt.Fprintf(os.Stderr, "warning: could not check %s: %v\n", p.NormalizedSource(), err)
			continue
		}
		results = append(results, r)
	}

	planner := providers.NewUpgradePlanner(results)
	candidates := planner.Candidates()

	if upgradeJSONOutput {
		return printUpgradeJSON(candidates)
	}
	return printUpgradeText(candidates, planner.HasMajorUpgrade())
}

func printUpgradeJSON(candidates []providers.UpgradeCandidate) error {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(candidates)
}

func printUpgradeText(candidates []providers.UpgradeCandidate, hasMajor bool) error {
	if len(candidates) == 0 {
		fmt.Println("All providers are up to date.")
		return nil
	}
	fmt.Printf("Found %d provider(s) with available upgrades:\n\n", len(candidates))
	for _, c := range candidates {
		fmt.Printf("  %s\n", c.String())
	}
	if hasMajor {
		// Use a more visible separator before the major version warning
		fmt.Println()
		fmt.Println("===========================================")
		fmt.Println("⚠  One or more upgrades include a major version bump. Review breaking changes before upgrading.")
		fmt.Println("===========================================")
	}
	return nil
}

func init() {
	upgradeCmd.Flags().BoolVar(&upgradeJSONOutput, "json", false, "Output results as JSON")
	rootCmd.AddCommand(upgradeCmd)
}
