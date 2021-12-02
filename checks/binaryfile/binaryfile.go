package binaryfile

import (
	mygit "github.com/MrWong99/gitscanner/git"
	"github.com/MrWong99/gitscanner/utils"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
)

func SearchBinaries(wrapRepo *mygit.ClonedRepo, output chan<- utils.SingleCheck) error {
	defer close(output)
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
			return nil
		}
		tree, err := commit.Tree()
		if err != nil {
			return err
		}
		tree.Files().ForEach(func(f *object.File) error {
			if isBin, err := f.IsBinary(); err == nil && isBin {
				output <- utils.SingleCheck{
					Origin:    f.Name,
					Branch:    branchRef.Name().String(),
					CheckName: utils.FunctionName(SearchBinaries),
					AdditionalInfo: map[string]interface{}{
						"filesize": utils.ByteCountDecimal(f.Size),
						"filemode": f.Mode,
					},
				}
			}
			return nil
		})
		return nil
	})
}
