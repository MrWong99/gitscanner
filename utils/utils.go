package utils

import (
	"fmt"
	"reflect"
	"regexp"
	"runtime"
	"time"

	"github.com/go-git/go-git/v5"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type SingleCheck struct {
	ID                        uint           `json:"id" gorm:"primarykey"`
	Origin                    string         `json:"origin"`
	Branch                    string         `json:"branch"`
	CheckName                 string         `json:"checkName"`
	Acknowledged              bool           `json:"acknowledged"`
	AdditionalInfo            datatypes.JSON `json:"additionalInfo"`
	CheckResultConsolidatedID uint           `json:"-"`
	CreatedAt                 time.Time      `json:"-"`
	UpdatedAt                 time.Time      `json:"-"`
	DeletedAt                 gorm.DeletedAt `json:"-" gorm:"index"`
}

type CheckResultConsolidated struct {
	ID         uint           `json:"id" gorm:"primarykey"`
	Date       time.Time      `json:"date"`
	Repository string         `json:"repository"`
	Error      string         `json:"error"`
	Checks     []SingleCheck  `json:"checks" gorm:"foreignKey:CheckResultConsolidatedID"`
	CreatedAt  time.Time      `json:"-"`
	UpdatedAt  time.Time      `json:"-"`
	DeletedAt  gorm.DeletedAt `json:"-" gorm:"index"`
}

type SearchRequestBody struct {
	Path       string   `json:"path"`
	CheckNames []string `json:"checkNames"`
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

func ExtractPattern(input string) (*regexp.Regexp, error) {
	if input == "" {
		return regexp.MustCompile(".*"), nil
	}
	pat, err := regexp.Compile(input)
	if err != nil {
		return nil, err
	}
	return pat, nil
}
