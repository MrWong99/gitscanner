package main

import (
	"embed"
	"encoding/json"
	"flag"
	"io/fs"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/MrWong99/gitscanner/checks"
	"github.com/MrWong99/gitscanner/checks/binaryfile"
	"github.com/MrWong99/gitscanner/checks/commitmeta"
	"github.com/MrWong99/gitscanner/checks/filesize"
	"github.com/MrWong99/gitscanner/checks/unicode"
	"github.com/MrWong99/gitscanner/db"
	"github.com/MrWong99/gitscanner/db/configrepo"
	mygit "github.com/MrWong99/gitscanner/git"
	"github.com/MrWong99/gitscanner/rest"
	"github.com/MrWong99/gitscanner/utils"
	"github.com/gorilla/mux"
)

//go:embed ui/dist/search-binary/*
var embedUi embed.FS

//go:embed VERSION
var embedVersion string

func main() {
	log.SetOutput(os.Stderr)
	repositoryPaths := flag.String("repositories", "",
		"A comma-separated list of repositories to perform checks against. Can be in these formats:\n"+
			"- http(s)://<remote URL>\n- git@<remote URL>\n- file://<path>")
	username := flag.String("username", "", "An optional username for http basic auth.")
	password := flag.String("password", "", "An optional password for http basic auth.")
	privateKeyFile := flag.String("ssh-private-key-file", "", "An optional path to a SSH private key file in PEM format.")
	keyPassword := flag.String("ssh-private-key-password", "", "An optional password if the given private key file is encrypted.")
	branchPattern := flag.String("branch-pattern", ".*", "Optional pattern to match refs against. Only matches will be processed in checks that rely on refs.")
	namePattern := flag.String("name-pattern", ".*", "Pattern to match all commiter and author names against. This will be used for the commitmeta.CheckCommits check.")
	emailPattern := flag.String("email-pattern", ".*", "Pattern to match all commiter and author emails against. This will be used for the commitmeta.CheckCommits check.")
	filesizeThreshold := flag.Int64("filesize-threshold-bytes", 81920, "Amout of bytes that a file should have maximum to trigger this check.")
	port := flag.Int("port", -1, "When provided this will startup a webserver including ui that can be used to perform the checks via browser.")
	sslKeyFile := flag.String("ssl-private-key-file", "", "An optional path to a TLS private key file in PEM format to enable HTTPS. Only used when port is set.")
	sslCertFile := flag.String("ssl-certificate-chain-file", "", "An optional path to a TLS certificate (chain) in PEM format to enable HTTPS. Only used when port is set.")
	flag.Parse()

	checks.AddCheck(&binaryfile.BinarySearchCheck{})
	checks.AddCheck(&commitmeta.CommitMetaInfoCheck{})
	checks.AddCheck(&unicode.UnicodeCharacterSearch{})
	checks.AddCheck(&filesize.FilesizeSearchCheck{})

	saveConfigToDB(branchPattern, namePattern, emailPattern, filesizeThreshold)

	if *repositoryPaths == "" && *port < 1 {
		log.Println("No repositories defined!")
		os.Exit(1)
	}

	if *username != "" {
		mygit.InitHttpBasicAuth(*username, *password)
	}
	if *privateKeyFile != "" {
		_, err := os.Stat(*privateKeyFile)
		if err != nil {
			log.Printf("Private key file %s could not be opened: %v\n", *privateKeyFile, err)
		} else {
			content, err := ioutil.ReadFile(*privateKeyFile)
			if err != nil {
				log.Printf("Private key file %s could not be opened: %v\n", *privateKeyFile, err)
			} else {
				if err = mygit.InitSshKey(content, *keyPassword); err != nil {
					log.Printf("Private key file %s could not be opened: %v\n", *privateKeyFile, err)
				}
			}
		}
	}

	if *port > 0 {
		log.SetOutput(os.Stdout)
		router := mux.NewRouter()
		rest.InitRouter(router)
		files, err := fs.Sub(embedUi, "ui/dist/search-binary")
		if err != nil {
			panic(err)
		}
		router.PathPrefix("/").Handler(http.FileServer(http.FS(files)))
		if *sslKeyFile != "" && *sslCertFile != "" {
			log.Printf("Starting gitscanner %s\nNavigate to https://localhost:%d in your browser!\n", embedVersion, *port)
			err = http.ListenAndServeTLS(":"+strconv.Itoa(*port), *sslCertFile, *sslKeyFile, router)
		} else {
			log.Printf("Starting gitscanner %s\nNavigate to http://localhost:%d in your browser!\n", embedVersion, *port)
			err = http.ListenAndServe(":"+strconv.Itoa(*port), router)
		}
		if err != nil {
			log.Printf("%v\n", err)
		}
	} else {
		allPaths := strings.Split(*repositoryPaths, ",")
		res := checks.CheckAllRepositories(allPaths)
		jsonStr, err := json.Marshal(res)
		if err == nil {
			log.Printf("%s\n\n", jsonStr)
		} else {
			log.Printf("JSON failed '%v'\n", err)
			for _, v := range res {
				log.Printf("%v\n\n", *v)
			}
		}
	}
}

