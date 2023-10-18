package git

import (
	"crypto/tls"
	"errors"
	"fmt"
	"github.com/gage-technologies/gitea-go/gitea"
	"net/http"
	"net/url"
)

type VCSClient struct {
	GiteaClient *gitea.Client
	HostUrl     string
	AdminID     int64
}

func CreateVCSClient(hosturl string, username string, password string, insecure bool) (*VCSClient, error) {
	// create http client
	httpClient := http.DefaultClient
	if insecure {
		httpClient.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
	}

	client, err := gitea.NewClient(hosturl, gitea.SetBasicAuth(username, password), gitea.SetHTTPClient(httpClient))
	if err != nil {
		return nil, errors.New(fmt.Sprintf("failed to create new gitea client, err: %v", err))
	}

	if client == nil {
		return nil, errors.New(fmt.Sprintf("failed to create gitea client, err: client returned nil"))
	}

	// retrieve user to get admin id
	user, res, err := client.GetUserInfo(username)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve admin user from gitea: %v\n    res: %+v", err, res)
	}

	return &VCSClient{
		GiteaClient: client,
		HostUrl:     hosturl,
		AdminID:     user.ID,
	}, nil
}

func (v *VCSClient) LoginAsUser(username string, password string) (*gitea.Client, error) {
	return gitea.NewClient(v.HostUrl, gitea.SetBasicAuth(username, password))
}

func (v *VCSClient) CreateRepo(userName string, repoName string, repoDescription string, isPrivate bool,
	readMePath string, gitIgnorePath string, licensePath string, defaultBranch string) (*gitea.Repository, error) {
	repo, res, err := v.GiteaClient.AdminCreateRepo(userName, gitea.CreateRepoOption{
		Name:          repoName,
		Description:   repoDescription,
		Private:       isPrivate,
		IssueLabels:   "",
		AutoInit:      true,
		Template:      false,
		Gitignores:    gitIgnorePath,
		License:       licensePath,
		Readme:        readMePath,
		DefaultBranch: defaultBranch,
		TrustModel:    "",
	})
	if err != nil {
		return nil, errors.New(fmt.Sprintf("failed to create repository: %v for user: %v, err: %v", repoName, userName, err))
	}

	if res.StatusCode != 201 && res.StatusCode != 200 {
		return nil, errors.New(fmt.Sprintf("failed to create repository: %v for user: %v, err: non 201/200 status code returned: %v", repoName, userName, res.StatusCode))
	}

	return repo, nil
}

func (v *VCSClient) DeleteRepo(ownerName string, repoName string) error {
	res, err := v.GiteaClient.DeleteRepo(ownerName, repoName)
	if err != nil {
		return errors.New(fmt.Sprintf("failed to delete repo: %v from user: %v, err: %v", repoName, ownerName, err))
	}

	if res.StatusCode != 200 && res.StatusCode != 202 {
		return errors.New(fmt.Sprintf("failed to delete repo: %v from user: %v, err: non 200/202 status code returned from delete: %v", repoName, ownerName, res.StatusCode))
	}

	return nil
}

// use snowflake model id for sourceid
func (v *VCSClient) CreateUser(loginName string, userName string, fullName string, email string,
	password string) (*gitea.User, error) {
	visibility := gitea.VisibleTypePrivate
	changePassword := false
	user, res, err := v.GiteaClient.AdminCreateUser(gitea.CreateUserOption{
		SourceID:           0,
		LoginName:          loginName,
		Username:           userName,
		FullName:           fullName,
		Email:              email,
		Password:           password,
		MustChangePassword: &changePassword,
		SendNotify:         false,
		Visibility:         &visibility,
	})
	if err != nil {
		return nil, errors.New(fmt.Sprintf("failed to create user: %v, err: %v", email, err))
	}

	if res.StatusCode != 201 && res.StatusCode != 200 {
		return nil, errors.New(fmt.Sprintf("failed to create user, err: status code is not 200/201 status code: %v", res.StatusCode))
	}

	if user == nil {
		return nil, errors.New(fmt.Sprintf("failed to create user, err: user returned nil"))
	}

	return user, nil

}

