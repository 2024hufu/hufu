package config

import (
	"encoding/hex"
	"fmt"
	"os"

	"github.com/FISCO-BCOS/go-sdk/v3/client"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Database struct {
		Username string `yaml:"username"`
		Password string `yaml:"password"`
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		DBName   string `yaml:"dbname"`
		Charset  string `yaml:"charset"`
	} `yaml:"database"`

	Fisco struct {
		IsSMCrypto  bool   `yaml:"is_sm_crypto"`
		GroupID     string `yaml:"group_id"`
		DisableSsl  bool   `yaml:"disable_ssl"`
		Host        string `yaml:"host"`
		Port        int    `yaml:"port"`
		TLSCaFile   string `yaml:"tls_ca_file"`
		TLSKeyFile  string `yaml:"tls_key_file"`
		TLSCertFile string `yaml:"tls_cert_file"`
	} `yaml:"fisco"`
}

var GlobalConfig Config

func LoadConfig(configPath string) (*Config, error) {
	config := &Config{}

	file, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("error reading config file: %v", err)
	}

	err = yaml.Unmarshal(file, config)
	if err != nil {
		return nil, fmt.Errorf("error parsing config file: %v", err)
	}

	GlobalConfig = *config
	return config, nil
}

func ReadConfig(pk string) *client.Config {
	privateKey, _ := hex.DecodeString(pk)

	config := &client.Config{
		IsSMCrypto:  GlobalConfig.Fisco.IsSMCrypto,
		GroupID:     GlobalConfig.Fisco.GroupID,
		DisableSsl:  GlobalConfig.Fisco.DisableSsl,
		PrivateKey:  privateKey,
		Host:        GlobalConfig.Fisco.Host,
		Port:        GlobalConfig.Fisco.Port,
		TLSCaFile:   GlobalConfig.Fisco.TLSCaFile,
		TLSKeyFile:  GlobalConfig.Fisco.TLSKeyFile,
		TLSCertFile: GlobalConfig.Fisco.TLSCertFile,
	}

	return config
}
