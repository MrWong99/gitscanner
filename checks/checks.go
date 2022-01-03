package checks

import (
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/MrWong99/gitscanner/config"
	mygit "github.com/MrWong99/gitscanner/git"
	"github.com/MrWong99/gitscanner/utils"
)

type Checker interface {
	fmt.Stringer
	Check(*mygit.ClonedRepo, chan<- utils.SingleCheck) error
}

type Configurer interface {
	GetConfig() map[string]interface{}
	SetConfig(map[string]interface{}) error
}

type ConfigurableChecker interface {
	Checker
	Configurer
}

var RepoChecks map[string]Checker = map[string]Checker{}

func AddCheck(check Checker) {
	RepoChecks[check.String()] = check
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

func Configure(config *config.CheckConfig) error {
	check, ok := RepoChecks[config.Name]
	if !ok {
		return errors.New("a check with name '" + config.Name + "' does not exist or is disabled")
	}
	configurableCheck, configurable := check.(ConfigurableChecker)
	if !configurable {
		return errors.New("the check '" + config.Name + "' does not need any configuration")
	}
	return configurableCheck.SetConfig(config.Config)
}

func GetCurrentConfig(name string) (*config.CheckConfig, error) {
	check, ok := RepoChecks[name]
	if !ok {
		return nil, errors.New("a check with name '" + name + "' does not exist or is disabled")
	}
	configurableCheck, configurable := check.(ConfigurableChecker)
	if !configurable {
		return nil, errors.New("the check '" + name + "' does not need any configuration")
	}
	return &config.CheckConfig{
		Name:    check.String(),
		Enabled: config.CurrentConfig().CheckIsEnabled(check.String()),
		Config:  configurableCheck.GetConfig(),
	}, nil
}

func matchingChecks(checkNames []string) map[string]Checker {
	res := map[string]Checker{}
	for _, name := range checkNames {
		check, ok := RepoChecks[name]
		if !ok {
			log.Printf("Requested to check using '%s' but that check is not known.\n", name)
			continue
		}
		res[name] = check
	}
	return res
}

func consolidateChecks(path string, checkFns map[string]Checker) *utils.CheckResultConsolidated {
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
		log.Printf("Error while cleaning up repo %s: %v\n", utils.RepoName(repo.Repo), err)
	}
	return res
}

func repositoryCheck(repo *mygit.ClonedRepo, checkFns map[string]Checker) *utils.CheckResultConsolidated {
	res := &utils.CheckResultConsolidated{
		Date:       time.Now(),
		Repository: utils.RepoName(repo.Repo),
		Checks:     []utils.SingleCheck{},
	}
	waitGroup := new(sync.WaitGroup)
	waitGroup.Add(len(checkFns))
	for name, check := range checkFns {
		if !config.CurrentConfig().CheckIsEnabled(name) {
			waitGroup.Done()
			continue
		}
		checkChan := make(chan utils.SingleCheck)
		go func(r *mygit.ClonedRepo, checker Checker, outputs chan<- utils.SingleCheck) {
			err := checker.Check(repo, checkChan)
			if err != nil {
				log.Printf("Error while checking repository %s with function %s: %v\n", utils.RepoName(r.Repo), checker.String(), err)
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
