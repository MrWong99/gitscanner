package checks

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	mygit "github.com/MrWong99/gitscanner/git"
	"github.com/MrWong99/gitscanner/utils"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type CheckConfiguration struct {
	ID        uint           `json:"-" gorm:"primarykey"`
	Config    datatypes.JSON `json:"config"`
	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

func (config *CheckConfiguration) GetConfig() map[string]interface{} {
	var result map[string]interface{}
	json.Unmarshal(config.Config, &result)
	return result
}

func (config *CheckConfiguration) SetConfig(c map[string]interface{}) error {
	res, err := json.Marshal(c)
	if err != nil {
		return err
	}
	config.Config = res
	return nil
}

type Checker interface {
	fmt.Stringer
	Check(*mygit.ClonedRepo, chan<- utils.SingleCheck) error
}

type Configurer interface {
	GetConfig() *CheckConfiguration
	SetConfig(*CheckConfiguration) error
}

type ConfigurableChecker interface {
	Checker
	Configurer
}

var RepoChecks []Checker

func AddCheck(check Checker) {
	RepoChecks = append(RepoChecks, check)
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

func matchingChecks(checkNames []string) []Checker {
	var res []Checker
	for _, name := range checkNames {
		for _, check := range RepoChecks {
			if name == check.String() {
				res = append(res, check)
			}
		}
	}
	return res
}

func consolidateChecks(path string, checkFns []Checker) *utils.CheckResultConsolidated {
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
		log.Printf("Error while cleaning up repo %s: %v", utils.RepoName(repo.Repo), err)
	}
	return res
}

func repositoryCheck(repo *mygit.ClonedRepo, checkFns []Checker) *utils.CheckResultConsolidated {
	res := &utils.CheckResultConsolidated{
		Date:       time.Now(),
		Repository: utils.RepoName(repo.Repo),
		Checks:     []utils.SingleCheck{},
	}
	waitGroup := new(sync.WaitGroup)
	waitGroup.Add(len(checkFns))
	for _, check := range checkFns {
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
