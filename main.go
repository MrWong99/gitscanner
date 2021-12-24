package main

import (
	"embed"
	"flag"
	"io/fs"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"strconv"

	"github.com/MrWong99/gitscanner/checks"
	"github.com/MrWong99/gitscanner/checks/binaryfile"
	"github.com/MrWong99/gitscanner/checks/commitmeta"
	"github.com/MrWong99/gitscanner/checks/filesize"
	"github.com/MrWong99/gitscanner/checks/unicode"
	"github.com/MrWong99/gitscanner/config"
	mygit "github.com/MrWong99/gitscanner/git"
	"github.com/MrWong99/gitscanner/rest"
	"github.com/gorilla/mux"
)

//go:embed ui/dist/search-binary/*
var embedUi embed.FS

//go:embed VERSION
var embedVersion string

func main() {
	configFile := flag.String("config", "GrootConfig.yml", "The absolute or relative path of the application configuration file.")
	flag.Parse()

	fileContent, err := ioutil.ReadFile(*configFile)
	if err != nil {
		log.Printf("Could not read file '%s'. Error: %v\n", *configFile, err)
		os.Exit(1)
	}
	config.ConfigLocation = *configFile
	var cfg *config.Config
	switch path.Ext(*configFile) {
	case ".yml", ".yaml":
		cfg, err = config.ReadYaml(fileContent)
	case ".json":
		cfg, err = config.ReadJson(fileContent)
	default:
		log.Printf("Could not determine file type of '%s', just trying to decode it as yaml...\n", *configFile)
		cfg, err = config.ReadYaml(fileContent)
	}
	if err != nil {
		log.Printf("Could not decode file '%s'. Error: %v\n", *configFile, err)
		os.Exit(1)
	}

	initializeChecks(cfg.Checks, []checks.Checker{
		&binaryfile.BinarySearchCheck{},
		&commitmeta.CommitMetaInfoCheck{},
		&unicode.UnicodeCharacterSearch{},
		&filesize.FilesizeSearchCheck{},
	})

	if cfg.Server == nil || cfg.Server.Port < 1 {
		log.Println("No port configured to start server.")
		os.Exit(1)
	}

	if cfg.Auth != nil {
		if cfg.Auth.BasicAuth != nil {
			mygit.InitHttpBasicAuth(cfg.Auth.BasicAuth.Username, cfg.Auth.BasicAuth.Password)
		}
		if cfg.Auth.Ssh != nil {
			_, err := os.Stat(cfg.Auth.Ssh.PrivateKeyFile)
			if err != nil {
				log.Printf("Private key file %s could not be opened: %v\n", cfg.Auth.Ssh.PrivateKeyFile, err)
			} else {
				content, err := ioutil.ReadFile(cfg.Auth.Ssh.PrivateKeyFile)
				if err != nil {
					log.Printf("Private key file %s could not be opened: %v\n", cfg.Auth.Ssh.PrivateKeyFile, err)
				} else {
					if err = mygit.InitSshKey(content, cfg.Auth.Ssh.KeyPassphrase); err != nil {
						log.Printf("Private key file %s could not be opened: %v\n", cfg.Auth.Ssh.PrivateKeyFile, err)
					}
				}
			}
		}
	}

	router := mux.NewRouter()
	rest.InitRouter(router)
	files, err := fs.Sub(embedUi, "ui/dist/search-binary")
	if err != nil {
		panic(err)
	}
	router.PathPrefix("/").Handler(http.FileServer(http.FS(files)))
	if cfg.Server.Tls != nil {
		log.Printf("Starting gitscanner %s\nNavigate to https://localhost:%d in your browser!\n", embedVersion, cfg.Server.Port)
		err = http.ListenAndServeTLS(":"+strconv.Itoa(cfg.Server.Port), cfg.Server.Tls.CertFile, cfg.Server.Tls.PrivateKeyFile, router)
	} else {
		log.Printf("Starting gitscanner %s\nNavigate to http://localhost:%d in your browser!\n", embedVersion, cfg.Server.Port)
		err = http.ListenAndServe(":"+strconv.Itoa(cfg.Server.Port), router)
	}
	if err != nil {
		log.Printf("%v\n", err)
	}
}

func initializeChecks(configuredChecks []config.CheckConfig, availableChecks []checks.Checker) {
	configMap := map[string]map[string]interface{}{}
	for _, cfgCheck := range configuredChecks {
		configMap[cfgCheck.Name] = cfgCheck.Config
	}
	for _, check := range availableChecks {
		checks.AddCheck(check)
		configurableCheck, configurable := check.(checks.ConfigurableChecker)
		if !configurable {
			continue
		}
		cfg, ok := configMap[check.String()]
		if !ok {
			configurableCheck.SetConfig(map[string]interface{}{})
			continue
		}
		err := configurableCheck.SetConfig(cfg)
		if err != nil {
			log.Printf("Error while configuring check '%s'. Error: %v", check.String(), err)
		}
	}
}
