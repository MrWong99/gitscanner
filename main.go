package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"sync"

	mygit "github.com/MrWong99/gitscanner/git"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

type RepoCheckFunc func(*git.Repository) error

var repoChecks = []RepoCheckFunc{
	searchBinaries,
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
	flag.Parse()

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
	waitGroup.Add(len(allPaths))
	for _, path := range allPaths {
		go func(p string) {
			defer waitGroup.Done()
			repo, err := mygit.LoadRepoInMemory(p)
			if err != nil {
				fmt.Printf("Error while opening repo '%s': %v\n", p, err)
				return
			}
			repositoryCheck(repo)
		}(path)
	}

	waitGroup.Wait()
}

func repositoryCheck(repo *git.Repository) error {
	fmt.Printf("Checking repository %v\n", repo)
	for _, check := range repoChecks {
		fmt.Printf("Checking repository %v with function %v\n", repo, check)
		go func(r *git.Repository, checkFn RepoCheckFunc) {
			fmt.Printf("Checking repository %v with function %v\n", r, checkFn)
			err := checkFn(repo)
			if err != nil {
				fmt.Printf("Error while checking repository %v with function %v: %v\n", r, checkFn, err)
			}
		}(repo, check)
	}
	return nil
}

func searchBinaries(repo *git.Repository) error {
	branchIt, err := repo.Branches()
	if err != nil {
		return err
	}
	worktree, err := repo.Worktree()
	if err != nil {
		return err
	}
	branchIt.ForEach(func(r *plumbing.Reference) error {
		fmt.Printf("Processing branch %s with Root %s\n", r.Name(), worktree.Filesystem.Root())
		err = worktree.Checkout(&git.CheckoutOptions{
			Branch: r.Name(),
			Force:  true,
		})
		if err != nil {
			return nil
		}
		fmt.Printf("Processing branch %s with Root %s\n", r.Name(), worktree.Filesystem.Root())
		return nil
	})
	return nil
}
