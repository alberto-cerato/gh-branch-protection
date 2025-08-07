/*
Copyright © 2025 Alberto Cerato <macros123@gmail.com>
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/cli/go-gh/v2/pkg/repository"
	graphql "github.com/cli/shurcooL-graphql"
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
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return &WrongArgsError{Arg: 1, Cmd: cmd}
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		pattern := args[0]

		currentRepo, err := repository.Current()
		if err != nil {
			return fmt.Errorf("Cannot set the branch protection: %w", err)
		}

		var rule github.CreateBranchProtectionRuleInput
		err = json.NewDecoder(os.Stdin).Decode(&rule)
		if err != nil {
			return fmt.Errorf("Cannot set the branch protection: %w", err)
		}

		repoId, err := github.GetRepoID(currentRepo.Owner, currentRepo.Name)
		if err != nil {
			return fmt.Errorf("Cannot set the branch protection: %w", err)
		}
		if err := github.CreateBranchProtectionRule(repoId, graphql.String(pattern), rule); err != nil {
			return fmt.Errorf("Cannot set the branch protection: %w", err)

		}
		return nil
	},
}
