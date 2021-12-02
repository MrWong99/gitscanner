package unicode

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"sync"
	"unicode/utf8"

	mygit "github.com/MrWong99/gitscanner/git"
	"github.com/MrWong99/gitscanner/utils"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
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

func SearchUnicode(wrapRepo *mygit.ClonedRepo, output chan<- utils.SingleCheck) error {
	defer close(output)
	repo := wrapRepo.Repo
	branchIt, err := repo.References()
	if err != nil {
		return err
	}
	wg := new(sync.WaitGroup)
	err = branchIt.ForEach(func(branchRef *plumbing.Reference) error {
		if !(branchRef.Name().IsBranch() || branchRef.Name().IsRemote()) || !utils.Config().BranchPattern.MatchString(branchRef.Name().String()) {
			return nil
		}
		commit, err := repo.CommitObject(branchRef.Hash())
		if err != nil {
			return err
		}
		tree, err := commit.Tree()
		if err != nil {
			return err
		}
		wg.Add(1)
		go func(t *object.Tree) {
			defer wg.Done()
			searchForIllegal(t, wrapRepo, branchRef, output)
		}(tree)
		return nil
	})
	wg.Wait()
	return err
}

func searchForIllegal(t *object.Tree, repo *mygit.ClonedRepo, branchRef *plumbing.Reference, output chan<- utils.SingleCheck) {
	t.Files().ForEach(func(f *object.File) error {
		reader, err := f.Reader()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Could not open file %s in repo %s with branch %s\n", f.Name, utils.RepoName(repo.Repo), branchRef.Name())
			return err
		}
		defer reader.Close()
		content, err := ioutil.ReadAll(reader)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Could not open file %s in repo %s with branch %s\n", f.Name, utils.RepoName(repo.Repo), branchRef.Name())
			return err
		}
		if !utf8.Valid(content) {
			return nil
		}
		for _, illegalRune := range illegalUnicodeChars {
			if strings.ContainsRune(string(content), illegalRune) {
				output <- utils.SingleCheck{
					Origin:    f.Name,
					Branch:    branchRef.Name().String(),
					CheckName: utils.FunctionName(SearchUnicode),
					AdditionalInfo: map[string]interface{}{
						"character": strconv.QuoteRuneToASCII(illegalRune),
						"filesize":  utils.ByteCountDecimal(f.Size),
						"filemode":  f.Mode.String(),
					},
				}
			}
		}
		return nil
	})
}
