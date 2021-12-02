package utils

import (
	"reflect"
	"regexp"
	"runtime"

	"github.com/go-git/go-git/v5"
)

var config *GlobalConfig

type GlobalConfig struct {
	BranchPattern *regexp.Regexp
}

// Initialize the global configuration with given config.
func InitConfig(cfg *GlobalConfig) {
	config = cfg
}

// Returns the name and package/module path of the given function.
func FunctionName(i interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}

// Returns the name of the given repository. Usually derived from its first remote.
func RepoName(repo *git.Repository) string {
	remotes, _ := repo.Remotes()
	if len(remotes) > 0 {
		return remotes[0].Config().URLs[0]
	} else {
		if wt, err := repo.Worktree(); err == nil {
			return wt.Filesystem.Root()
		}
	}
	return "Unknown Repo"
}

// Retrieve the global configuration.
func Config() GlobalConfig {
	return *config
}
