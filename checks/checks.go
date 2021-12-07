package checks

import (
	"encoding/json"
	"errors"
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
	CheckName string         `json:"checkName" gorm:"primaryKey"`
	Config    datatypes.JSON `json:"config"`
	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

func (config *CheckConfiguration) SetConfigMap(cMap map[string]interface{}) error {
	cfg, err := json.Marshal(&cMap)
	if err != nil {
		return err
	}
	config.Config = cfg
	return nil
}

func (config *CheckConfiguration) ParseConfigMap() (map[string]interface{}, error) {
	var result map[string]interface{}
	err := json.Unmarshal(config.Config, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (config *CheckConfiguration) MustParseConfigMap() map[string]interface{} {
	var result map[string]interface{}
	err := json.Unmarshal(config.Config, &result)
	if err != nil {
		panic(err)
	}
	return result
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

func Configure(config *CheckConfiguration) error {
	for _, check := range RepoChecks {
		if check.String() == config.CheckName {
			switch ctype := check.(type) {
			case ConfigurableChecker:
				return ctype.SetConfig(config)
			default:
				return errors.New("Check '" + config.CheckName + "' is not configurable.")
			}
		}
	}
	return errors.New("A check with name '" + config.CheckName + "' is not registered.")
}

func GetCurrentConfig(checkname string) (*CheckConfiguration, error) {
	for _, check := range RepoChecks {
		if check.String() == checkname {
			switch ctype := check.(type) {
			case ConfigurableChecker:
				return ctype.GetConfig(), nil
			default:
				return nil, errors.New("Check '" + checkname + "' is not configurable.")
			}
		}
	}
	return nil, errors.New("A check with name '" + checkname + "' is not registered.")
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
