package api

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/mkrakowitzer/ghsettings/config"
	"github.com/shurcooL/githubv4"
)

type RepoPayload struct {
	Organization struct {
		Repository struct {
			ID githubv4.ID
		}
	}
}

func GetRepoID(client *Client, org string, reponame string) (*RepoPayload, error) {
	query := `
	query($org: String!, $name: String!) {
	  organization(login: $org) {
	    repository(name: $name) {
	      id
	    }
	  }
	}`
	variables := map[string]interface{}{"org": org, "name": reponame}
	result := RepoPayload{}

	err := client.GraphQL(query, variables, &result)
	return &result, err
}

func UpdateRepository(apiClient *Client, repo *RepoPayload, config config.C) error {
	variables := map[string]interface{}{
		"id":          repo.Organization.Repository.ID,
		"homepage":    config.Repository.Homepage,
		"issues":      config.Repository.HasIssues,
		"wiki":        config.Repository.HasWiki,
		"projects":    config.Repository.HasProjects,
		"description": config.Repository.Description,
	}

	err := UpdateRepositoryV4(apiClient, variables)
	if err != nil {
		return err
	}

	err = UpdateRepositoryV3(apiClient, repo, config)
	return err
}

func UpdateRepositoryV4(client *Client, variables map[string]interface{}) error {
	mutation := `
	mutation ($id: ID!, $homepage: String!, $wiki: Boolean!, $projects: Boolean!, $issues: Boolean!, $description: String!) {	
		updateRepository(input: {
			repositoryId: $id,
			homepageUrl: $homepage,
			hasWikiEnabled: $wiki,
			hasProjectsEnabled: $projects,
			hasIssuesEnabled: $issues,
			description: $description,
		}) {
		  clientMutationId
		}
	  }`
	err := client.GraphQL(mutation, variables, nil)
	return err
}

type Repository struct {
	Private             bool   `json:"private"`
	DefaultBranch       string `json:"default_branch"`
	AllowRebaseMerge    bool   `json:"allow_rebase_merge"`
	AllowSquashMerge    bool   `json:"allow_squash_merge"`
	AllowMergeCommit    bool   `json:"allow_merge_commit"`
	DeleteBranchOnMerge bool   `json:"delete_branch_on_merge"`
}

func UpdateRepositoryV3(client *Client, repo *RepoPayload, config config.C) error {

	path := fmt.Sprintf("repos/%s/%s", Org, config.Repository.Name)
	result := Repository{}

	updateRepo := Repository{
		Private:             config.Repository.Private,
		DefaultBranch:       "master",
		AllowRebaseMerge:    config.Repository.AllowRebaseMerge,
		AllowSquashMerge:    config.Repository.AllowSquashMerge,
		AllowMergeCommit:    config.Repository.AllowMergeCommit,
		DeleteBranchOnMerge: config.Repository.DeleteBranchOnMerge,
	}

	j, _ := json.Marshal(updateRepo)

	err := client.REST("PATCH", path, bytes.NewBuffer(j), &result)
	return err

}
