/*
Copyright © 2025 Alberto Cerato <macros123@gmail.com>
*/
package cmd

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/cli/go-gh/v2/pkg/api"
	"github.com/cli/go-gh/v2/pkg/auth"
	"github.com/cli/go-gh/v2/pkg/repository"
	graphql "github.com/cli/shurcooL-graphql"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		currentRepo, err := repository.Current()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Cannot list protected branches: %s", err)
			os.Exit(1)
		}
		fmt.Println(currentRepo.Host)
		token, _ := auth.TokenForHost(currentRepo.Host)

		branches, err := ListProtectedBranches(currentRepo.Host, token, currentRepo.Owner, currentRepo.Name)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Cannot list protected branches: %s", err)
			os.Exit(1)
		}
		for _, b := range branches {
			fmt.Println(b)
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func ListProtectedBranches(host string, token string, repoOwner string, repoName string) ([]string, error) {
	branches := []string{}
	first := 100 // TODO: allow customization of page size

	client, err := api.DefaultGraphQLClient()
	if err != nil {
		return nil, fmt.Errorf("ListProtectedBranches: %w", err)
	}

	var query struct {
		Repository struct {
			Refs struct {
				Edges []struct {
					Node struct {
						Name                 graphql.String
						BranchProtectionRule struct {
							ID graphql.ID
						}
					}
					Cursor string
				}
				PageInfo struct {
					EndCursor   graphql.String
					HasNextPage graphql.Boolean
				}
			} `graphql:"refs(refPrefix: \"refs/heads/\", first: $first, after: $cursor)"`
		} `graphql:"repository(owner: $repoOwner, name: $repoName)"`
	}

	variables := map[string]interface{}{
		"repoOwner": graphql.String(repoOwner),
		"repoName":  graphql.String(repoName),

		"cursor": graphql.String(""),
		"first":  graphql.Int(first),
	}
	for {
		err = client.Query("ListProtectedBranches", &query, variables)
		if err != nil {
			return nil, fmt.Errorf("ListProtectedBranches: %w", err)
		}

		for _, e := range query.Repository.Refs.Edges {
			b := e.Node

			if b.BranchProtectionRule.ID == nil {
				continue
			}
			slog.Debug("ListProtectedBranches", "Branch", b.Name)
			branches = append(branches, string(b.Name))
		}
		variables["cursor"] = query.Repository.Refs.PageInfo.EndCursor

		if !query.Repository.Refs.PageInfo.HasNextPage {
			break
		}
	}

	return branches, nil
}
