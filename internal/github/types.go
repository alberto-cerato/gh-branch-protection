/*
Copyright © 2025 Alberto Cerato <macros123@gmail.com>
*/
package github

import graphql "github.com/cli/shurcooL-graphql"

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
