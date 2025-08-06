/*
Copyright © 2025 Alberto Cerato <macros123@gmail.com>
*/
package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gh-branch-protection",
	Short: "List, get, set, and delete protections on branches",
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		argsError := &WrongArgsError{}
		if errors.As(err, &argsError) {
			argsError.Cmd.Usage()
		} else {
			fmt.Fprintf(os.Stderr, "%s\n", err)
		}

		os.Exit(1)
	}
}

func init() {
	// Suppress usage and error output for subcommands when an error occurs.
	rootCmd.SilenceUsage = true
	rootCmd.SilenceErrors = true
}
