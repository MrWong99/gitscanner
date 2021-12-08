package binaryfile

import (
	"errors"
	"encoding/json"
	"regexp"

	"github.com/MrWong99/gitscanner/checks"
	mygit "github.com/MrWong99/gitscanner/git"
	"github.com/MrWong99/gitscanner/utils"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"gorm.io/datatypes"
)

type BinarySearchCheck struct {
	cfg checks.CheckConfiguration
}

func (*BinarySearchCheck) String() string {
	return "SearchBinaries"
}

func (bins *BinarySearchCheck) GetConfig() *checks.CheckConfiguration {
	return &bins.cfg
}

func (bins *BinarySearchCheck) SetConfig(c *checks.CheckConfiguration) error {
	cfg, err := c.ParseConfigMap()
	if err != nil {
		return err
	}
	pat, ok := cfg["branchPattern"]
	if !ok {
		return errors.New("Given configuration for '" + bins.String() + "' did not contain mandatory config 'branchPattern'!")
	}
	switch strPat := pat.(type) {
	case string:
		if _, err := utils.ExtractPattern(strPat); err != nil {
			return err
		}
		bins.cfg = *c
	default:
		return errors.New("Given configuration for '" + bins.String() + "' didn't have a string as 'branchPattern'!")
	}
	return nil
}

func (bins *BinarySearchCheck) getPat() *regexp.Regexp {
	pat, ok := bins.cfg.MustParseConfigMap()["branchPattern"]
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

func (check *BinarySearchCheck) Check(wrapRepo *mygit.ClonedRepo, output chan<- utils.SingleCheck) error {
	defer close(output)
	repo := wrapRepo.Repo
	branchIt, err := repo.References()
	if err != nil {
		return err
	}
	return branchIt.ForEach(func(branchRef *plumbing.Reference) error {
		if !(branchRef.Name().IsBranch() || branchRef.Name().IsRemote()) || !check.getPat().MatchString(branchRef.Name().String()) {
			return nil
		}
		commit, err := repo.CommitObject(branchRef.Hash())
		if err != nil {
			return nil
		}
		tree, err := commit.Tree()
		if err != nil {
			return nil
		}
		tree.Files().ForEach(func(f *object.File) error {
			if isBin, err := f.IsBinary(); err == nil && isBin {
				output <- utils.SingleCheck{
					Origin:         f.Name,
					Branch:         branchRef.Name().String(),
					CheckName:      check.String(),
					AdditionalInfo: getAdditionalInfo(f),
				}
			}
			return nil
		})
		return nil
	})
}

func getAdditionalInfo(f *object.File) datatypes.JSON {
	bytes, err := json.Marshal(map[string]interface{}{
		"filesize": utils.ByteCountDecimal(f.Size),
		"filemode": f.Mode.String(),
	})
	if err != nil {
		return datatypes.JSONF([]byte(`{"err": "` + strings.ReplaceAll(err.Error(), "\\", "\\\\") + `"}`)
	}
	return datatypes.JSON(bytes)
}
