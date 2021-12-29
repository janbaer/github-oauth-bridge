package config

import (
	"fmt"
	"log"
	"os"
	"strings"
)

// Config defines the sub config object
type Config struct {
	RedirectURL    string `json:"redirectUrl"`
	ClientID       string `json:"clientId"`
	ClientSecretID string `json:"clientSecretId"`
}

// ReadConfigFromEnv - Reads the config for the given clientID from the Environment variable
func ReadConfigFromEnv(clientID string) (Config, error) {
	var config Config

	configValue := os.Getenv(fmt.Sprintf("CLIENT_%s", clientID))
	if len(configValue) == 0 {
		return config, fmt.Errorf("No configuration for clientID %s found", clientID)
	}

	config = parseConfig(clientID, configValue)
	return config, nil
}

func parseConfig(clientID string, configValue string) Config {
	var clientSecretID string
	var redirectURL string

	values := strings.Split(configValue, "|")
	clientSecretID, redirectURL = values[0], values[1]

	log.Println("RedirectURL", configValue)

	return Config{
		ClientID:       clientID,
		ClientSecretID: clientSecretID,
		RedirectURL:    redirectURL,
	}
}
