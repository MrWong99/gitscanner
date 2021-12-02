package main

import (
	"flag"
	"fmt"
	"os"
	"regexp"
	"strings"
	"sync"

	"github.com/MrWong99/gitscanner/checks/binaryfile"
	"github.com/MrWong99/gitscanner/checks/unicode"
	mygit "github.com/MrWong99/gitscanner/git"
	"github.com/MrWong99/gitscanner/utils"
)

type RepoCheckFunc func(*mygit.ClonedRepo) error

var repoChecks = []RepoCheckFunc{
	binaryfile.SearchBinaries,
	unicode.SearchUnicode,
}

var waitGroup *sync.WaitGroup

func main() {
	repositoryPaths := flag.String("repositories", "",
		"A comma-separated list of repositories to perform checks against. Can be in these formats:\n"+
			"- http(s)://<remote URL>\n- git@<remote URL>\n- file://<path>")
	username := flag.String("username", "", "An optional username for http basic auth.")
	password := flag.String("password", "", "An optional password for http basic auth.")
	privateKeyFile := flag.String("ssh-private-key-file", "", "An optional path to a SSH private key file in PEM format.")
	keyPassword := flag.String("ssh-private-key-password", "", "An optional password if the given private key file is encrypted.")
	branchPattern := flag.String("branch-pattern", "", "Optional pattern to match refs against. Only matches will be processed in checks that rely on refs.")
	flag.Parse()

	var pat *regexp.Regexp
	if *branchPattern == "" {
		pat = regexp.MustCompile(".*")
	} else {
		var err error
		pat, err = regexp.Compile(*branchPattern)
		if err != nil {
			fmt.Printf("Given branch pattern '%s' is not a valid regex: %v", *branchPattern, err)
			os.Exit(1)
		}
	}
	utils.InitConfig(&utils.GlobalConfig{
		BranchPattern: pat,
	})

	if *repositoryPaths == "" {
		fmt.Println("No repositories defined!")
		os.Exit(1)
	}

	if *username != "" {
		mygit.InitHttpBasicAuth(*username, *password)
	}
	if *privateKeyFile != "" {
		mygit.InitSshKey(*privateKeyFile, *keyPassword)
	}

	allPaths := strings.Split(*repositoryPaths, ",")
	waitGroup = new(sync.WaitGroup)
	for _, path := range allPaths {
		repo, err := mygit.CloneRepo(path)
		if err != nil {
			fmt.Printf("Error while opening repo '%s': %v\n", path, err)
			return
		}
		repositoryCheck(repo)
		if err := repo.Cleanup(); err != nil {
			fmt.Printf("Error while cleaning up repo %s: %v", utils.RepoName(repo.Repo), err)
		}
	}

	waitGroup.Wait()
}

func repositoryCheck(repo *mygit.ClonedRepo) error {
	waitGroup.Add(len(repoChecks))
	for _, check := range repoChecks {
		go func(r *mygit.ClonedRepo, checkFn RepoCheckFunc) {
			defer waitGroup.Done()
			fmt.Printf("Checking repository %v with function %s\n", utils.RepoName(r.Repo), utils.FunctionName(checkFn))
			err := checkFn(repo)
			if err != nil {
				fmt.Printf("Error while checking repository %s with function %s: %v\n", utils.RepoName(r.Repo), utils.FunctionName(checkFn), err)
			}
		}(repo, check)
	}
	return nil
}
