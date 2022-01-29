package rest

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
	"fmt"

	"github.com/MrWong99/gitscanner/checks"
	"github.com/MrWong99/gitscanner/config"
	"github.com/MrWong99/gitscanner/config/encryption"
	"github.com/MrWong99/gitscanner/db/checkrepo"
	mygit "github.com/MrWong99/gitscanner/git"
	"github.com/MrWong99/gitscanner/utils"
	"github.com/gorilla/mux"
)

func InitRouter(router *mux.Router) {
	router.HandleFunc("/api/v1/checkRepos", handleCheckRequest).Methods("POST")
	router.HandleFunc("/api/v1/staticCodeAnalysis", handleStaticCodeAnalysisRequest).Methods("POST")
	router.HandleFunc("/api/v1/config/{checkName}", handleGetConfig).Methods("GET")
	router.HandleFunc("/api/v1/config", handlePutConfig).Methods("PUT")
	router.HandleFunc("/api/v1/config/sshkey", handlePutSshKey).Methods("PUT")
	router.HandleFunc("/api/v1/config/basicauth", handlePutBasicAuth).Methods("PUT")
	router.HandleFunc("/api/v1/checkDefinitions", handleGetChecks).Methods("GET")
	router.HandleFunc("/api/v1/acknowledged/{singleCheckId}", handleAcknowledge).Methods("PUT")
	router.HandleFunc("/api/v1/checks", handleRetrieveChecksByDate).Methods("GET").Queries(
		"from", "{from:[0-9]+}",
		"to", "{to:[0-9]+}",
		"checkNames", "{checkNames:.+}",
	)
}

type SshPrivateKeyInfo struct {
	Key      string `json:"key"`
	Password string `json:"password"`
}

type BasicAuth struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Acknowledge struct {
	Acknowledged bool `json:"acknowledged"`
}

// POST /api/v1/checkRepos
func handleCheckRequest(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	body, err := ioutil.ReadAll(r.Body)
	if handleError(err, 400, w, r) {
		return
	}
	var request utils.SearchRequestBody
	err = json.Unmarshal(body, &request)
	if handleError(err, 400, w, r) {
		return
	}
	checks := checks.CheckAllRepositoriesSpecificChecks(strings.Split(request.Path, ","), request.CheckNames)
	err = checkrepo.SaveChecks(checks)
	if handleError(err, 400, w, r) {
		return
	}
	err = json.NewEncoder(w).Encode(checks)
	if err != nil {
		log.Printf("Error encoding response: %v\n", err)
	}
}

// POST /api/v1/staticCodeAnalysis
func handleStaticCodeAnalysisRequest(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	body, err := ioutil.ReadAll(r.Body)
	if handleError(err, 400, w, r) {
		return
	}
	var request utils.SearchRequestBody
	err = json.Unmarshal(body, &request)
	if handleError(err, 400, w, r) {
		return
	}
	repos := strings.Split(request.Path, ",")
	fmt.Println("Repos are: %v", repos)
	for _, repo := range repos{
		clonedRepo := mygit.ClonedRepo(repo)
		fmt.Println("Scanning repo %v %v", repo, clonedRepo.LocalDir)
		out, err := sast.semgrepScan(request.ConfigFiles, clonedRepo.LocalDir)	
	}
}

// GET /api/v1/config/{checkName}
func handleGetConfig(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	checkName := vars["checkName"]
	w.Header().Add("Content-Type", "application/json")
	cfg, err := checks.GetCurrentConfig(checkName)
	if handleError(err, 400, w, r) {
		return
	}
	json.NewEncoder(w).Encode(cfg)
}

// PUT /api/v1/config
func handlePutConfig(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	body, err := ioutil.ReadAll(r.Body)
	if handleError(err, 400, w, r) {
		return
	}
	var request config.CheckConfig
	err = json.Unmarshal(body, &request)
	if handleError(err, 400, w, r) {
		return
	}
	err = checks.Configure(&request)
	if handleError(err, 400, w, r) {
		return
	}
	cfg := config.CurrentConfig()
	checkCfg, _ := checks.GetCurrentConfig(request.Name)
	checkCfg.Enabled = request.Enabled
	cfg.AddOrUpdateCheckConfig(checkCfg)
	err = config.UpdateConfigFile()
	if handleError(err, 400, w, r) {
		return
	}
	w.WriteHeader(200)
}

// PUT /api/v1/config/sshkey
func handlePutSshKey(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	body, err := ioutil.ReadAll(r.Body)
	if handleError(err, 400, w, r) {
		return
	}
	var request SshPrivateKeyInfo
	err = json.Unmarshal(body, &request)
	if handleError(err, 400, w, r) {
		return
	}
	err = updateSshKey(&request)
	if handleError(err, 400, w, r) {
		return
	}
	w.WriteHeader(200)
}

