package api

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/janbaer/github-oauth-bridge/config"
	"github.com/janbaer/github-oauth-bridge/github"
	"github.com/janbaer/github-oauth-bridge/state"
)

// Login - Handles the login request
func Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, fmt.Sprintf("The %s method is not allowed", r.Method), http.StatusBadRequest)
		return
	}

	clientID := r.URL.Query().Get("clientId")
	if len(clientID) == 0 {
		http.Error(w, "The query parameter clientId is required!", http.StatusBadRequest)
		return
	}

	configValue, err := config.ReadConfigFromEnv(clientID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	secret := os.Getenv("SECRET")
	state := state.EncryptState(clientID, secret)

	url := github.AuthCodeURL(configValue.ClientID, configValue.ClientSecretID, state)

	log.Printf("Redirect to %s", url)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}
