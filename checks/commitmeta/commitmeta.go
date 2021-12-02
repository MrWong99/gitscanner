package commitmeta

import (
	"fmt"

	mygit "github.com/MrWong99/gitscanner/git"
	"github.com/MrWong99/gitscanner/utils"
	"github.com/go-git/go-git/v5/plumbing/object"
)

func CheckCommits(wrapRepo *mygit.ClonedRepo) error {
	commits, err := wrapRepo.Repo.CommitObjects()
	if err != nil {
		return err
	}
	cfg := utils.Config()
	return commits.ForEach(func(commit *object.Commit) error {
		if !cfg.EmailPattern.MatchString(commit.Author.Email) {
			fmt.Printf("Commit %s in repo %s has illegal author email '%s' with commit message %s", commit.Hash, utils.RepoName(wrapRepo.Repo), commit.Author.Email, commit.Message)
		}
		if !cfg.EmailPattern.MatchString(commit.Committer.Email) {
			fmt.Printf("Commit %s in repo %s has illegal commiter email '%s' with commit message %s", commit.Hash, utils.RepoName(wrapRepo.Repo), commit.Committer.Email, commit.Message)
		}
		if !cfg.NamePattern.MatchString(commit.Author.Name) {
			fmt.Printf("Commit %s in repo %s has illegal author name '%s' with commit message %s", commit.Hash, utils.RepoName(wrapRepo.Repo), commit.Author.Name, commit.Message)
		}
		if !cfg.NamePattern.MatchString(commit.Committer.Name) {
			fmt.Printf("Commit %s in repo %s has illegal commiter name '%s' with commit message %s", commit.Hash, utils.RepoName(wrapRepo.Repo), commit.Committer.Name, commit.Message)
		}
		return nil
	})
}
