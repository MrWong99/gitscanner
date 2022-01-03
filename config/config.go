package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server *ServerConfig `json:"server,omitempty" yaml:"server,omitempty"`
	Auth   *AuthConfig   `json:"auth,omitempty" yaml:"auth,omitempty"`
	Checks []CheckConfig `json:"checks,omitempty" yaml:"checks,omitempty"`
}

type ServerConfig struct {
	Port int        `json:"port,omitempty" yaml:"port,omitempty"`
	Tls  *TlsConfig `json:"tls,omitempty" yaml:"tls,omitempty"`
}

type TlsConfig struct {
	CertFile       string `json:"certFile,omitempty" yaml:"certFile,omitempty"`
	PrivateKeyFile string `json:"privateKeyFile,omitempty" yaml:"privateKeyFile,omitempty"`
}

type AuthConfig struct {
	Ssh       *SshConfig       `json:"ssh,omitempty" yaml:"ssh,omitempty"`
	BasicAuth *BasicAuthConfig `json:"basicAuth,omitempty" yaml:"basicAuth,omitempty"`
}

type SshConfig struct {
	PrivateKeyFile string `json:"privateKeyFile,omitempty" yaml:"privateKeyFile,omitempty"`
	KeyPassphrase  string `json:"keyPassphrase,omitempty" yaml:"keyPassphrase,omitempty"`
}

type BasicAuthConfig struct {
	Username string `json:"username,omitempty" yaml:"username,omitempty"`
	Password string `json:"password,omitempty" yaml:"password,omitempty"`
}

type CheckConfig struct {
	Name    string                 `json:"name"`
	Enabled bool                   `json:"enabled"`
	Config  map[string]interface{} `json:"config,omitempty" yaml:"config,omitempty"`
}

var ConfigLocation string = "GrootConfig.yml"

var lastReadConfig *Config

func CurrentConfig() *Config {
	return lastReadConfig
}

func ReadJson(input []byte) (*Config, error) {
	cfg := &Config{}
	err := json.Unmarshal(input, cfg)
	if err != nil {
		return nil, err
	}
	lastReadConfig = cfg
	return cfg, nil
}

func ReadYaml(input []byte) (*Config, error) {
	cfg := &Config{}
	err := yaml.Unmarshal(input, cfg)
	if err != nil {
		return nil, err
	}
	lastReadConfig = cfg
	return cfg, nil
}

func (cfg *Config) AsJson(isPretty bool) ([]byte, error) {
	if isPretty {
		return json.MarshalIndent(cfg, "", "  ")
	} else {
		return json.Marshal(cfg)
	}
}

func (cfg *Config) AsYaml() ([]byte, error) {
	return yaml.Marshal(cfg)
}

func (cfg *Config) AddOrUpdateCheckConfig(checkCfg *CheckConfig) {
	for i, c := range cfg.Checks {
		if c.Name == checkCfg.Name {
			newSlice := cfg.Checks[:i]
			newSlice = append(newSlice, *checkCfg)
			newSlice = append(newSlice, cfg.Checks[i+1:]...)
			cfg.Checks = newSlice
			return
		}
	}
	cfg.Checks = append(cfg.Checks, *checkCfg)
}

func UpdateConfigFile() error {
	fileinfo, err := os.Stat(ConfigLocation)
	if err != nil {
		return err
	}
	var parsedBytes []byte
	if path.Ext(ConfigLocation) == ".json" {
		parsedBytes, err = lastReadConfig.AsJson(true)
	} else {
		parsedBytes, err = lastReadConfig.AsYaml()
	}
	if err != nil {
		return err
	}
	return ioutil.WriteFile(ConfigLocation, parsedBytes, fileinfo.Mode())
}

func (cfg *Config) CheckIsEnabled(name string) bool {
	for _, check := range cfg.Checks {
		if check.Name == name {
			return check.Enabled
		}
	}
	return true
}
