/*
Copyright © 2025 Alberto Cerato <macros123@gmail.com>
*/
package cmd

import (
	"fmt"

	"github.com/cli/go-gh/v2/pkg/repository"
	"github.com/spf13/cobra"

	"github.com/alberto-cerato/gh-branch-protection/internal/github"
)

func init() {
	rootCmd.AddCommand(listCmd)
}

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List branch protection rules",
	RunE: func(cmd *cobra.Command, args []string) error {
		currentRepo, err := repository.Current()
		if err != nil {
			return fmt.Errorf("Cannot list protected branches: %w\n", err)
			
		}

		branches, err := github.ListBranchProtectionRules(currentRepo.Owner, currentRepo.Name)
		if err != nil {
			return fmt.Errorf("Cannot list protected branches: %w\n", err)
		}
		for _, b := range branches {
			fmt.Printf("{id: %s, pattern: %s}\n", b.ID, b.Pattern)
		}

		return nil
	},
}
