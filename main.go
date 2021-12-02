package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/MrWong99/gitscanner/checks"
	mygit "github.com/MrWong99/gitscanner/git"
	"github.com/MrWong99/gitscanner/utils"
)

func main() {
	repositoryPaths := flag.String("repositories", "",
		"A comma-separated list of repositories to perform checks against. Can be in these formats:\n"+
			"- http(s)://<remote URL>\n- git@<remote URL>\n- file://<path>")
	username := flag.String("username", "", "An optional username for http basic auth.")
	password := flag.String("password", "", "An optional password for http basic auth.")
	privateKeyFile := flag.String("ssh-private-key-file", "", "An optional path to a SSH private key file in PEM format.")
	keyPassword := flag.String("ssh-private-key-password", "", "An optional password if the given private key file is encrypted.")
	branchPattern := flag.String("branch-pattern", "", "Optional pattern to match refs against. Only matches will be processed in checks that rely on refs.")
	namePattern := flag.String("name-pattern", "", "Pattern to match all commiter and author names against. This will be used for the commitmeta.CheckCommits check.")
	emailPattern := flag.String("email-pattern", "", "Pattern to match all commiter and author emails against. This will be used for the commitmeta.CheckCommits check.")
	flag.Parse()

	utils.InitConfig(&utils.GlobalConfig{
		BranchPattern: extractPattern(*branchPattern),
		NamePattern:   extractPattern(*namePattern),
		EmailPattern:  extractPattern(*emailPattern),
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
	res := checks.CheckAllRepositories(allPaths)
	jsonStr, err := json.Marshal(res)
	if err == nil {
		fmt.Printf("%s\n\n", jsonStr)
	} else {
		fmt.Fprintf(os.Stderr, "JSON failed '%v'\n", err)
		for _, v := range res {
			fmt.Printf("%v\n\n", *v)
		}
	}
}

func extractPattern(input string) *regexp.Regexp {
	if input == "" {
		return regexp.MustCompile(".*")
	}
	pat, err := regexp.Compile(input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Given pattern '%s' is not a valid regex: %v\n", input, err)
		os.Exit(1)
	}
	return pat
}
