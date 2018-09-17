package handlers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/janbaer/github-oauth-bridge/config"
	"github.com/janbaer/github-oauth-bridge/github"
	"github.com/janbaer/github-oauth-bridge/utils"
)

// HandleLogin implements the handler for the login request from the client and redirect to Github
func HandleLogin(configEntries *[]config.Config, keyStore map[string]string, authCodeURLFunc github.AuthCodeURLFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, fmt.Sprintf("The %s method is not allowed", r.Method), http.StatusBadRequest)
			return
		}

		clientID := r.URL.Query().Get("clientId")
		if len(clientID) == 0 {
			http.Error(w, "The query parameter clientId is required", http.StatusBadRequest)
			return
		}

		config, err := utils.GetConfigByClientID(configEntries, clientID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		randomKey := utils.RandomString(20)
		keyStore[randomKey] = clientID

		url := authCodeURLFunc(config.ClientID, config.ClientSecretID, randomKey)
		log.Printf("Redirect to %s", url)
		http.Redirect(w, r, url, http.StatusTemporaryRedirect)
	})
}

// HandleAuthCallback implements the handler for the callback from Github
func HandleAuthCallback(configEntries *[]config.Config, keyStore map[string]string, exchangeFunc github.ExchangeFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		code := r.FormValue("code")
		randomKey := r.FormValue("state")

		clientID, exists := keyStore[randomKey]
		if !exists {
			http.Error(w, "Validation of state failed", http.StatusBadRequest)
			return
		}
		delete(keyStore, randomKey)

		config, _ := utils.GetConfigByClientID(configEntries, clientID)

		token, err := exchangeFunc(config.ClientID, config.ClientSecretID, code)
		if err != nil {
			http.Error(w, "Code was not accepted by the Oauth provider", http.StatusBadRequest)
			return
		}

		redirectURLWithToken := fmt.Sprintf("%s?token=%s", config.RedirectURL, token)

		w.Header().Set("Location", redirectURLWithToken)
		w.WriteHeader(http.StatusTemporaryRedirect)
	})
}
