package api

import (
	"fmt"
	"net/http"
	"os"

	"github.com/janbaer/github-oauth-bridge/config"
)

var keyMap map[string]string
var configEntries []config.Config

func init() {
	configFile := "./../config.prod.json"
	if os.Getenv("ENV") == "DEV" {
		configFile = "./config.dev.json"
	}

	dir, _ := os.Getwd()
	fmt.Printf("Working in %s", dir)

	configEntries = config.ReadConfig(configFile)
	fmt.Printf("Reading %d configEntries", len(configEntries))
}

// Login - Handles the login request
func Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, fmt.Sprintf("The %s method is not allowed", r.Method), http.StatusBadRequest)
		return
	}

	clientID := r.URL.Query().Get("clientId")
	if len(clientID) == 0 {
		http.Error(w, "The query parameter clientId is required", http.StatusBadRequest)
		return
	}

	fmt.Fprintf(w, "Hello from loginHandler with clientId %s", clientID)
}
