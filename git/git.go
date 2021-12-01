package git

import (
	"errors"
	"os"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"github.com/go-git/go-git/v5/storage/memory"
)

var sshKeyAuth transport.AuthMethod

// Set the ssh private key in PEM format to use for authentication. If needed provide the password.
func InitSshKey(privateKeyFile, password string) error {
	_, err := os.Stat(privateKeyFile)
	if err != nil {
		return err
	}
	publicKeys, err := ssh.NewPublicKeysFromFile("git", privateKeyFile, password)
	if err != nil {
		return err
	}
	sshKeyAuth = publicKeys
	return nil
}

var httpAuth transport.AuthMethod

// Set HTTP basic authentication user and password.
func InitHttpBasicAuth(username, password string) error {
	httpAuth = &http.BasicAuth{
		Username: username,
		Password: password,
	}
	return nil
}

// Loads a given repository in a memory db and return it. The url can be given in these formats:
// file://<path> -> path to repository in local filesystem
// git@<remote repo> -> ssh repository clone URL
// http(s)://<remote repo>
func LoadRepoInMemory(url string) (*git.Repository, error) {
	if strings.HasPrefix(url, "http") {
		return git.Clone(memory.NewStorage(), nil, &git.CloneOptions{
			URL:               url,
			Auth:              httpAuth,
			Tags:              git.NoTags,
			RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
		})
	} else if strings.HasPrefix(url, "git@") {
		return git.Clone(memory.NewStorage(), nil, &git.CloneOptions{
			URL:               url,
			Auth:              sshKeyAuth,
			Tags:              git.NoTags,
			RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
		})
	} else if strings.HasPrefix(url, "file://") {
		return git.PlainOpen(strings.TrimPrefix(url, "file://"))
	} else {
		return nil, errors.New("Given url does not match any known url type: " + url)
	}
}
