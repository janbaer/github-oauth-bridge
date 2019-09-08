package state_test

import (
	"testing"

	"github.com/janbaer/github-oauth-bridge/state"
	"github.com/stretchr/testify/assert"
)

func Test_EncryptState(t *testing.T) {
	key := "123password"
	encryptedState := state.EncryptState("clientID", key)
	assert.NotEmpty(t, encryptedState)
}

func Test_DecryptState(t *testing.T) {
	clientID := "clientID"
	key := "123password"

	encryptedState := state.EncryptState(clientID, key)
	decryptedClientID, _ := state.DecryptState(encryptedState, key)
	assert.Equal(t, clientID, decryptedClientID, "The decrypted state should be the same as the expected")
}
