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

// AuthCallback - Handles the authCallback request
func AuthCallback(w http.ResponseWriter, r *http.Request) {
	code := r.FormValue("code")
	callbackState := r.FormValue("state")

	clientID, err := state.DecryptState(callbackState, os.Getenv("SECRET"))
	if err != nil {
		http.Error(w, "State could not be verified", http.StatusBadRequest)
		return
	}

	configValue, err := config.ReadConfigFromEnv(clientID)
	if err != nil {
		log.Printf("Error while verifying state: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	token, err := github.Exchange(configValue.ClientID, configValue.ClientSecretID, code)
	if err != nil {
		log.Printf("Error while exchange code %s for client %s with Github: %v", code, configValue.ClientID, err)
		http.Error(w, "Code was not accepted by the Oauth provider", http.StatusBadRequest)
		return
	}

	redirectURLWithToken := fmt.Sprintf("%s?token=%s", configValue.RedirectURL, token)

	w.Header().Set("Location", redirectURLWithToken)
	w.WriteHeader(http.StatusTemporaryRedirect)
}
