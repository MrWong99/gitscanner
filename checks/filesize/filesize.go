package filesize

import (
	"encoding/json"
	"errors"
	"math"
	"regexp"
	"strings"

	mygit "github.com/MrWong99/gitscanner/git"
	"github.com/MrWong99/gitscanner/utils"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"gorm.io/datatypes"
)

type FilesizeSearchCheck struct {
	cfg map[string]interface{}
}

func (*FilesizeSearchCheck) String() string {
	return "SearchBigFiles"
}

func (bins *FilesizeSearchCheck) GetConfig() map[string]interface{} {
	return map[string]interface{}{
		"branchPattern":         bins.getPat().String(),
		"filesizeThresholdByte": bins.getThreshold(),
	}
}

func (bins *FilesizeSearchCheck) SetConfig(cfg map[string]interface{}) error {
	pat, ok := cfg["branchPattern"]
	branchPattern := regexp.MustCompile(".*")
	if ok {
		switch strPat := pat.(type) {
		case string:
			pat, err := utils.ExtractPattern(strPat)
			if err != nil {
				return err
			}
			branchPattern = pat
		default:
			return errors.New("given configuration didn't have a string as 'branchPattern'")
		}
	}

	t, ok := cfg["filesizeThresholdByte"]
	if !ok {
		bins.cfg = map[string]interface{}{
			"branchPattern":         branchPattern,
			"filesizeThresholdByte": int64(81920),
		}
		return nil
	}
	var t64 int64
	switch threshold := t.(type) {
	case int:
		t64 = int64(threshold)
	case int8:
		t64 = int64(threshold)
	case int16:
		t64 = int64(threshold)
	case int32:
		t64 = int64(threshold)
	case int64:
		t64 = int64(threshold)
	case float32:
		t64 = int64(math.Trunc(float64(threshold)))
	case float64:
		t64 = int64(math.Trunc(threshold))
	default:
		return errors.New("given configuration didn't have an integer for optional 'commitSizeThresholdByte'")
	}
	bins.cfg = map[string]interface{}{
		"branchPattern":         branchPattern,
		"filesizeThresholdByte": t64,
	}
	return nil
}

func (bins *FilesizeSearchCheck) getPat() *regexp.Regexp {
	pat, _ := bins.cfg["branchPattern"].(*regexp.Regexp)
	return pat
}

func (bins *FilesizeSearchCheck) getThreshold() int64 {
	i, _ := bins.cfg["filesizeThresholdByte"].(int64)
	return i
}

func (check *FilesizeSearchCheck) Check(wrapRepo *mygit.ClonedRepo, output chan<- utils.SingleCheck) error {
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
			if f.Size >= check.getThreshold() {
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
