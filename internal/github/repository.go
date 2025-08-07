/*
Copyright © 2025 Alberto Cerato <macros123@gmail.com>
*/
package github

import (
	"fmt"

	"github.com/cli/go-gh/v2/pkg/api"
	graphql "github.com/cli/shurcooL-graphql"
)

// GetRepoID returns the GraphQL repository ID for the given owner and repository name.
// It queries the GitHub GraphQL API and returns the repository ID or an error if the query fails.
func GetRepoID(owner string, name string) (graphql.ID, error) {
	client, err := api.DefaultGraphQLClient()
	if err != nil {
		return nil, fmt.Errorf("GetRepoID: %w", err)
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
	err = client.Query("GetRepoID", &query, variables)
	if err != nil {
		return nil, fmt.Errorf("GetRepoID: %w", err)
	}

	return query.Repository.ID, nil
}
