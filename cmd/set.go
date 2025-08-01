/*
Copyright © 2025 Alberto Cerato <macros123@gmail.com>
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"

	"github.com/cli/go-gh/v2/pkg/repository"
	"github.com/spf13/cobra"

	"github.com/alberto-cerato/gh-branch-protection/internal/github"
)

func init() {
	rootCmd.AddCommand(setCmd)
}

// setCmd represents the set command
var setCmd = &cobra.Command{
	Use:   "set <branch>",
	Short: "Set the protection for a branch",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		branch := args[0]

		currentRepo, err := repository.Current()
		if err != nil {
			return fmt.Errorf("Cannot set branch protection configuration: %w", err)
		}

		var rule github.BranchProtectionRule
		err = json.NewDecoder(os.Stdin).Decode(&rule)
		if err != nil {
			return fmt.Errorf("Cannot set branch protection configuration: %w", err)
		}
		slog.Debug("setCmd", "operation", "decode JSON input", "rule", rule)

		repoId, err := github.GetRepoID(currentRepo.Owner, currentRepo.Name)
		if err != nil {
			return fmt.Errorf("Cannot set branch protection configuration: %w", err)

		}
		if err := github.CreateBranchProtectionRule(repoId, branch, rule); err != nil {
			return fmt.Errorf("Cannot set branch protection configuration: %w", err)

		}
		return nil
	},
}
