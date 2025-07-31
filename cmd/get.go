/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/cli/go-gh/v2/pkg/api"
	"github.com/cli/go-gh/v2/pkg/auth"
	"github.com/cli/go-gh/v2/pkg/repository"
	graphql "github.com/cli/shurcooL-graphql"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(getCmd)
}

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Get the branch protection configuration",
	Run: func(cmd *cobra.Command, args []string) {
		currentRepo, err := repository.Current()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Cannot list protected branches: %s", err)
			os.Exit(1)
		}
		token, _ := auth.TokenForHost(currentRepo.Host)

		branch := args[0]
		branches, err := GetBranchProtection(currentRepo.Host, token, currentRepo.Owner, currentRepo.Name, branch)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Cannot list protected branches: %s\n", err)
			os.Exit(1)
		}
		for _, b := range branches {
			fmt.Println(b)
		}
	},
}

func GetBranchProtection(host string, token string, repoOwner string, repoName string, branch string) ([]string, error) {
	branches := []string{}

	client, err := api.DefaultGraphQLClient()
	if err != nil {
		return nil, fmt.Errorf("GetBranchProtection: %w", err)
	}

	var query struct {
		Repository struct {
			Ref struct {
				Name                 graphql.String
				BranchProtectionRule struct {
					ID                             graphql.ID
					AllowsDeletions                graphql.Boolean
					AllowsForcePushes              graphql.Boolean
					RequiredApprovingReviewCount   graphql.Int
					RequiresApprovingReviews       graphql.Boolean
					RequiresCodeOwnerReviews       graphql.Boolean
					RequiresCommitSignatures       graphql.Boolean
					RequiresLinearHistory          graphql.Boolean
					RequiresConversationResolution graphql.Boolean
					IsAdminEnforced                graphql.Boolean
					RestrictsPushes                graphql.Boolean
					RestrictsReviewDismissals      graphql.Boolean
				}
			} `graphql:"ref(qualifiedName: $branch)"`
		} `graphql:"repository(owner: $repoOwner, name: $repoName)"`
	}

	variables := map[string]interface{}{
		"repoOwner": graphql.String(repoOwner),
		"repoName":  graphql.String(repoName),
		"branch":    graphql.String("refs/heads/" + branch),
	}

	err = client.Query("GetBranchProtection", &query, variables)
	if err != nil {
		return nil, fmt.Errorf("GetBranchProtection: %w", err)
	}

	b, err := json.MarshalIndent(query.Repository, "", "  ")
	fmt.Println(string(b[:]))
	if err != nil {
		return nil, err
	}

	return branches, nil
}
