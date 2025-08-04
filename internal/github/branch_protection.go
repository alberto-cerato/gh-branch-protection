/*
Copyright © 2025 Alberto Cerato <macros123@gmail.com>
*/
package github

import (
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/cli/go-gh/v2/pkg/api"
	graphql "github.com/cli/shurcooL-graphql"
)

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

func branchProtectionRule(repoOwner string, repoName string, branch string) (*BranchProtectionRule, error) {
	client, err := api.DefaultGraphQLClient()
	if err != nil {
		return nil, fmt.Errorf("branchProtectionRule: %w", err)
	}

	var query struct {
		Repository struct {
			Name graphql.String
			Ref  struct {
				Name                 graphql.String
				BranchProtectionRule BranchProtectionRule
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
		return nil, fmt.Errorf("branchProtectionRule: %w", err)
	}

	return &query.Repository.Ref.BranchProtectionRule, nil
}

func GetBranchProtectionRule(repoOwner string, repoName string, branch string) (string, error) {
	rule, err := branchProtectionRule(repoOwner, repoName, branch)
	if err != nil {
		return "", nil
	}

	b, err := json.MarshalIndent(rule, "", "  ")
	if err != nil {
		return "", err
	}

	return string(b), nil
}

func DeleteBranchProtectionRule(repoOwner string, repoName string, branch string) error {
	rule, err := branchProtectionRule(repoOwner, repoName, branch)
	if err != nil {
		return fmt.Errorf("DeleteBranchProtectionRule: %w", err)
	}
	if rule.ID == nil {
		return fmt.Errorf("DeleteBranchProtectionRule: no rule to delete")
	}

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
		"branchProtectionRuleId": rule.ID,
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
func CreateBranchProtectionRule(repositoryId graphql.ID, branch graphql.String, rule CreateBranchProtectionRuleInput) error {
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
	rule.Pattern = branch
	variables := map[string]interface{}{
		"input": rule,
	}

	if err = client.Mutate("CreateBranchProtectionRule", &mutation, variables); err != nil {
		// TODO: if already exist maybe we should update the existing rule
		return fmt.Errorf("CreateBranchProtectionRule: %w", err)
	}
	return nil
}
