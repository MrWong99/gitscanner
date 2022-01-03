package commitmeta

import (
	"encoding/json"
	"errors"
	"math"
	"regexp"
	"strconv"
	"strings"

	mygit "github.com/MrWong99/gitscanner/git"
	"github.com/MrWong99/gitscanner/utils"
	"github.com/go-git/go-git/v5/plumbing/object"
	"gorm.io/datatypes"
)

type CommitMetaInfoCheck struct {
	cfg map[string]interface{}
}

func (*CommitMetaInfoCheck) String() string {
	return "CheckCommitMetaInformation"
}

func (bins *CommitMetaInfoCheck) GetConfig() map[string]interface{} {
	return bins.cfg
}

func (bins *CommitMetaInfoCheck) SetConfig(cfg map[string]interface{}) error {
	pat, ok := cfg["emailPattern"]
	emailPattern := ".*"
	if ok {
		switch strPat := pat.(type) {
		case string:
			if _, err := utils.ExtractPattern(strPat); err != nil {
				return err
			}
			emailPattern = strPat
		default:
			return errors.New("given configuration for  didn't have a string as 'emailPattern'")
		}
	}

	pat, ok = cfg["namePattern"]
	namePattern := ".*"
	if ok {
		switch strPat := pat.(type) {
		case string:
			if _, err := utils.ExtractPattern(strPat); err != nil {
				return err
			}
			namePattern = strPat
		default:
			return errors.New("given configuration for didn't have a string as 'namePattern'")
		}
	}

	t, ok := cfg["commitSizeThresholdByte"]
	if !ok {
		bins.cfg = map[string]interface{}{
			"emailPattern":            emailPattern,
			"namePattern":             namePattern,
			"commitSizeThresholdByte": math.MaxInt64,
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
		"emailPattern":            emailPattern,
		"namePattern":             namePattern,
		"commitSizeThresholdByte": t64,
	}
	return nil
}

func (bins *CommitMetaInfoCheck) getPat(name string) *regexp.Regexp {
	pat, _ := bins.cfg[name].(string)
	return regexp.MustCompile(pat)
}

func (bins *CommitMetaInfoCheck) getSizeThreshold() int64 {
	i, _ := bins.cfg["commitSizeThresholdByte"].(int64)
	return i
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
		commitSize := getCommitSize(commit)
		if commitSize >= check.getSizeThreshold() {
			output <- utils.SingleCheck{
				Origin:         "Commit " + commit.Hash.String(),
				Branch:         "",
				CheckName:      check.String(),
				AdditionalInfo: getAdditionalInfo(commit, commitSize),
			}
		}
		if !emailPat.MatchString(commit.Author.Email) || !emailPat.MatchString(commit.Committer.Email) ||
			!namePat.MatchString(commit.Author.Name) || !namePat.MatchString(commit.Committer.Name) {
			output <- utils.SingleCheck{
				Origin:         "Commit " + commit.Hash.String(),
				Branch:         "",
				CheckName:      check.String(),
				AdditionalInfo: getAdditionalInfo(commit, commitSize),
			}
		}
		return nil
	})
}

func getAdditionalInfo(c *object.Commit, commitSize int64) datatypes.JSON {
	bytes, err := json.Marshal(map[string]interface{}{
		"commitMessage":   c.Message,
		"authorName":      c.Author.Name,
		"authorEmail":     c.Author.Email,
		"commiterName":    c.Committer.Name,
		"commiterEmail":   c.Committer.Email,
		"commitSize":      utils.ByteCountDecimal(commitSize),
		"numberOfParents": strconv.Itoa(c.NumParents()),
	})
	if err != nil {
		return datatypes.JSON([]byte(`{"err": "` + strings.ReplaceAll(err.Error(), "\\", "\\\\") + `"}`))
	}
	return datatypes.JSON(bytes)
}

func getCommitSize(c *object.Commit) int64 {
	parent, err := c.Parent(0)
	if err != nil {
		if files, err := c.Files(); err == nil {
			return calculateFilesSize(files)
		} else {
			return 0
		}
	}
	patch, err := c.Patch(parent)
	if err != nil {
		return -1
	}
	var size int64
	for _, fp := range patch.FilePatches() {
		for _, chunk := range fp.Chunks() {
			size += int64(len(chunk.Content()))
		}
	}
	return size
}

func calculateFilesSize(files *object.FileIter) int64 {
	var size int64
	files.ForEach(func(f *object.File) error {
		size += f.Size
		return nil
	})
	return size
}
