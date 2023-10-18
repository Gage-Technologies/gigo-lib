package models

type AutoGitSettings struct {
	RunOnStart     bool   `json:"runOnStart"`
	UpdateInterval int    `json:"updateInterval"`
	Logging        bool   `json:"logging"`
	Silent         bool   `json:"silent"`
	CommitMessage  string `json:"commitMessage"`
	Locale         string `json:"locale"`
	TimeZone       string `json:"timeZone"`
}

var DefaultAutoGitSettings = AutoGitSettings{
	RunOnStart:     true,
	UpdateInterval: 18,
	Logging:        true,
	Silent:         false,
	CommitMessage:  "--- Auto Git Commit ---",
	Locale:         "en-US",
	TimeZone:       "America/Chicago",
}

type WorkspaceSettings struct {
	AutoGit AutoGitSettings `json:"auto_git"`
}

var DefaultWorkspaceSettings = WorkspaceSettings{
	AutoGit: DefaultAutoGitSettings,
}
