/*
Copyright © 2025 Alberto Cerato <macros123@gmail.com>
*/
package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/alberto-cerato/gh-branch-protection/internal/github"
)

func init() {
	rootCmd.AddCommand(getCmd)
}

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get <rule_id>",
	Short: "Get a branch protection rule definition",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return &WrongArgsError{Arg: 1, Cmd: cmd}
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		id := args[0]
		rule, err := github.GetBranchProtectionRule(id)
		if err != nil {
			return fmt.Errorf("Cannot get the branch protection: %w", err)
		}
		b, err := json.MarshalIndent(rule, "", "  ")
		if err != nil {
			return fmt.Errorf("Cannot get the branch protection: %w", err)
		}
		fmt.Println(string(b))

		return nil
	},
}
