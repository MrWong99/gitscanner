package rest

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/MrWong99/gitscanner/checks"
	"github.com/MrWong99/gitscanner/db/checkrepo"
	"github.com/MrWong99/gitscanner/db/configrepo"
	mygit "github.com/MrWong99/gitscanner/git"
	"github.com/MrWong99/gitscanner/utils"
	"github.com/gorilla/mux"
)

func InitRouter(router *mux.Router) {
	router.HandleFunc("/api/v1/checkRepos", handleCheckRequest).Methods("POST")
	router.HandleFunc("/api/v1/config", handleGetConfig).Methods("GET")
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
	if err != nil {
		log.Println(err)
		w.WriteHeader(400)
		w.Write([]byte(err.Error()))
		return
	}
	var request utils.SearchRequestBody
	err = json.Unmarshal(body, &request)
	if err != nil {
		log.Println(err)
		w.WriteHeader(400)
		w.Write([]byte(err.Error()))
		return
	}
	checks := checks.CheckAllRepositoriesSpecificChecks(strings.Split(request.Path, ","), request.CheckNames)
	err = checkrepo.SaveChecks(checks)
	if err != nil {
		log.Println(err)
		w.WriteHeader(400)
		w.Write([]byte(err.Error()))
		return
	}
	json.NewEncoder(w).Encode(checks)
}

// GET /api/v1/config
func handleGetConfig(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	cfg, err := configrepo.ReadConfig()
	if err != nil {
		log.Println(err)
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		return
	}
	json.NewEncoder(w).Encode(cfg)
}

// PUT /api/v1/config
func handlePutConfig(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		w.WriteHeader(400)
		return
	}
	var request utils.GlobalConfig
	err = json.Unmarshal(body, &request)
	if err != nil {
		log.Println(err)
		w.WriteHeader(400)
		return
	}
	err = updateConfig(&request)
	if err != nil {
		log.Println(err)
		w.WriteHeader(400)
		w.Write([]byte(err.Error()))
		return
	} else {
		w.WriteHeader(200)
	}
}

// PUT /api/v1/config/sshkey
func handlePutSshKey(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		w.WriteHeader(400)
		return
	}
	var request SshPrivateKeyInfo
	err = json.Unmarshal(body, &request)
	if err != nil {
		log.Println(err)
		w.WriteHeader(400)
		return
	}
	err = updateSshKey(&request)
	if err != nil {
		log.Println(err)
		w.WriteHeader(400)
		w.Write([]byte(err.Error()))
		return
	} else {
		w.WriteHeader(200)
	}
}

// PUT /api/v1/config/basicauth
func handlePutBasicAuth(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		w.WriteHeader(400)
		return
	}
	var request BasicAuth
	err = json.Unmarshal(body, &request)
	if err != nil {
		log.Println(err)
		w.WriteHeader(400)
		return
	}
	err = updateBasicAuth(&request)
	if err != nil {
		log.Println(err)
		w.WriteHeader(400)
		w.Write([]byte(err.Error()))
		return
	} else {
		w.WriteHeader(200)
	}
}

// GET /api/v1/checks
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
	if err != nil {
		log.Println(err)
		w.WriteHeader(400)
		w.Write([]byte(err.Error()))
	}
	err = json.Unmarshal(body, &requestAck)
	if err != nil {
		log.Println(err)
		w.WriteHeader(400)
		w.Write([]byte(err.Error()))
		return
	}
	id, err := strconv.ParseUint(vars["singleCheckId"], 10, 64)
	if err != nil {
		log.Println(err)
		w.WriteHeader(400)
		w.Write([]byte(err.Error()))
		return
	}
	if err = checkrepo.AcknowledgeCheck(uint(id), requestAck.Acknowledged); err != nil {
		log.Println(err)
		w.WriteHeader(400)
		w.Write([]byte(err.Error()))
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
	if err != nil {
		log.Println(err)
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		return
	}
	json.NewEncoder(w).Encode(checks)
}

func updateConfig(cfg *utils.GlobalConfig) error {
	var err error
	if _, err = utils.ExtractPattern(cfg.BranchPattern); err != nil {
		return err
	}
	if _, err = utils.ExtractPattern(cfg.NamePattern); err != nil {
		return err
	}
	if _, err = utils.ExtractPattern(cfg.EmailPattern); err != nil {
		return err
	}
	err = configrepo.UpdateConfig(cfg)
	return err
}

func updateSshKey(keyInfo *SshPrivateKeyInfo) error {
	return mygit.InitSshKey([]byte(keyInfo.Key), keyInfo.Password)
}

func updateBasicAuth(auth *BasicAuth) error {
	return mygit.InitHttpBasicAuth(auth.Username, auth.Password)
}
