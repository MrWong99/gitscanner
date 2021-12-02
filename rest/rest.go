package rest

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"

	"github.com/MrWong99/gitscanner/checks"
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
	router.HandleFunc("/api/v1/checks", handleGetChecks).Methods("GET")
}

type ConfigDto struct {
	BranchPattern string `json:"branchPattern"`
	NamePattern   string `json:"namePattern"`
	EmailPattern  string `json:"emailPattern"`
}

type SshPrivateKeyInfo struct {
	Key      string `json:"key"`
	Password string `json:"password"`
}

type BasicAuth struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// POST /api/v1/checkRepos
func handleCheckRequest(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(400)
		return
	}
	var request utils.SearchRequestBody
	err = json.Unmarshal(body, &request)
	if err != nil {
		w.WriteHeader(400)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(checks.CheckAllRepositoriesSpecificChecks(strings.Split(request.Path, ","), request.CheckNames))
}

// GET /api/v1/config
func handleGetConfig(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(convert(utils.Config()))
}

// PUT /api/v1/config
func handlePutConfig(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(400)
		return
	}
	var request ConfigDto
	err = json.Unmarshal(body, &request)
	if err != nil {
		w.WriteHeader(400)
		return
	}
	err = updateConfig(&request)
	if err != nil {
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(err)
		return
	} else {
		w.WriteHeader(200)
	}
}

// PUT /api/v1/config/sshkey
func handlePutSshKey(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(400)
		return
	}
	var request SshPrivateKeyInfo
	err = json.Unmarshal(body, &request)
	if err != nil {
		w.WriteHeader(400)
		return
	}
	err = updateSshKey(&request)
	if err != nil {
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(err)
		return
	} else {
		w.WriteHeader(200)
	}
}

// PUT /api/v1/config/basicauth
func handlePutBasicAuth(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(400)
		return
	}
	var request BasicAuth
	err = json.Unmarshal(body, &request)
	if err != nil {
		w.WriteHeader(400)
		return
	}
	err = updateBasicAuth(&request)
	if err != nil {
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(err)
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
		checkNames = append(checkNames, utils.FunctionName(v))
	}
	json.NewEncoder(w).Encode(checkNames)
}

func convert(cfg utils.GlobalConfig) *ConfigDto {
	return &ConfigDto{
		BranchPattern: cfg.BranchPattern.String(),
		NamePattern:   cfg.NamePattern.String(),
		EmailPattern:  cfg.EmailPattern.String(),
	}
}

func updateConfig(cfg *ConfigDto) error {
	var err error
	var bPat *regexp.Regexp
	var nPat *regexp.Regexp
	var ePat *regexp.Regexp
	if bPat, err = utils.ExtractPattern(cfg.BranchPattern); err != nil {
		return err
	}
	if nPat, err = utils.ExtractPattern(cfg.NamePattern); err != nil {
		return err
	}
	if ePat, err = utils.ExtractPattern(cfg.EmailPattern); err != nil {
		return err
	}
	utils.InitConfig(&utils.GlobalConfig{
		BranchPattern: bPat,
		NamePattern:   nPat,
		EmailPattern:  ePat,
	})
	return nil
}

func updateSshKey(keyInfo *SshPrivateKeyInfo) error {
	return mygit.InitSshKey([]byte(keyInfo.Key), keyInfo.Password)
}

func updateBasicAuth(auth *BasicAuth) error {
	return mygit.InitHttpBasicAuth(auth.Username, auth.Password)
}
