package filesize

import (
	"errors"
	"math"
	"regexp"

	"github.com/MrWong99/gitscanner/checks"
	mygit "github.com/MrWong99/gitscanner/git"
	"github.com/MrWong99/gitscanner/utils"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"gorm.io/datatypes"
)

type FilesizeSearchCheck struct {
	cfg checks.CheckConfiguration
}

func (*FilesizeSearchCheck) String() string {
	return "SearchBigFiles"
}

func (bins *FilesizeSearchCheck) GetConfig() *checks.CheckConfiguration {
	return &bins.cfg
}

func (bins *FilesizeSearchCheck) SetConfig(c *checks.CheckConfiguration) error {
	cfg, err := c.ParseConfigMap()
	if err != nil {
		return err
	}
	pat, ok := cfg["branchPattern"]
	if !ok {
		return errors.New("Given configuration for '" + bins.String() + "' did not contain mandatory config 'branchPattern'!")
	}
	threshold, ok := cfg["filesizeThresholdByte"]
	if !ok {
		return errors.New("Given configuration for '" + bins.String() + "' did not contain mandatory config 'filesizeThresholdByte'!")
	}
	switch threshold.(type) {
	case int, int8, int16, int32, int64, float32, float64:
	default:
		return errors.New("Given configuration for '" + bins.String() + "' didn't have a int as 'filesizeThresholdByte'!")
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

func (bins *FilesizeSearchCheck) getPat() *regexp.Regexp {
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

func (bins *FilesizeSearchCheck) getThreshold() int64 {
	threshold, ok := bins.cfg.MustParseConfigMap()["filesizeThresholdByte"]
	if !ok {
		return 81920 // 80 KB
	}
	switch thInt := threshold.(type) {
	case int:
		return int64(thInt)
	case int8:
		return int64(thInt)
	case int16:
		return int64(thInt)
	case int32:
		return int64(thInt)
	case int64:
		return thInt
	case float32:
		return int64(math.Trunc(float64(thInt)))
	case float64:
		return int64(math.Trunc(thInt))
	default:
		return 81920 // 80 KB
	}
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
	return datatypes.JSON([]byte(`{"filesize": "` + utils.ByteCountDecimal(f.Size) + `", "filemode": "` + f.Mode.String() + `"}`))
}
