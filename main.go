package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/MrWong99/gitscanner/checks"
	mygit "github.com/MrWong99/gitscanner/git"
	"github.com/MrWong99/gitscanner/utils"
	"github.com/gorilla/mux"
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
	port := flag.Int("port", -1, "When provided this will startup a webserver including ui that can be used to perform the checks via browser.")
	sslKeyFile := flag.String("ssl-private-key-file", "", "An optional path to a TLS private key file in PEM format to enable HTTPS. Only used when port is set.")
	sslCertFile := flag.String("ssl-certificate-chain-file", "", "An optional path to a TLS certificate (chain) in PEM format to enable HTTPS. Only used when port is set.")
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

	if *port > 0 {
		router := mux.NewRouter().StrictSlash(true)
		router.HandleFunc("/api/v1/checkRepos", handleCheckRequest).Methods("POST")
		var err error
		if *sslKeyFile != "" && *sslCertFile != "" {
			err = http.ListenAndServeTLS(":"+strconv.Itoa(*port), *sslCertFile, *sslKeyFile, router)
		} else {
			err = http.ListenAndServe(":"+strconv.Itoa(*port), router)
		}
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
		}
	} else {
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
}

func handleCheckRequest(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(400)
	}
	var request utils.SearchRequestBody
	err = json.Unmarshal(body, &request)
	if err != nil {
		w.WriteHeader(400)
	}

	json.NewEncoder(w).Encode(checks.CheckAllRepositoriesSpecificChecks(strings.Split(request.Path, ","), request.CheckNames))
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
