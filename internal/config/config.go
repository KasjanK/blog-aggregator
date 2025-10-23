package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

const configFileName = ".gatorconfig.json"

type Config struct {
	DatabaseUrl string `json:"db_url"`	
	CurrentUserName string `json:"current_user_name"`
}

func getConfigFilePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	fullConfigFilePath := filepath.Join(home, configFileName)	
	return fullConfigFilePath, nil
}

func write(cfg Config) error {
	configFilePath, err := getConfigFilePath()
	if err != nil {
		return err
	}

	configFile, err := os.Create(configFilePath)
	if err != nil {
		return err
	}

	defer configFile.Close()

	encoder := json.NewEncoder(configFile)
	err = encoder.Encode(cfg)
	if err != nil {
		return err
	}
	return nil
}

func Read() (Config, error) {
	configFilePath, err := getConfigFilePath()
	if err != nil {
		return Config{}, err
	}

	configFile, err := os.Open(configFilePath)
	if err != nil {
		return Config{}, err
	}

	defer configFile.Close()

	decoder := json.NewDecoder(configFile)
	cfg := Config{}
	err = decoder.Decode(&cfg)
	if err != nil {
		return Config{}, err
	}

	return cfg, nil 
}

func (cfg *Config) SetUser(username string) error {
	cfg.CurrentUserName = username
	return write(*cfg)
}