func (v *VCSClient) EditUser(loginName string, userName string, fullName string, email string,
	password string) error {
	visibility := gitea.VisibleTypePrivate
	changePassword := false
	res, err := v.GiteaClient.AdminEditUser(userName, gitea.EditUserOption{
		SourceID:           0,
		LoginName:          loginName,
		FullName:           &fullName,
		Email:              &email,
		Password:           password,
		MustChangePassword: &changePassword,
		Visibility:         &visibility,
	})
	if err != nil {
		return errors.New(fmt.Sprintf("failed to edit pre-existing user: %v, err: %v", email, err))
	}

	if res.StatusCode != 201 && res.StatusCode != 200 {
		return errors.New(fmt.Sprintf("failed to edit pre-existing user, err: status code is not 200/201 status code: %v", res.StatusCode))
	}

	return nil

}

func (v *VCSClient) DeleteUser(userName string) error {
	res, err := v.GiteaClient.AdminDeleteUser(userName)
	if err != nil {
		return errors.New(fmt.Sprintf("failed to delete user, err: %v", err))
	}

	if res.StatusCode != 200 && res.StatusCode != 202 {
		return errors.New(fmt.Sprintf("failed to delete user, err: status code is not 200/202 status coder: %v", res.StatusCode))
	}

	return nil
}

func (v *VCSClient) Branch(ownerName string, repoName string, branchName string, oldBranchName string) (*gitea.Branch, error) {
	branch, res, err := v.GiteaClient.CreateBranch(ownerName, repoName, gitea.CreateBranchOption{
		BranchName:    branchName,
		OldBranchName: oldBranchName,
	})
	if err != nil {
		return nil, errors.New(fmt.Sprintf("failed to create branch: %v on repo: %v from user: %v, err: %v",
			branchName, repoName, ownerName, err))
	}

	if res.StatusCode != 200 && res.StatusCode != 201 {
		return nil, errors.New(fmt.Sprintf("failed to create branch: %v on repo: %v from user: %v, "+
			"err: non 200/201 status code returned status code: %v", branchName, repoName, ownerName, res.StatusCode))
	}

	if branch == nil {
		return nil, errors.New(fmt.Sprintf("failed to create branch: %v on repo: %v from user: %v, "+
			"err: branch returned nil", branchName, repoName, ownerName))
	}

	return branch, nil
}

func (v *VCSClient) DeleteBranch(ownerName string, repoName string, branchName string) error {
	deleted, res, err := v.GiteaClient.DeleteRepoBranch(ownerName, repoName, branchName)
	if err != nil {
		return errors.New(fmt.Sprintf("failed to delete branch: %v from user: %v on repo: %v, err: %v",
			branchName, ownerName, repoName, err))
	}

	if res.StatusCode != 200 && res.StatusCode != 202 {
		return errors.New(fmt.Sprintf("failed to delete branch: %v from user: %v on repo: %v, "+
			"err: status code is not 200/202 status coder: %v", branchName, ownerName, repoName, res.StatusCode))
	}

	if !deleted {
		return errors.New(fmt.Sprintf("failed to delete branch: %v from user: %v on repo: %v, "+
			"err: boolean returned false for delete", branchName, ownerName, repoName))
	}

	return nil
}

