/*
Copyright © 2025 Alberto Cerato <macros123@gmail.com>
*/
package github

import (
	graphql "github.com/cli/shurcooL-graphql"
)

// TODO: we could generate the 2 flat structures using go:generate as they are almost the same

// https://docs.github.com/en/graphql/reference/input-objects#createbranchprotectionruleinput
type CreateBranchProtectionRuleInput struct {
	RepositoryID                   graphql.ID      `json:"repositoryId"`
	Pattern                        graphql.String  `json:"pattern"`
	AllowsDeletions                graphql.Boolean `json:"allowsDeletions"`
	AllowsForcePushes              graphql.Boolean `json:"allowsForcePushes"`
	RequiredApprovingReviewCount   graphql.Int     `json:"requiredApprovingReviewCount"`
	RequiresApprovingReviews       graphql.Boolean `json:"requiresApprovingReviews"`
	RequiresCodeOwnerReviews       graphql.Boolean `json:"requiresCodeOwnerReviews"`
	RequiresCommitSignatures       graphql.Boolean `json:"requiresCommitSignatures"`
	RequiresLinearHistory          graphql.Boolean `json:"requiresLinearHistory"`
	RequiresConversationResolution graphql.Boolean `json:"requiresConversationResolution"`
	IsAdminEnforced                graphql.Boolean `json:"isAdminEnforced"`
	RestrictsPushes                graphql.Boolean `json:"restrictsPushes"`
	RestrictsReviewDismissals      graphql.Boolean `json:"restrictsReviewDismissals"`
}

// https://docs.github.com/en/graphql/reference/objects#branchprotectionrule
type BranchProtectionRule struct {
	ID                             graphql.ID
	Pattern                        graphql.String
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
