package unicode

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"unicode/utf8"

	mygit "github.com/MrWong99/gitscanner/git"
	"github.com/MrWong99/gitscanner/utils"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"gorm.io/datatypes"
)

var illegalUnicodeChars = []rune{
	'\u202a',
	'\u202b',
	'\u202d',
	'\u202e',
	'\u2066',
	'\u2067',
	'\u2068',
	'\u2069',
	'\u202c',
}

type UnicodeCharacterSearch struct {
	cfg map[string]interface{}
}

func (*UnicodeCharacterSearch) String() string {
	return "SearchIllegalUnicodeCharacters"
}

func (bins *UnicodeCharacterSearch) GetConfig() map[string]interface{} {
	return bins.cfg
}

func (bins *UnicodeCharacterSearch) SetConfig(cfg map[string]interface{}) error {
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

func (check *UnicodeCharacterSearch) getPat() *regexp.Regexp {
	pat, _ := check.cfg["branchPattern"].(string)
	return regexp.MustCompile(pat)
}

func (check *UnicodeCharacterSearch) Check(wrapRepo *mygit.ClonedRepo, output chan<- utils.SingleCheck) error {
	defer close(output)
	repo := wrapRepo.Repo
	branchIt, err := repo.References()
	if err != nil {
		return err
	}
	wg := new(sync.WaitGroup)
	err = branchIt.ForEach(func(branchRef *plumbing.Reference) error {
		if !(branchRef.Name().IsBranch() || branchRef.Name().IsRemote()) || !check.getPat().MatchString(branchRef.Name().String()) {
			return nil
		}
		commit, err := repo.CommitObject(branchRef.Hash())
		if err != nil {
			return err
		}
		tree, err := commit.Tree()
		if err != nil {
			return nil
		}
		wg.Add(1)
		go func(t *object.Tree) {
			defer wg.Done()
			check.searchForIllegal(t, wrapRepo, branchRef, output)
		}(tree)
		return nil
	})
	wg.Wait()
	return err
}

func getAdditionalInfo(f *object.File, illegalChar rune) datatypes.JSON {
	bytes, err := json.Marshal(map[string]interface{}{
		"filesize":  utils.ByteCountDecimal(f.Size),
		"filemode":  f.Mode.String(),
		"character": strings.ReplaceAll(strconv.QuoteRuneToASCII(illegalChar), "\\", "\\\\"),
	})
	if err != nil {
		return datatypes.JSON([]byte(`{"err": "` + strings.ReplaceAll(err.Error(), "\\", "\\\\") + `"}`))
	}
	return datatypes.JSON(bytes)
}

func (check *UnicodeCharacterSearch) searchForIllegal(t *object.Tree, repo *mygit.ClonedRepo, branchRef *plumbing.Reference, output chan<- utils.SingleCheck) {
	t.Files().ForEach(func(f *object.File) error {
		reader, err := f.Reader()
		if err != nil {
			log.Printf("Could not open file %s in repo %s with branch %s\n", f.Name, utils.RepoName(repo.Repo), branchRef.Name())
			return nil
		}
		defer reader.Close()
		content, err := ioutil.ReadAll(reader)
		if err != nil {
			log.Printf("Could not open file %s in repo %s with branch %s\n", f.Name, utils.RepoName(repo.Repo), branchRef.Name())
			return nil
		}
		if !utf8.Valid(content) {
			return nil
		}
		for _, illegalRune := range illegalUnicodeChars {
			if strings.ContainsRune(string(content), illegalRune) {
				output <- utils.SingleCheck{
					Origin:         f.Name,
					Branch:         branchRef.Name().String(),
					CheckName:      check.String(),
					AdditionalInfo: getAdditionalInfo(f, illegalRune),
				}
			}
		}
		return nil
	})
}