func saveConfigToDB(branchPattern, namePattern, emailPattern *string, filesizeThreshold *int64) {
	var err error

	err = db.InitDb()
	if err != nil {
		log.Printf("Error while initializing database: %v\n", err)
		os.Exit(1)
	}

	if _, err = utils.ExtractPattern(*branchPattern); err != nil {
		log.Printf("Error with given pattern %s: %v\n", *branchPattern, err)
		os.Exit(1)
	}
	if _, err = utils.ExtractPattern(*namePattern); err != nil {
		log.Printf("Error with given pattern %s: %v\n", *namePattern, err)
		os.Exit(1)
	}
	if _, err = utils.ExtractPattern(*emailPattern); err != nil {
		log.Printf("Error with given pattern %s: %v\n", *emailPattern, err)
		os.Exit(1)
	}

	for _, check := range checks.RepoChecks {
		configurableCheck, configurable := check.(checks.ConfigurableChecker)
		if !configurable {
			continue
		}

		var currentCfg *checks.CheckConfiguration

		currentCfg, err = configrepo.ReadConfig(configurableCheck.String())
		if err != nil {
			log.Printf("Error while reading config '%s' from database: %v\n", check.String(), err)
		}
		if currentCfg != nil {
			configurableCheck.SetConfig(currentCfg)
			continue
		}
		newCfg := &checks.CheckConfiguration{
			CheckName: configurableCheck.String(),
		}
		var err error
		switch configurableCheck.String() {
		case "SearchBinaries":
			err = newCfg.SetConfigMap(map[string]interface{}{
				"branchPattern": *branchPattern,
			})
		case "CheckCommitMetaInformation":
			err = newCfg.SetConfigMap(map[string]interface{}{
				"emailPattern": *emailPattern,
				"namePattern":  *namePattern,
			})
		case "SearchBigFiles":
			err = newCfg.SetConfigMap(map[string]interface{}{
				"branchPattern":         *branchPattern,
				"filesizeThresholdByte": *filesizeThreshold,
			})
		case "SearchIllegalUnicodeCharacters":
			err = newCfg.SetConfigMap(map[string]interface{}{
				"branchPattern": *branchPattern,
			})
		default:
			log.Printf("Error while updating config for check '%s': Name not registered in saveConfigToDB.\n", configurableCheck.String())
		}
		if err != nil {
			log.Printf("Error while updating config for check '%s': %v\n", configurableCheck.String(), err)
			continue
		}
		configurableCheck.SetConfig(newCfg)
		err = configrepo.UpdateConfig(newCfg)
		if err != nil {
			log.Printf("Error while updating config in database: %v\n", err)
			os.Exit(1)
		}
	}
}
