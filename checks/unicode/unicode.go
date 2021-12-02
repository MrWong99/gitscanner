package unicode

import (
	"fmt"
	"io/ioutil"
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

func SearchUnicode(wrapRepo *mygit.ClonedRepo) error {
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
			searchForIllegal(t, wrapRepo, branchRef)
		}(tree)
		return nil
	})
	wg.Wait()
	return err
}

func searchForIllegal(t *object.Tree, repo *mygit.ClonedRepo, branchRef *plumbing.Reference) {
	t.Files().ForEach(func(f *object.File) error {
		reader, err := f.Reader()
		if err != nil {
			fmt.Printf("Could not open file %s in repo %s with branch %s\n", f.Name, utils.RepoName(repo.Repo), branchRef.Name())
			return nil
		}
		defer reader.Close()
		content, err := ioutil.ReadAll(reader)
		if err != nil {
			fmt.Printf("Could not open file %s in repo %s with branch %s\n", f.Name, utils.RepoName(repo.Repo), branchRef.Name())
			return nil
		}
		if !utf8.Valid(content) {
			return nil
		}
		for _, illegalRune := range illegalUnicodeChars {
			if strings.ContainsRune(string(content), illegalRune) {
				fmt.Printf("Found file '%s' that contains illegal unicode character %s in repo %s, branch %s\n",
					f.Name, strconv.QuoteRuneToASCII(illegalRune), utils.RepoName(repo.Repo), branchRef.Name())
			}
		}
		return nil
	})
}
