package main

import (
	"embed"
	"flag"
	"fmt"
	"io/fs"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"strconv"
	"syscall"

	"github.com/MrWong99/gitscanner/checks"
	"github.com/MrWong99/gitscanner/checks/binaryfile"
	"github.com/MrWong99/gitscanner/checks/commitmeta"
	"github.com/MrWong99/gitscanner/checks/filesize"
	"github.com/MrWong99/gitscanner/checks/unicode"
	"github.com/MrWong99/gitscanner/config"
	"github.com/MrWong99/gitscanner/config/encryption"
	mygit "github.com/MrWong99/gitscanner/git"
	"github.com/MrWong99/gitscanner/rest"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/ssh/terminal"
)

//go:embed ui/dist/search-binary/*
var embedUi embed.FS

//go:embed VERSION
var embedVersion string

func main() {
	configFile := flag.String("config", "GrootConfig.yml", "The absolute or relative path of the application configuration file.")
	encryptionKey := flag.String("encryptionKey", "",
		"Key to use for en-/decrypting sensitive data. Can also be provided via environment variable 'ENCRYPTION_KEY' or by typing into console after start.")
	encrypt := flag.String("encrypt", "",
		"When set this tool will simply encrypt the given input and exit afterwards. Can be used to encrypt any value for the config file.")
	decrypt := flag.String("decrypt", "",
		"When set this tool will simply decrypt the given input and exit afterwards. Can be used to decrypt any value for the config file given the correct key.")
	flag.Parse()

	if *encryptionKey == "" {
		key, found := os.LookupEnv("ENCRYPTION_KEY")
		if found {
			encryptionKey = &key
		} else {
			encryptionKey = readKeyFromUserInput()
		}
	}
	encryption.SetEncryptionKey(*encryptionKey)

	if *encrypt != "" {
		res, err := encryption.EncryptConfigString(*encrypt)
		if err != nil {
			log.Printf("Error during encryption: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Encrypted input is:\n%s\n", res)
		os.Exit(0)
	}

	if *decrypt != "" {
		res, err := encryption.DecryptConfigString(*decrypt)
		if err != nil {
			log.Printf("Error during decryption: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Decrypted input is:\n%s\n", res)
		os.Exit(0)
	}

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
			name, err := encryption.DecryptConfigString(cfg.Auth.BasicAuth.Username)
			if err != nil {
				log.Printf("Could not decrypt auth.basicAuth.username, Error: %v\n", err)
				os.Exit(1)
			}
			password, err := encryption.DecryptConfigString(cfg.Auth.BasicAuth.Password)
			if err != nil {
				log.Printf("Could not decrypt auth.basicAuth.password, Error: %v\n", err)
				os.Exit(1)
			}
			mygit.InitHttpBasicAuth(name, password)
		}
		if cfg.Auth.Ssh != nil {
			sshFile, err := encryption.DecryptConfigString(cfg.Auth.Ssh.PrivateKeyFile)
			if err != nil {
				log.Printf("Could not decrypt auth.ssh.privateKeyFile, Error: %v\n", err)
				os.Exit(1)
			}
			passphrase, err := encryption.DecryptConfigString(cfg.Auth.Ssh.KeyPassphrase)
			if err != nil {
				log.Printf("Could not decrypt auth.ssh.keyPassphrase, Error: %v\n", err)
				os.Exit(1)
			}
			_, err = os.Stat(sshFile)
			if err != nil {
				log.Printf("Private key file %s could not be opened: %v\n", sshFile, err)
			} else {
				content, err := ioutil.ReadFile(sshFile)
				if err != nil {
					log.Printf("Private key file %s could not be opened: %v\n", sshFile, err)
				} else {
					if err = mygit.InitSshKey(content, passphrase); err != nil {
						log.Printf("Private key file %s could not be opened: %v\n", sshFile, err)
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
		privateKeyFile, err := encryption.DecryptConfigString(cfg.Server.Tls.PrivateKeyFile)
		if err != nil {
			log.Printf("Could not decrypt server.tls.privateKeyFile, Error: %v\n", err)
			os.Exit(1)
		}
		certFile, err := encryption.DecryptConfigString(cfg.Server.Tls.CertFile)
		if err != nil {
			log.Printf("Could not decrypt server.tls.certFile, Error: %v\n", err)
			os.Exit(1)
		}
		log.Printf("Starting gitscanner %s\nNavigate to https://localhost:%d in your browser!\n", embedVersion, cfg.Server.Port)
		err = http.ListenAndServeTLS(":"+strconv.Itoa(cfg.Server.Port), privateKeyFile, certFile, router)
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

func readKeyFromUserInput() *string {
	fmt.Println("No encryption key provided to store configuration securly. Type it in now:")
	password, err := terminal.ReadPassword(int(syscall.Stdin))
	if err != nil {
		log.Printf("Error while reading password: %v\n", err)
		os.Exit(1)
	}
	pw := string(password)
	return &pw
}
