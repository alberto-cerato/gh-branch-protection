/*
Copyright © 2025 Alberto Cerato <macros123@gmail.com>
*/
package github

import (
	"fmt"
	"log/slog"

	"github.com/cli/go-gh/v2/pkg/api"
	graphql "github.com/cli/shurcooL-graphql"
)

// ListBranchProtectionRules returns a list of branch protection rules for the given repository owner and name.
// It queries the GitHub GraphQL API and paginates through all branch protection rules.
func ListBranchProtectionRules(repoOwner string, repoName string) ([]BranchProtectionRule, error) {
	rules := []BranchProtectionRule{}
	first := 100 // TODO: allow customization of page size

	client, err := api.DefaultGraphQLClient()
	if err != nil {
		return nil, fmt.Errorf("ListBranchProtectionRules: %w", err)
	}

	var query struct {
		Repository struct {
			BranchProtectionRules struct {
				Edges []struct {
					Node BranchProtectionRule
					Cursor string
				}
				PageInfo struct {
					EndCursor   graphql.String
					HasNextPage graphql.Boolean
				}
			} `graphql:"branchProtectionRules(first: $first, after: $cursor)"`
		} `graphql:"repository(owner: $repoOwner, name: $repoName)"`
	}

	variables := map[string]interface{}{
		"repoOwner": graphql.String(repoOwner),
		"repoName":  graphql.String(repoName),

		"cursor": graphql.String(""),
		"first":  graphql.Int(first),
	}
	for {
		err = client.Query("ListBranchProtectionRules", &query, variables)
		if err != nil {
			return nil, fmt.Errorf("ListBranchProtectionRules: %w", err)
		}

		for _, e := range query.Repository.BranchProtectionRules.Edges {
			r := e.Node

			if r.ID == nil {
				continue
			}
			slog.Debug("ListBranchProtectionRules", "Pattern", r.Pattern)
			rules = append(rules, r)
		}
		variables["cursor"] = query.Repository.BranchProtectionRules.PageInfo.EndCursor

		if !query.Repository.BranchProtectionRules.PageInfo.HasNextPage {
			break
		}
	}

	return rules, nil
}

// GetBranchProtectionRule retrieves the branch protection rule by its GraphQL node ID.
// Returns a pointer to BranchProtectionRule or an error if the rule is not found or the query fails.
func GetBranchProtectionRule(id string) (*BranchProtectionRule, error) {
	client, err := api.DefaultGraphQLClient()
	if err != nil {
		return nil, fmt.Errorf("branchProtectionRule: %w", err)
	}
	/*
	query GetRule {
	node(id: "BPR_kwDOPIadBM4D8Txc") {
		... on BranchProtectionRule {
		id
		pattern
		# Add any other fields you need here
		}
	}
	}*/
	var query struct {
		Node struct {
			Rule BranchProtectionRule `graphql:"... on BranchProtectionRule"`
		}  `graphql:"node(id: $id)"`
	}

	variables := map[string]interface{}{
		"id": graphql.ID(id),
	}
	err = client.Query("GetBranchProtection", &query, variables)
	if err != nil {
		return nil, fmt.Errorf("GetBranchProtectionRule: %w", err)
	}

	return &query.Node.Rule, nil
}

// DeleteBranchProtectionRule deletes the branch protection rule with the specified GraphQL node ID.
// Returns an error if the rule does not exist or the deletion fails.
func DeleteBranchProtectionRule(id string) error {
	client, err := api.DefaultGraphQLClient()
	if err != nil {
		return fmt.Errorf("DeleteBranchProtectionRule: %w", err)
	}

	var mutation struct {
		DeleteBranchProtectionRule struct {
			ClientMutationId graphql.ID
		} `graphql:" deleteBranchProtectionRule(input:{branchProtectionRuleId: $branchProtectionRuleId})"`
	}

	variables := map[string]interface{}{
		"branchProtectionRuleId": id,
	}

	if err = client.Mutate("DeleteBranchProtectionRule", &mutation, variables); err != nil {
		return fmt.Errorf("DeleteBranchProtectionRule: %w", err)
	}

	return nil
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
// CreateBranchProtectionRule creates a new branch protection rule for the specified repository and pattern.
// The rule parameter should be populated with the desired protection settings.
// Returns an error if the creation fails.
func CreateBranchProtectionRule(repositoryId graphql.ID, pattern graphql.String, rule CreateBranchProtectionRuleInput) error {
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
			}
		} `graphql:"createBranchProtectionRule(input: $input)"`
	}

	rule.RepositoryID = repositoryId
	rule.Pattern = pattern
	variables := map[string]interface{}{
		"input": rule,
	}

	if err = client.Mutate("CreateBranchProtectionRule", &mutation, variables); err != nil {
		// TODO: if already exist maybe we should update the existing rule
		return fmt.Errorf("CreateBranchProtectionRule: %w", err)
	}
	return nil
}
