package utils

import (
	"fmt"
	"reflect"
	"regexp"
	"runtime"
	"time"

	"github.com/go-git/go-git/v5"
)

var config *GlobalConfig

type GlobalConfig struct {
	BranchPattern *regexp.Regexp
	NamePattern   *regexp.Regexp
	EmailPattern  *regexp.Regexp
}

type SingleCheck struct {
	Origin         string                 `json:"origin"`
	Branch         string                 `json:"branch"`
	CheckName      string                 `json:"checkName"`
	Acknowledged   bool                   `json:"acknowledged"`
	AdditionalInfo map[string]interface{} `json:"additionalInfo"`
}

type CheckResultConsolidated struct {
	Date       time.Time     `json:"date"`
	Repository string        `json:"repository"`
	Error      string        `json:"error"`
	Checks     []SingleCheck `json:"checks"`
}

type SearchRequestBody struct {
	Path       string   `json:"path"`
	CheckNames []string `json:"checkNames"`
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

func ByteCountDecimal(b int64) string {
	const unit = 1000
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(b)/float64(div), "kMGTPE"[exp])
}
