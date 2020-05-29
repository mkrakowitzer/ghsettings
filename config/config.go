package config

type C struct {
	Repository struct {
		Name                string `yaml:"name"`
		Description         string `yaml:"description"`
		Homepage            string `yaml:"homepage"`
		Private             bool   `yaml:"private"`
		HasIssues           bool   `yaml:"has_issues"`
		HasProjects         bool   `yaml:"has_projects"`
		HasWiki             bool   `yaml:"has_wiki"`
		HasDownloads        bool   `yaml:"has_downloads"`
		DefaultBranch       string `yaml:"default_branch"`
		AllowSquashMerge    bool   `yaml:"allow_squash_merge"`
		AllowMergeCommit    bool   `yaml:"allow_merge_commit"`
		AllowRebaseMerge    bool   `yaml:"allow_rebase_merge"`
		DeleteBranchOnMerge bool   `yaml:"delete_branch_on_merge"`
	} `yaml:"repository"`
	Collaborators []struct {
		Username   string `yaml:"username"`
		Permission string `yaml:"permission"`
	} `yaml:"collaborators"`
	Teams []struct {
		Name       string `yaml:"name"`
		Permission string `yaml:"permission"`
	} `yaml:"teams"`
	Branches []struct {
		Name                         string   `yaml:"name"`
		RequiredApprovingReviewCount int      `yaml:"requiredApprovingReviewCount"`
		RequiresStatusChecks         bool     `yaml:"requiresStatusChecks"`
		RequiredStatusCheckContexts  []string `yaml:"requiredStatusCheckContexts"`
		RequiresApprovingReviews     bool     `yaml:"requiresApprovingReviews"`
		RequiresCodeOwnerReviews     bool     `yaml:"requiresCodeOwnerReviews"`
		RequiresCommitSignatures     bool     `yaml:"requiresCommitSignatures"`
		RequiresStrictStatusChecks   bool     `yaml:"requiresStrictStatusChecks"`
		RestrictsPushes              bool     `yaml:"restrictsPushes"`
		IsAdminEnforced              bool     `yaml:"isAdminEnforced"`
		DismissesStaleReviews        bool     `yaml:"dismissesStaleReviews"`
		PushActorIds                 []string `yaml:"pushActorIds"`
	} `yaml:"branches"`
}
