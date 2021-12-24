package binaryfile

import (
	"encoding/json"
	"errors"
	"regexp"
	"strings"

	mygit "github.com/MrWong99/gitscanner/git"
	"github.com/MrWong99/gitscanner/utils"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"gorm.io/datatypes"
)

type BinarySearchCheck struct {
	cfg map[string]interface{}
}

func (*BinarySearchCheck) String() string {
	return "SearchBinaries"
}

func (bins *BinarySearchCheck) GetConfig() map[string]interface{} {
	return bins.cfg
}

func (bins *BinarySearchCheck) SetConfig(cfg map[string]interface{}) error {
	pat, ok := cfg["branchPattern"]
	branchPattern := ".*"
	if ok {
		switch strPat := pat.(type) {
		case string:
			if _, err := utils.ExtractPattern(strPat); err != nil {
				return err
			}
			branchPattern = strPat
		default:
			return errors.New("given configuration didn't have a string as 'branchPattern'")
		}
	}
	bins.cfg = map[string]interface{}{
		"branchPattern": branchPattern,
	}
	return nil
}

func (check *BinarySearchCheck) getPat() *regexp.Regexp {
	pat, _ := check.cfg["branchPattern"].(string)
	return regexp.MustCompile(pat)
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
		return datatypes.JSON([]byte(`{"err": "` + strings.ReplaceAll(err.Error(), "\\", "\\\\") + `"}`))
	}
	return datatypes.JSON(bytes)
}
