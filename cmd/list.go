/*
Copyright © 2025 Alberto Cerato <macros123@gmail.com>
*/
package cmd

import (
	"fmt"
	"log/slog"

	"github.com/cli/go-gh/v2/pkg/api"
	"github.com/cli/go-gh/v2/pkg/repository"
	graphql "github.com/cli/shurcooL-graphql"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(listCmd)
}

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List protected branches",
	RunE: func(cmd *cobra.Command, args []string) error {
		currentRepo, err := repository.Current()
		if err != nil {
			return fmt.Errorf("Cannot list protected branches: %w\n", err)
			
		}

		branches, err := ListProtectedBranches(currentRepo.Owner, currentRepo.Name)
		if err != nil {
			return fmt.Errorf("Cannot list protected branches: %w\n", err)
		}
		for _, b := range branches {
			fmt.Println(b)
		}

		return nil
	},
}

func ListProtectedBranches(repoOwner string, repoName string) ([]string, error) {
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
