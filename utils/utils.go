package utils

import (
	"fmt"
	"math/rand"

	"github.com/janbaer/github-oauth-bridge/config"
)

// RandomString creates a random string with the desired length
func RandomString(keyLength int) string {
	letter := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	b := make([]rune, keyLength)
	for i := range b {
		b[i] = letter[rand.Intn(len(letter))]
	}
	return string(b)
}

// GetConfigByClientID returns the config with the given clientID from the list of configs
func GetConfigByClientID(configEntries *[]config.Config, clientID string) (config.Config, error) {
	for _, config := range *configEntries {
		if config.ClientID == clientID {
			return config, nil
		}
	}

	var emptyConfig config.Config
	return emptyConfig, fmt.Errorf("No configuration found for clientId %s", clientID)
}
