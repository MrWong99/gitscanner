package commitmeta

import (
	mygit "github.com/MrWong99/gitscanner/git"
	"github.com/MrWong99/gitscanner/utils"
	"github.com/go-git/go-git/v5/plumbing/object"
)

func CheckCommitAuthor(wrapRepo *mygit.ClonedRepo, output chan<- utils.SingleCheck) error {
	defer close(output)
	commits, err := wrapRepo.Repo.CommitObjects()
	if err != nil {
		return err
	}
	cfg := utils.Config()
	return commits.ForEach(func(commit *object.Commit) error {
		if !cfg.EmailPattern.MatchString(commit.Author.Email) || !cfg.EmailPattern.MatchString(commit.Committer.Email) ||
			!cfg.NamePattern.MatchString(commit.Author.Name) || !cfg.NamePattern.MatchString(commit.Committer.Name) {
			output <- utils.SingleCheck{
				Origin:    "Commit " + commit.Hash.String(),
				Branch:    "",
				CheckName: utils.FunctionName(CheckCommitAuthor),
				AdditionalInfo: map[string]interface{}{
					"commitMessage":   commit.Message,
					"authorName":      commit.Author.Name,
					"authorEmail":     commit.Author.Email,
					"commiterName":    commit.Committer.Name,
					"commiterEmail":   commit.Committer.Email,
					"numberOfParents": commit.NumParents(),
				},
			}
		}
		return nil
	})
}
