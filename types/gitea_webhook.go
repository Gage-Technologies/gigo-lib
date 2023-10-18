package types

import "time"

// GiteaWebhookPush
//
//	Struct modeling the return of a Gitea webhook fired
//	on a repository push
type GiteaWebhookPush struct {
	After   string `json:"after"`
	Before  string `json:"before"`
	Commits []struct {
		Added  []string `json:"added"`
		Author struct {
			Email    string `json:"email"`
			Name     string `json:"name"`
			Username string `json:"username"`
		} `json:"author"`
		Committer struct {
			Email    string `json:"email"`
			Name     string `json:"name"`
			Username string `json:"username"`
		} `json:"committer"`
		ID           string      `json:"id"`
		Message      string      `json:"message"`
		Modified     []string    `json:"modified"`
		Removed      []string    `json:"removed"`
		Timestamp    time.Time   `json:"timestamp"`
		URL          string      `json:"url"`
		Verification interface{} `json:"verification"`
	} `json:"commits"`
	CompareURL string `json:"compare_url"`
	HeadCommit struct {
		Added  []string `json:"added"`
		Author struct {
			Email    string `json:"email"`
			Name     string `json:"name"`
			Username string `json:"username"`
		} `json:"author"`
		Committer struct {
			Email    string `json:"email"`
			Name     string `json:"name"`
			Username string `json:"username"`
		} `json:"committer"`
		ID           string      `json:"id"`
		Message      string      `json:"message"`
		Modified     []string    `json:"modified"`
		Removed      []string    `json:"removed"`
		Timestamp    time.Time   `json:"timestamp"`
		URL          string      `json:"url"`
		Verification interface{} `json:"verification"`
	} `json:"head_commit"`
	Pusher struct {
		Active            bool      `json:"active"`
		AvatarURL         string    `json:"avatar_url"`
		Created           time.Time `json:"created"`
		Description       string    `json:"description"`
		Email             string    `json:"email"`
		FollowersCount    int64     `json:"followers_count"`
		FollowingCount    int64     `json:"following_count"`
		FullName          string    `json:"full_name"`
		ID                int64     `json:"id"`
		IsAdmin           bool      `json:"is_admin"`
		Language          string    `json:"language"`
		LastLogin         time.Time `json:"last_login"`
		Location          string    `json:"location"`
		Login             string    `json:"login"`
		ProhibitLogin     bool      `json:"prohibit_login"`
		Restricted        bool      `json:"restricted"`
		StarredReposCount int64     `json:"starred_repos_count"`
		Username          string    `json:"username"`
		Visibility        string    `json:"visibility"`
		Website           string    `json:"website"`
	} `json:"pusher"`
	Ref        string `json:"ref"`
	Repository struct {
		AllowMergeCommits         bool      `json:"allow_merge_commits"`
		AllowRebase               bool      `json:"allow_rebase"`
		AllowRebaseExplicit       bool      `json:"allow_rebase_explicit"`
		AllowSquashMerge          bool      `json:"allow_squash_merge"`
		Archived                  bool      `json:"archived"`
		AvatarURL                 string    `json:"avatar_url"`
		CloneURL                  string    `json:"clone_url"`
		CreatedAt                 time.Time `json:"created_at"`
		DefaultBranch             string    `json:"default_branch"`
		DefaultMergeStyle         string    `json:"default_merge_style"`
		Description               string    `json:"description"`
		Empty                     bool      `json:"empty"`
		Fork                      bool      `json:"fork"`
		ForksCount                int64     `json:"forks_count"`
		FullName                  string    `json:"full_name"`
		HasIssues                 bool      `json:"has_issues"`
		HasProjects               bool      `json:"has_projects"`
		HasPullRequests           bool      `json:"has_pull_requests"`
		HasWiki                   bool      `json:"has_wiki"`
		HtmlURL                   string    `json:"html_url"`
		ID                        int64     `json:"id"`
		IgnoreWhitespaceConflicts bool      `json:"ignore_whitespace_conflicts"`
		int64ernal                bool      `json:"int64ernal"`
		int64ernalTracker         struct {
			AllowOnlyContributorsToTrackTime bool `json:"allow_only_contributors_to_track_time"`
			EnableIssueDependencies          bool `json:"enable_issue_dependencies"`
			EnableTimeTracker                bool `json:"enable_time_tracker"`
		} `json:"int64ernal_tracker"`
		Language         string    `json:"language"`
		LanguagesURL     string    `json:"languages_url"`
		Mirror           bool      `json:"mirror"`
		Mirrorint64erval string    `json:"mirror_int64erval"`
		MirrorUpdated    time.Time `json:"mirror_updated"`
		Name             string    `json:"name"`
		OpenIssuesCount  int64     `json:"open_issues_count"`
		OpenPrCounter    int64     `json:"open_pr_counter"`
		OriginalURL      string    `json:"original_url"`
		Owner            struct {
			Active            bool      `json:"active"`
			AvatarURL         string    `json:"avatar_url"`
			Created           time.Time `json:"created"`
			Description       string    `json:"description"`
			Email             string    `json:"email"`
			FollowersCount    int64     `json:"followers_count"`
			FollowingCount    int64     `json:"following_count"`
			FullName          string    `json:"full_name"`
			ID                int64     `json:"id"`
			IsAdmin           bool      `json:"is_admin"`
			Language          string    `json:"language"`
			LastLogin         time.Time `json:"last_login"`
			Location          string    `json:"location"`
			Login             string    `json:"login"`
			ProhibitLogin     bool      `json:"prohibit_login"`
			Restricted        bool      `json:"restricted"`
			StarredReposCount int64     `json:"starred_repos_count"`
			Username          string    `json:"username"`
			Visibility        string    `json:"visibility"`
			Website           string    `json:"website"`
		} `json:"owner"`
		Parent      interface{} `json:"parent"`
		Permissions struct {
			Admin bool `json:"admin"`
			Pull  bool `json:"pull"`
			Push  bool `json:"push"`
		} `json:"permissions"`
		Private        bool        `json:"private"`
		ReleaseCounter int64       `json:"release_counter"`
		RepoTransfer   interface{} `json:"repo_transfer"`
		Size           int64       `json:"size"`
		SSHURL         string      `json:"ssh_url"`
		StarsCount     int64       `json:"stars_count"`
		Template       bool        `json:"template"`
		UpdatedAt      time.Time   `json:"updated_at"`
		WatchersCount  int64       `json:"watchers_count"`
		Website        string      `json:"website"`
	} `json:"repository"`
	Sender struct {
		Active            bool      `json:"active"`
		AvatarURL         string    `json:"avatar_url"`
		Created           time.Time `json:"created"`
		Description       string    `json:"description"`
		Email             string    `json:"email"`
		FollowersCount    int64     `json:"followers_count"`
		FollowingCount    int64     `json:"following_count"`
		FullName          string    `json:"full_name"`
		ID                int64     `json:"id"`
		IsAdmin           bool      `json:"is_admin"`
		Language          string    `json:"language"`
		LastLogin         time.Time `json:"last_login"`
		Location          string    `json:"location"`
		Login             string    `json:"login"`
		ProhibitLogin     bool      `json:"prohibit_login"`
		Restricted        bool      `json:"restricted"`
		StarredReposCount int64     `json:"starred_repos_count"`
		Username          string    `json:"username"`
		Visibility        string    `json:"visibility"`
		Website           string    `json:"website"`
	} `json:"sender"`
	TotalCommits int64 `json:"total_commits"`
}
