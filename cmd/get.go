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
	rootCmd.AddCommand(getCmd)
}

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get <branch>",
	Short: "Get the branch protection of a branch",
	RunE: func(cmd *cobra.Command, args []string) error {
		currentRepo, err := repository.Current()
		if err != nil {
			return fmt.Errorf("Cannot get the branch protection: %w", err)
		}

		branch := args[0]
		branches, err := github.GetBranchProtectionRule(currentRepo.Owner, currentRepo.Name, branch)
		if err != nil {
			return fmt.Errorf("Cannot get the branch protection: %w", err)
		}
		for _, b := range branches {
			fmt.Println(b)
		}
		return nil
	},
}
