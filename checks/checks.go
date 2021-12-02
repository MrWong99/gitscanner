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

var RepoChecks = []RepoCheckFunc{
	binaryfile.SearchBinaries,
	unicode.SearchUnicode,
	commitmeta.CheckCommitAuthor,
}

func CheckAllRepositories(repos []string) []*utils.CheckResultConsolidated {
	var results []*utils.CheckResultConsolidated
	for _, path := range repos {
		results = append(results, consolidateChecks(path, RepoChecks))
	}
	return results
}

func CheckAllRepositoriesSpecificChecks(repos, checks []string) []*utils.CheckResultConsolidated {
	var results []*utils.CheckResultConsolidated
	checkFns := matchingChecks(checks)
	for _, path := range repos {
		results = append(results, consolidateChecks(path, checkFns))
	}
	return results
}

func matchingChecks(checkNames []string) []RepoCheckFunc {
	var res []RepoCheckFunc
	for _, name := range checkNames {
		for _, fn := range RepoChecks {
			if name == utils.FunctionName(fn) {
				res = append(res, fn)
			}
		}
	}
	return res
}

func consolidateChecks(path string, checkFns []RepoCheckFunc) *utils.CheckResultConsolidated {
	repo, err := mygit.CloneRepo(path)
	if err != nil {
		return &utils.CheckResultConsolidated{
			Date:       time.Now(),
			Repository: path,
			Error:      err.Error(),
		}
	}
	res := repositoryCheck(repo, checkFns)
	if err := repo.Cleanup(); err != nil {
		fmt.Fprintf(os.Stderr, "Error while cleaning up repo %s: %v", utils.RepoName(repo.Repo), err)
	}
	return res
}

func repositoryCheck(repo *mygit.ClonedRepo, checkFns []RepoCheckFunc) *utils.CheckResultConsolidated {
	res := &utils.CheckResultConsolidated{
		Date:       time.Now(),
		Repository: utils.RepoName(repo.Repo),
		Checks:     []utils.SingleCheck{},
	}
	waitGroup := new(sync.WaitGroup)
	waitGroup.Add(len(checkFns))
	for _, check := range checkFns {
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
