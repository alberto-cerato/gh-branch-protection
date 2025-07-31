/*
Copyright © 2025 Alberto Cerato <macros123@gmail.com>
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"

	"github.com/cli/go-gh/v2/pkg/api"
	"github.com/cli/go-gh/v2/pkg/repository"
	graphql "github.com/cli/shurcooL-graphql"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(setCmd)
}

type BranchProtectionRule struct {
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

// setCmd represents the set command
var setCmd = &cobra.Command{
	Use:   "set <branch>",
	Short: "Set the protection for a branch",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		branch := args[0]

		currentRepo, err := repository.Current()
		if err != nil {
			return fmt.Errorf("Cannot set branch protection configuration: %w", err)
		}

		var rule BranchProtectionRule
		err = json.NewDecoder(os.Stdin).Decode(&rule)
		if err != nil {
			return fmt.Errorf("Cannot set branch protection configuration: %w", err)
		}
		slog.Debug("setCmd", "operation", "decode JSON input", "rule", rule)

		repoId, err := FindRepoID(currentRepo.Owner, currentRepo.Name)
		if err != nil {
			return fmt.Errorf("Cannot set branch protection configuration: %w", err)

		}
		if err := CreateBranchProtectionRule(repoId, branch, rule); err != nil {
			return fmt.Errorf("Cannot set branch protection configuration: %w", err)

		}
		return nil
	},
}

func FindRepoID(owner string, name string) (graphql.ID, error) {
	client, err := api.DefaultGraphQLClient()
	if err != nil {
		return nil, fmt.Errorf("CreateBranchProtectionRule: %w", err)
	}

	var query struct {
		Repository struct {
			ID graphql.ID
		} `graphql:"repository(owner: $owner, name: $name)"`
	}
	variables := map[string]interface{}{
		"owner": graphql.String(owner),
		"name":  graphql.String(name),
	}
	err = client.Query("GetOwnerID", &query, variables)
	if err != nil {
		return nil, fmt.Errorf("CreateRepo: %w", err)
	}

	return query.Repository.ID, nil
}

/*
This is the GraphQL that inspired the implementation of the below function.

mutation AddProtection {
  createBranchProtectionRule(input:{repositoryId: "R_kgDOPIadBA", pattern:"master", requiresLinearHistory: true}) {
    branchProtectionRule {
      repository {
        id
      }
      requiresLinearHistory
      allowsForcePushes
      reviewDismissalAllowances {
        nodes {
          actor {
            __typename
          }
        }
      }
    }
  }
}
*/
func CreateBranchProtectionRule(repositoryId graphql.ID, branch string, rule BranchProtectionRule) error {
	client, err := api.DefaultGraphQLClient()
	if err != nil {
		return fmt.Errorf("CreateBranchProtectionRule: %w", err)
	}

	var mutation struct {
		CreateBranchProtectionRule struct {
			BranchProtectionRule struct {
				Repository struct {
					ID graphql.ID
				}
				RequiresLinearHistory graphql.Boolean
			}
		} `graphql:"createBranchProtectionRule(input: {repositoryId: $repositoryId, pattern: $pattern, requiresLinearHistory: $requiresLinearHistory})"`
	}
	variables := map[string]interface{}{
		"repositoryId":          repositoryId,
		"pattern":               graphql.String(branch),
		"requiresLinearHistory": graphql.Boolean(rule.RequiresLinearHistory),
		// TODO: add missing branch protections
	}

	if err = client.Mutate("CreateBranchProtectionRule", &mutation, variables); err != nil {
		return fmt.Errorf("CreateBranchProtectionRule: %w", err)
	}
	return nil
}

