package api

import (
	"github.com/mkrakowitzer/ghsettings/config"
	"github.com/mkrakowitzer/ghsettings/utils"
	"github.com/shurcooL/githubv4"
)

type BranchProtectionRules struct {
	Organization struct {
		Repository struct {
			BranchProtectionRules struct {
				Nodes []struct {
					ID                           githubv4.ID `json:"id"`
					Pattern                      string      `json:"pattern"`
					RequiredApprovingReviewCount int         `json:"requiredApprovingReviewCount"`
					RequiredStatusCheckContexts  []string    `json:"requiredStatusCheckContexts"`
					RequiresApprovingReviews     bool        `json:"requiresApprovingReviews"`
					RequiresCodeOwnerReviews     bool        `json:"requiresCodeOwnerReviews"`
					RequiresCommitSignatures     bool        `json:"requiresCommitSignatures"`
					RequiresStatusChecks         bool        `json:"requiresStatusChecks"`
					RequiresStrictStatusChecks   bool        `json:"requiresStrictStatusChecks"`
					RestrictsPushes              bool        `json:"restrictsPushes"`
					IsAdminEnforced              bool        `json:"isAdminEnforced"`
					DismissesStaleReviews        bool        `json:"dismissesStaleReviews"`
					PushActorIds                 []string    `json:"pushActorIds"`
				} `json:"nodes"`
			} `json:"branchProtectionRules"`
		} `json:"repository"`
	} `json:"organization"`
}

func BranchProtections(apiClient *Client, repo *RepoPayload, config config.C) error {

	rules, err := GetBranchProtectionRules(apiClient, config.Repository.Name)
	if err != nil {
		return err
	}

	var variables map[string]interface{}

	for _, s := range config.Branches {

		found := false
		variables = map[string]interface{}{
			"name":                         s.Name,
			"requiredApprovingReviewCount": s.RequiredApprovingReviewCount,
			"requiresStatusChecks":         s.RequiresStatusChecks,
			"requiredStatusCheckContexts":  s.RequiredStatusCheckContexts,
			"requiresApprovingReviews":     s.RequiresApprovingReviews,
			"requiresCodeOwnerReviews":     s.RequiresCodeOwnerReviews,
			"requiresCommitSignatures":     s.RequiresCommitSignatures,
			"requiresStrictStatusChecks":   s.RequiresStrictStatusChecks,
			"restrictsPushes":              s.RestrictsPushes,
			"isAdminEnforced":              s.IsAdminEnforced,
			"dismissesStaleReviews":        s.DismissesStaleReviews,
			"pushActorIds":                 s.PushActorIds,
			"repositoryId":                 repo.Organization.Repository.ID,
		}
		for _, k := range rules.Organization.Repository.BranchProtectionRules.Nodes {
			if k.Pattern == s.Name {
				variables["branchProtectionRuleId"] = k.ID
				err := UpdateBranchProtections(apiClient, variables)
				if err != nil {
					return err
				}
				found = true
			}
		}
		if found {
			continue
		}
		err := CreateBranchProtections(apiClient, variables)
		if err != nil {
			return err
		}
	}
	err = DeleteBranchProtections(apiClient, config, rules)
	if err != nil {
		return err
	}
	return nil
}

func GetBranchProtectionRules(client *Client, reponame string) (*BranchProtectionRules, error) {
	query := `query($org: String!, $name: String!) {
		organization(login: $org) {
			repository(name: $name) {
					branchProtectionRules(first: 100) {
					nodes {
						id
						pattern
						requiredApprovingReviewCount
						requiredStatusCheckContexts
						requiresApprovingReviews
						requiresCodeOwnerReviews
						requiresCommitSignatures
						requiresStatusChecks
						requiresStrictStatusChecks
						restrictsPushes
						isAdminEnforced
						dismissesStaleReviews
					}
			    }
			}
		}
	}`

	variables := map[string]interface{}{"org": Org, "name": reponame}
	result := BranchProtectionRules{}

	err := client.GraphQL(query, variables, &result)
	return &result, err
}

