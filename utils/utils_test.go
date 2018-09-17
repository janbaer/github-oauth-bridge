package utils_test

import (
	"testing"

	"github.com/janbaer/github-oauth-bridge/config"
	"github.com/janbaer/github-oauth-bridge/utils"
	"github.com/stretchr/testify/assert"
)

func TestRandomStringLength(t *testing.T) {
	expectedLength := 20
	assert.Equal(t, len(utils.RandomString(expectedLength)), expectedLength)
}

func TestConfigByClientID_with_known_clientID(t *testing.T) {
	clientID := "12345678"
	configEntries := []config.Config{config.Config{ClientID: clientID}}

	result, err := utils.GetConfigByClientID(&configEntries, clientID)

	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, result.ClientID, clientID)
}

func TestConfigByClientID_with_unknown_clientID(t *testing.T) {
	clientID := "12345678"
	configEntries := []config.Config{config.Config{ClientID: clientID}}

	result, err := utils.GetConfigByClientID(&configEntries, "535252523")

	assert.Error(t, err)
	assert.Empty(t, result)
}