func (v *VCSClient) ListRepoBranches(ownerName string, repoName string,
	pageNum *int, pageSize *int) ([]*gitea.Branch, error) {
	pNum := 0
	pSize := 1000
	if pageNum != nil {
		pNum = *pageNum
	}

	if pageSize != nil {
		pSize = *pageSize
	}

	branches, res, err := v.GiteaClient.ListRepoBranches(ownerName, repoName, gitea.ListRepoBranchesOptions{
		gitea.ListOptions{
			Page:     pNum,
			PageSize: pSize,
		},
	})
	if err != nil {
		return nil, errors.New(fmt.Sprintf("failed to list branches for repo: %v from user: %v, err: %v",
			repoName, ownerName, err))
	}

	if res.StatusCode != 200 {
		return nil, errors.New(fmt.Sprintf("failed to list branches for repo: %v from user: %v,"+
			" err: non 200 status code returned: %v", repoName, ownerName, res.StatusCode))
	}

	if branches == nil || len(branches) < 1 {
		return nil, errors.New(fmt.Sprintf("failed to list branches for repo: %v from user: %v, "+
			"err: no branches found or returned", repoName, ownerName))
	}

	return branches, nil
}

func (v *VCSClient) ListRepoContents(ownerName string, repoName string, ref string,
	directory string) ([]*gitea.ContentsResponse, error) {
	content, res, err := v.GiteaClient.ListContents(ownerName, repoName, ref, directory)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("failed to list contents of repo: %v made by user: %v, err: %v",
			repoName, ownerName, err))
	}

	if res.StatusCode != 200 {
		return nil, errors.New(fmt.Sprintf("failed to list contents of repo: %v made by user: %v, err: non 200 "+
			"status code returned: %v", repoName, ownerName, res.StatusCode))
	}

	if content == nil || len(content) < 1 {
		return nil, errors.New(fmt.Sprintf("failed to list contents of repo: %v made by user: %v, "+
			"err: no data returned", repoName, ownerName))
	}

	for _, c := range content {
		if c.DownloadURL == nil || *c.DownloadURL == "" {
			continue
		}

		nUrl, err := url.Parse(*c.DownloadURL)
		if err != nil {
			return nil, errors.New(fmt.Sprintf("failed to parse url inside of ListRepoContents, err: %v", err))
		}
		c.DownloadURL = &nUrl.Path
	}

	return content, nil
}

func (v *VCSClient) ListAllRepos(ownerName string, pageNum *int, pageSize *int) ([]*gitea.Repository, error) {
	pNum := 0
	pSize := 1000

	if pageNum != nil {
		pNum = *pageNum
	}

	if pageSize != nil {
		pSize = *pageSize
	}

	repos, res, err := v.GiteaClient.ListUserRepos(ownerName, gitea.ListReposOptions{
		gitea.ListOptions{
			Page:     pNum,
			PageSize: pSize,
		},
	})
	if err != nil {
		return nil, errors.New(fmt.Sprintf("failed to list all repos for user: %v, err: %v", ownerName, err))
	}

	if res.StatusCode != 200 {
		return nil, errors.New(fmt.Sprintf("failed to list all repos for user: %v, "+
			"err: non 200 status code returned: %v", ownerName, res.StatusCode))
	}

	if repos == nil || len(repos) < 1 {
		return nil, errors.New(fmt.Sprintf("failed to list all repos for user: %v, "+
			"err: no repos found or results returned nil", ownerName))
	}

	return repos, nil
}

func (v *VCSClient) Fork(ownerName string, repoName string, orgName *string) (*gitea.Repository, error) {
	repo, res, err := v.GiteaClient.CreateFork(ownerName, repoName, gitea.CreateForkOption{
		Organization: orgName,
	})
	if err != nil {
		return nil, errors.New(fmt.Sprintf("failed to create fork on repo: %v from user: %v,"+
			" err: %v", repoName, ownerName, err))
	}

	if res.StatusCode != 200 && res.StatusCode != 201 {
		return nil, errors.New(fmt.Sprintf("failed to create fork on repo: %v from user: %v,"+
			" err: non 200/201 status code returned: %v", repoName, ownerName, res.StatusCode))
	}

	if repo == nil {
		return nil, errors.New(fmt.Sprintf("failed to create fork on repo: %v from user: %v,"+
			" err: forked repo returned nil", repoName, ownerName))
	}

	return repo, nil
}
