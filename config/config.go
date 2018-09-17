package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

// Config defines the sub config object
type Config struct {
	RedirectURL    string `json:"redirectUrl"`
	ClientID       string `json:"clientId"`
	ClientSecretID string `json:"clientSecretId"`
}

// ReadConfig reads the config from the given config file
func ReadConfig(configFile string) []Config {
	_, err := os.Stat(configFile)
	if err != nil {
		log.Fatalf("Configurationfile %s not found", configFile)
	}

	content, err := ioutil.ReadFile(configFile)
	if err != nil {
		log.Fatalf("Configurationfile %s could not be read", configFile)
	}

	var configEntries []Config
	err = json.Unmarshal(content, &configEntries)
	if err != nil {
		log.Fatalf("Configurationfile %s could not be parsed to the expected json structure", configFile)
	}

	return configEntries
}