func CreateBranchProtections(client *Client, variables map[string]interface{}) error {
	mutation := `
	  mutation (
		$repositoryId: ID!,
		$name: String!
		$isAdminEnforced: Boolean!,
		$dismissesStaleReviews: Boolean!,
		$requiredApprovingReviewCount: Int!,
		$requiresApprovingReviews: Boolean!,
		$requiresCodeOwnerReviews: Boolean!,
		$requiresCommitSignatures: Boolean!,
		$requiresStatusChecks: Boolean!,
		$requiresStrictStatusChecks: Boolean!,
		$restrictsPushes: Boolean!,
		$requiredStatusCheckContexts: [String!]
		$pushActorIds: [String!]
		) {
		createBranchProtectionRule(input: {
			repositoryId: $repositoryId,
			pattern: $name,
			restrictsPushes: $restrictsPushes,
			isAdminEnforced: $isAdminEnforced,
			dismissesStaleReviews: $dismissesStaleReviews,
			requiredApprovingReviewCount: $requiredApprovingReviewCount,
			requiresApprovingReviews: $requiresApprovingReviews,
			requiresCodeOwnerReviews: $requiresCodeOwnerReviews,
			requiresCommitSignatures: $requiresCommitSignatures,
			requiresStatusChecks: $requiresStatusChecks,
			requiresStrictStatusChecks: $requiresStrictStatusChecks,
			requiredStatusCheckContexts: $requiredStatusCheckContexts,
			pushActorIds: $pushActorIds,
		}) {
		clientMutationId
		}
	}`
	err := client.GraphQL(mutation, variables, nil)
	return err
}

func UpdateBranchProtections(client *Client, variables map[string]interface{}) error {

	mutation := `
	  mutation (
		$branchProtectionRuleId: ID!,
		$name: String!
		$isAdminEnforced: Boolean!,
		$dismissesStaleReviews: Boolean!,
		$requiredApprovingReviewCount: Int!,
		$requiresApprovingReviews: Boolean!,
		$requiresCodeOwnerReviews: Boolean!,
		$requiresCommitSignatures: Boolean!,
		$requiresStatusChecks: Boolean!,
		$requiresStrictStatusChecks: Boolean!,
		$restrictsPushes: Boolean!,
		$requiredStatusCheckContexts: [String!]
		$pushActorIds: [String!]
		) {
		updateBranchProtectionRule(input: {
			branchProtectionRuleId: $branchProtectionRuleId,
			pattern: $name,
			restrictsPushes: $restrictsPushes,
			isAdminEnforced: $isAdminEnforced,
			dismissesStaleReviews: $dismissesStaleReviews,
			requiredApprovingReviewCount: $requiredApprovingReviewCount,
			requiresApprovingReviews: $requiresApprovingReviews,
			requiresCodeOwnerReviews: $requiresCodeOwnerReviews,
			requiresCommitSignatures: $requiresCommitSignatures,
			requiresStatusChecks: $requiresStatusChecks,
			requiresStrictStatusChecks: $requiresStrictStatusChecks,
			requiredStatusCheckContexts: $requiredStatusCheckContexts,
			pushActorIds: $pushActorIds,
		}) {
		clientMutationId
		}
	  }`
	err := client.GraphQL(mutation, variables, nil)
	return err
}

// Refactor this, was in a hurry
func DeleteBranchProtections(client *Client, config config.C, rules *BranchProtectionRules) error {

	var yml_rules []string
	var gh_rules []string
	for _, s := range config.Branches {
		yml_rules = append(yml_rules, s.Name)
	}
	for _, k := range rules.Organization.Repository.BranchProtectionRules.Nodes {
		gh_rules = append(gh_rules, k.Pattern)
	}
	delete := utils.Missing(yml_rules, gh_rules)

	for _, s := range delete {
		for _, k := range rules.Organization.Repository.BranchProtectionRules.Nodes {
			if k.Pattern == s {
				variables := map[string]interface{}{
					"branchProtectionRuleId": k.ID,
				}
				mutation := `
				mutation (
				  $branchProtectionRuleId: ID!,
				  ) {
					  deleteBranchProtectionRule(input: {branchProtectionRuleId: $branchProtectionRuleId}) {
						  clientMutationId
					  }
				  }`
				err := client.GraphQL(mutation, variables, nil)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}