// PUT /api/v1/config/basicauth
func handlePutBasicAuth(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	body, err := ioutil.ReadAll(r.Body)
	if handleError(err, 400, w, r) {
		return
	}
	var request BasicAuth
	err = json.Unmarshal(body, &request)
	if handleError(err, 400, w, r) {
		return
	}
	err = updateBasicAuth(&request)
	if handleError(err, 400, w, r) {
		return
	}
	w.WriteHeader(200)
}

// GET /api/v1/checkDefinitions
func handleGetChecks(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	var checkNames []string
	for _, v := range checks.RepoChecks {
		checkNames = append(checkNames, v.String())
	}
	json.NewEncoder(w).Encode(checkNames)
}

// PUT /api/v1/acknowledged/{singleCheckId}
func handleAcknowledge(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	vars := mux.Vars(r)
	var requestAck Acknowledge
	body, err := ioutil.ReadAll(r.Body)
	if handleError(err, 400, w, r) {
		return
	}
	err = json.Unmarshal(body, &requestAck)
	if handleError(err, 400, w, r) {
		return
	}
	id, err := strconv.ParseUint(vars["singleCheckId"], 10, 64)
	if handleError(err, 400, w, r) {
		return
	}
	err = checkrepo.AcknowledgeCheck(uint(id), requestAck.Acknowledged)
	if handleError(err, 400, w, r) {
		return
	}
	w.WriteHeader(200)
}

// GET /api/v1/checks?from=4358905432&to=53987518&checkNames=SearchBinaries,CheckCommitMetaInformation
func handleRetrieveChecksByDate(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	vars := mux.Vars(r)
	from, _ := strconv.ParseInt(vars["from"], 10, 64)
	to, _ := strconv.ParseInt(vars["to"], 10, 64)
	checks, err := checkrepo.ReadSavedChecks(strings.Split(vars["checkNames"], ","), time.UnixMilli(from), time.UnixMilli(to))
	if handleError(err, 500, w, r) {
		return
	}
	json.NewEncoder(w).Encode(checks)
}

func handleError(err error, statusCode int, w http.ResponseWriter, r *http.Request) bool {
	if err == nil {
		return false
	}
	log.Printf("%s on %s failed with error: %v\n", r.Method, r.URL, err)
	log.Println(err)
	w.WriteHeader(statusCode)
	w.Write([]byte(err.Error()))
	return true
}

func updateSshKey(keyInfo *SshPrivateKeyInfo) error {
	privateKeyFileName := "privatekey.pem"
	err := mygit.InitSshKey([]byte(keyInfo.Key), keyInfo.Password)
	if err != nil {
		return err
	}
	password, err := encryption.EncryptConfigString(keyInfo.Password)
	if err != nil {
		return err
	}
	cfg := config.CurrentConfig()
	if cfg.Auth == nil {
		cfg.Auth = &config.AuthConfig{
			Ssh: &config.SshConfig{
				PrivateKeyFile: "./"+privateKeyFileName,
				KeyPassphrase:  password,
			},
		}
	} else {
		cfg.Auth.Ssh = &config.SshConfig{
			PrivateKeyFile: "./"+privateKeyFileName,
			KeyPassphrase:  password,
		}
	}

	keyHandle, err := os.OpenFile(privateKeyFileName, os.O_RDWR|os.O_CREATE, 0600)
	if err != nil{
		log.Printf("Could not open private key file %v", err)
	}
	defer keyHandle.Close()
	_, err = keyHandle.WriteString(keyInfo.Key)
	if err != nil{
		log.Printf("Could not write bytes in private key file %v", err)
	}
	
	return config.UpdateConfigFile()
}

func updateBasicAuth(auth *BasicAuth) error {
	err := mygit.InitHttpBasicAuth(auth.Username, auth.Password)
	if err != nil {
		return err
	}
	username, err := encryption.EncryptConfigString(auth.Username)
	if err != nil {
		return err
	}
	password, err := encryption.EncryptConfigString(auth.Password)
	if err != nil {
		return err
	}
	cfg := config.CurrentConfig()
	if cfg.Auth == nil {
		cfg.Auth = &config.AuthConfig{
			BasicAuth: &config.BasicAuthConfig{
				Username: username,
				Password: password,
			},
		}
	} else {
		cfg.Auth.BasicAuth = &config.BasicAuthConfig{
			Username: username,
			Password: password,
		}
	}
	return config.UpdateConfigFile()
}
