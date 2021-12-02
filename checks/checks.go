package checks

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/MrWong99/gitscanner/checks/binaryfile"
	"github.com/MrWong99/gitscanner/checks/commitmeta"
	"github.com/MrWong99/gitscanner/checks/unicode"
	mygit "github.com/MrWong99/gitscanner/git"
	"github.com/MrWong99/gitscanner/utils"
)

type RepoCheckFunc func(*mygit.ClonedRepo, chan<- utils.SingleCheck) error

var repoChecks = []RepoCheckFunc{
	binaryfile.SearchBinaries,
	unicode.SearchUnicode,
	commitmeta.CheckCommitAuthor,
}

func CheckAllRepositories(repos []string) []*utils.CheckResultConsolidated {
	var results []*utils.CheckResultConsolidated
	for _, path := range repos {
		results = append(results, consolidateChecks(path))
	}
	return results
}

func consolidateChecks(path string) *utils.CheckResultConsolidated {
	repo, err := mygit.CloneRepo(path)
	if err != nil {
		return &utils.CheckResultConsolidated{
			Date:       time.Now(),
			Repository: path,
			Error:      err.Error(),
		}
	}
	res := repositoryCheck(repo)
	if err := repo.Cleanup(); err != nil {
		fmt.Fprintf(os.Stderr, "Error while cleaning up repo %s: %v", utils.RepoName(repo.Repo), err)
	}
	return res
}

func repositoryCheck(repo *mygit.ClonedRepo) *utils.CheckResultConsolidated {
	res := &utils.CheckResultConsolidated{
		Date:       time.Now(),
		Repository: utils.RepoName(repo.Repo),
		Checks:     []utils.SingleCheck{},
	}
	waitGroup := new(sync.WaitGroup)
	waitGroup.Add(len(repoChecks))
	for _, check := range repoChecks {
		checkChan := make(chan utils.SingleCheck)
		go func(r *mygit.ClonedRepo, checkFn RepoCheckFunc, outputs chan<- utils.SingleCheck) {
			err := checkFn(repo, checkChan)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error while checking repository %s with function %s: %v\n", utils.RepoName(r.Repo), utils.FunctionName(checkFn), err)
			}
		}(repo, check, checkChan)
		go awaitCheckResults(checkChan, res, waitGroup)
	}
	waitGroup.Wait()
	return res
}

func awaitCheckResults(inputs <-chan utils.SingleCheck, results *utils.CheckResultConsolidated, waitGroup *sync.WaitGroup) {
	defer waitGroup.Done()
	for {
		res, ok := <-inputs
		if ok {
			results.Checks = append(results.Checks, res)
		} else {
			break
		}
	}
}
