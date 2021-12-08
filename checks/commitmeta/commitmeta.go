package commitmeta

import (
	"errors"
	"encoding/json"
	"math"
	"regexp"
	"strconv"
	"strings"

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
	cfg, err := c.ParseConfigMap()
	if err != nil {
		return err
	}
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
	threshold, ok := cfg["commitSizeThresholdByte"]
	if !ok {
		return nil
	}
	switch threshold.(type) {
	case int, int8, int16, int32, int64, float32, float64:
	default:
		return errors.New("Given configuration for '" + bins.String() + "' didn't have a int for optional 'commitSizeThresholdByte'!")
	}
	return nil
}

func (bins *CommitMetaInfoCheck) getPat(name string) *regexp.Regexp {
	pat, ok := bins.cfg.MustParseConfigMap()[name]
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

func (bins *CommitMetaInfoCheck) getSizeThreshold() int64 {
	t, ok := bins.cfg.MustParseConfigMap()["commitSizeThresholdByte"]
	if !ok {
		return math.MaxInt64
	}
	switch threshold := t.(type) {
	case int:
		return int64(threshold)
	case int8:
		return int64(threshold)
	case int16:
		return int64(threshold)
	case int32:
		return int64(threshold)
	case int64:
		return int64(threshold)
	case float32:
		return int64(math.Trunc(float64(threshold)))
	case float64:
		return int64(math.Trunc(threshold))
	default:
		return math.MaxInt64
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
		"commitMessage": c.Message,
		"authorName": c.Author.Name,
		"authorEmail": c.Author.Email,
		"commiterName": c.Committer.Name,
		"commiterEmail": c.Committer.Email,
		"commitSize": utils.ByteCountDecimal(commitSize),
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
