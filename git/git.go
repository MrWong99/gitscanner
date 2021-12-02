package git

import (
	"errors"
	"os"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
)

var sshKeyAuth, httpAuth transport.AuthMethod

// Set the ssh private key in PEM format to use for authentication. If needed provide the password.
func InitSshKey(privateKeyFileContent []byte, password string) error {
	publicKeys, err := ssh.NewPublicKeys("git", privateKeyFileContent, password)
	if err != nil {
		return err
	}
	sshKeyAuth = publicKeys
	return nil
}

// Set HTTP basic authentication user and password.
func InitHttpBasicAuth(username, password string) error {
	httpAuth = &http.BasicAuth{
		Username: username,
		Password: password,
	}
	return nil
}

type ClonedRepo struct {
	Repo     *git.Repository
	LocalDir string
	Auth     transport.AuthMethod
}

// Deletes any local content in TempDir.
func (repo *ClonedRepo) Cleanup() error {
	if repo.LocalDir != "" && !strings.HasPrefix(repo.LocalDir, "file://") {
		return os.RemoveAll(repo.LocalDir)
	}
	return nil
}

// Clone a given repository and return it. The url can be given in these formats:
// file://<path> -> path to repository in local filesystem
// git@<remote repo> -> ssh repository clone URL
// http(s)://<remote repo>
func CloneRepo(url string) (*ClonedRepo, error) {
	clonedPath, err := os.MkdirTemp("", "git")
	if err != nil {
		return nil, err
	}
	var repo *git.Repository
	var auth transport.AuthMethod
	if strings.HasPrefix(url, "http") {
		auth = httpAuth
		repo, err = git.PlainClone(clonedPath, false, &git.CloneOptions{
			URL:      url,
			Auth:     httpAuth,
			Progress: os.Stderr,
		})
	} else if strings.HasPrefix(url, "git@") {
		auth = sshKeyAuth
		repo, err = git.PlainClone(clonedPath, false, &git.CloneOptions{
			URL:      url,
			Auth:     sshKeyAuth,
			Progress: os.Stderr,
		})
	} else if strings.HasPrefix(url, "file://") {
		clonedPath = url
		repo, err = git.PlainOpen(strings.TrimPrefix(url, "file://"))
	} else {
		return nil, errors.New("Given url does not match any known url type: " + url)
	}
	return wrap(clonedPath, repo, auth, err)
}

func wrap(tmpPath string, repo *git.Repository, auth transport.AuthMethod, err error) (*ClonedRepo, error) {
	if err != nil {
		if !strings.HasPrefix(tmpPath, "file://") {
			os.RemoveAll(tmpPath)
		}
		return nil, err
	}
	return &ClonedRepo{
		Repo:     repo,
		LocalDir: tmpPath,
		Auth:     auth,
	}, nil
}
