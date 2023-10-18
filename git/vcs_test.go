package git

import (
	"fmt"
	"testing"
)

func TestCreateVCSClient(t *testing.T) {
	vcs, err := CreateVCSClient("http://0.0.0.0:4541", "test", "test", true)
	if err != nil {
		t.Errorf("failed to create vcs client, err: %v", err)
		return
	}

	if vcs == nil {
		t.Errorf("failed to create vcs client, err: client returned nil")
		return
	}

}

func TestVCSClient_CreateRepo(t *testing.T) {
	vcs, err := CreateVCSClient("http://0.0.0.0:4541", "test", "test", true)
	if err != nil {
		t.Errorf("failed to create vcs client, err: %v", err)
		return
	}

	if vcs == nil {
		t.Errorf("failed to create vcs client, err: client returned nil")
		return
	}

	repo, err := vcs.CreateRepo("test", "testingRepo", "", false,
		"", "", "", "")
	if err != nil {
		t.Errorf("failed to create repo, err: %v", err)
		return
	}

	if repo == nil {
		t.Errorf("failed to create repo, err: repo returned nil")
		return
	}

}

func TestVCSClient_ListRepoContents(t *testing.T) {
	vcs, err := CreateVCSClient("http://0.0.0.0:4541", "test", "test", true)
	if err != nil {
		t.Errorf("failed to create vcs client, err: %v", err)
		return
	}

	if vcs == nil {
		t.Errorf("failed to create vcs client, err: client returned nil")
		return
	}

	contents, err := vcs.ListRepoContents("test", "testingRepo", "", "/")
	if err != nil {
		t.Errorf("failed to list repo contents, err: %v", err)
		return
	}

	for _, c := range contents {
		fmt.Println(c)
		fmt.Println(c.Path)
		if c.DownloadURL != nil {
			fmt.Println(*c.DownloadURL)
		}

	}
}
