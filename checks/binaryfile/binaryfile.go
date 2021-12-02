package binaryfile

import (
	"fmt"

	mygit "github.com/MrWong99/gitscanner/git"
	"github.com/MrWong99/gitscanner/utils"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
)

func SearchBinaries(wrapRepo *mygit.ClonedRepo) error {
	repo := wrapRepo.Repo
	branchIt, err := repo.References()
	if err != nil {
		return err
	}
	return branchIt.ForEach(func(branchRef *plumbing.Reference) error {
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
		return tree.Files().ForEach(func(f *object.File) error {
			if isBin, err := f.IsBinary(); err == nil && isBin {
				fmt.Printf("Found binary file in repo %s, branch %s: %s\n", utils.RepoName(repo), branchRef.Name(), f.Name)
			}
			return nil
		})
	})
}
