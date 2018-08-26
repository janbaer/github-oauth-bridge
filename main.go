package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/stretchr/gomniauth"
	"github.com/stretchr/gomniauth/oauth2"
	"github.com/stretchr/gomniauth/providers/github"
	"github.com/stretchr/objx"
	"github.com/stretchr/signature"
)

// Config defines the sub config object
type Config struct {
	APIURL         string `json:"apiUrl"`
	RedirectURL    string `json:"redirectUrl"`
	ClientID       string `json:"clientId"`
	ClientSecretID string `json:"clientSecretId"`
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	isProd := false

	if os.Getenv("ENV") == "PROD" {
		isProd = true
		log.Println("Reading the configuration for PROD")
	}

	configEntries := readConfig(isProd)

	http.HandleFunc("/healthcheck", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			respondErr(w, r, http.StatusBadRequest, r.Method, " is not allowed")
			return
		}

		respond(w, r, http.StatusOK, fmt.Sprintf("Request to %s OK", r.RequestURI))
	})

	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			respondErr(w, r, http.StatusBadRequest, r.Method, " is not allowed")
			return
		}

		clientID := r.URL.Query().Get("clientId")
		if len(clientID) == 0 {
			respondErr(w, r, http.StatusBadRequest, "The query parameter clientId is required")
			return
		}

		config, err := getConfigByClientID(configEntries, clientID)
		if err != nil {
			respondErr(w, r, http.StatusBadRequest, err)
			return
		}

		setupGomniauth(config)

		providerName := "github"
		provider, err := gomniauth.Provider(providerName)

		if err != nil {
			log.Fatalln("Error trying to get provider", providerName)
		}

		loginURL, err := provider.GetBeginAuthURL(nil, objx.MSI(oauth2.OAuth2KeyScope, "public_repo"))

		if err != nil {
			log.Fatalln("Error trying to GetBeginAuthURL", providerName, "-", err)
		}

		w.Header().Set("Location", loginURL)
		w.WriteHeader(http.StatusTemporaryRedirect)
	})

	http.HandleFunc("/auth/callback", func(w http.ResponseWriter, r *http.Request) {
		providerName := "github"
		provider, err := gomniauth.Provider(providerName)
		if err != nil {
			log.Fatalln("Error when trying to get provider", providerName, "-", err)
		}

		// get the credentials
		creds, err := provider.CompleteAuth(objx.MustFromURLQuery(r.URL.RawQuery))
		if err != nil {
			log.Fatalln("Error when trying to complete auth for", providerName, "-", err)
		}

		clientID := r.URL.Query().Get("clientId")
		config, _ := getConfigByClientID(configEntries, clientID)

		redirectURLWithToken := config.RedirectURL + "?token=" + creds.Get("access_token").String()

		w.Header().Set("Location", redirectURLWithToken)
		w.WriteHeader(http.StatusTemporaryRedirect)
	})

	log.Println("Listening for requests on port: " + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func setupGomniauth(config Config) {
	gomniauth.SetSecurityKey(signature.RandomKey(64))

	callbackURL := fmt.Sprintf("%s/auth/callback?clientId=%s", config.APIURL, config.ClientID)
	clientID := config.ClientID
	clientSecret := config.ClientSecretID
	gomniauth.WithProviders(github.New(clientID, clientSecret, callbackURL))
}

func getConfigByClientID(configEntries []Config, clientID string) (Config, error) {
	for _, config := range configEntries {
		if config.ClientID == clientID {
			return config, nil
		}
	}

	var emptyConfig Config
	return emptyConfig, fmt.Errorf("No configuration found for clientId %s", clientID)
}

func readConfig(isProd bool) []Config {
	configFile := "./config.prod.json"
	if !isProd {
		configFile = "./config.dev.json"
	}

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

func encodeBody(w http.ResponseWriter, r *http.Request, v interface{}) error {
	return json.NewEncoder(w).Encode(v)
}

func respond(w http.ResponseWriter, r *http.Request, status int, data interface{}) {
	w.WriteHeader(status)
	if data != nil {
		encodeBody(w, r, data)
	}
}

func respondErr(w http.ResponseWriter, r *http.Request, status int, args ...interface{}) {
	respond(w, r, status, map[string]interface{}{
		"error": map[string]interface{}{
			"message": fmt.Sprint(args...),
		},
	})
}
