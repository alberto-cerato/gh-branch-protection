/*
Copyright © 2025 Alberto Cerato <macros123@gmail.com>
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gh-branch-protection",
	Short: "Allow to get and set branch protections",
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.gh-branch-protection.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.

	// TODO
	//rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}


