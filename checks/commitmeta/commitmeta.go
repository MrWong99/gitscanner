package commitmeta

import (
	"errors"
	"regexp"
	"strconv"

	"github.com/MrWong99/gitscanner/checks"
	mygit "github.com/MrWong99/gitscanner/git"
	"github.com/MrWong99/gitscanner/utils"
	"github.com/go-git/go-git/v5/plumbing/object"
	"gorm.io/datatypes"
)

type CommitMetaInfoCheck struct {
	cfg checks.CheckConfiguration
}

func (*CommitMetaInfoCheck) String() string {
	return "CheckCommitMetaInformation"
}

func (bins *CommitMetaInfoCheck) GetConfig() *checks.CheckConfiguration {
	return &bins.cfg
}

func (bins *CommitMetaInfoCheck) SetConfig(c *checks.CheckConfiguration) error {
	cfg := c.GetConfig()
	pat, ok := cfg["emailPattern"]
	if !ok {
		return errors.New("Given configuration for '" + bins.String() + "' did not contain mandatory config 'emailPattern'!")
	}
	switch strPat := pat.(type) {
	case string:
		if _, err := utils.ExtractPattern(strPat); err != nil {
			return err
		}
	default:
		return errors.New("Given configuration for '" + bins.String() + "' didn't have a string as 'emailPattern'!")
	}
	pat, ok = cfg["namePattern"]
	if !ok {
		return errors.New("Given configuration for '" + bins.String() + "' did not contain mandatory config 'namePattern'!")
	}
	switch strPat := pat.(type) {
	case string:
		if _, err := utils.ExtractPattern(strPat); err != nil {
			return err
		}
		bins.cfg = *c
	default:
		return errors.New("Given configuration for '" + bins.String() + "' didn't have a string as 'namePattern'!")
	}
	return nil
}

func (bins *CommitMetaInfoCheck) getPat(name string) *regexp.Regexp {
	pat, ok := bins.cfg.GetConfig()[name]
	if !ok {
		return regexp.MustCompile(".*")
	}
	switch strPat := pat.(type) {
	case string:
		return regexp.MustCompile(strPat)
	default:
		return regexp.MustCompile(".*")
	}
}

func (check *CommitMetaInfoCheck) Check(wrapRepo *mygit.ClonedRepo, output chan<- utils.SingleCheck) error {
	defer close(output)
	emailPat := check.getPat("emailPattern")
	namePat := check.getPat("namePattern")
	commits, err := wrapRepo.Repo.CommitObjects()
	if err != nil {
		return err
	}
	return commits.ForEach(func(commit *object.Commit) error {
		if !emailPat.MatchString(commit.Author.Email) || !emailPat.MatchString(commit.Committer.Email) ||
			!namePat.MatchString(commit.Author.Name) || !namePat.MatchString(commit.Committer.Name) {
			output <- utils.SingleCheck{
				Origin:         "Commit " + commit.Hash.String(),
				Branch:         "",
				CheckName:      check.String(),
				AdditionalInfo: getAdditionalInfo(commit),
			}
		}
		return nil
	})
}

func getAdditionalInfo(c *object.Commit) datatypes.JSON {
	return datatypes.JSON([]byte(`{"commitMessage": "` + c.Message +
		`", "authorName": "` + c.Author.Name +
		`", "authorEmail": "` + c.Author.Email +
		`", "commiterName": "` + c.Committer.Name +
		`", "commiterEmail": "` + c.Committer.Email +
		`", "numberOfParents": ` + strconv.Itoa(c.NumParents()) + `}`))
}
