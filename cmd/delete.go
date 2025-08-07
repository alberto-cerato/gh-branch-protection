/*
Copyright © 2025 Alberto Cerato <macros123@gmail.com>
*/
package cmd

import (
	"fmt"

	"github.com/alberto-cerato/gh-branch-protection/internal/github"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(deleteCmd)
}

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete <branch>",
	Short: "Delete the protection from a branch",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return &WrongArgsError{Arg: 1, Cmd: cmd}
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		id := args[0]
		if err := github.DeleteBranchProtectionRule(id); err != nil {
			return fmt.Errorf("Cannot delete the branch protection: %w", err)
		}

		return nil
	},
}
