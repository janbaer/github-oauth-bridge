package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

// Config defines the sub config object
type Config struct {
	RedirectURL    string `json:"redirectUrl"`
	ClientID       string `json:"clientId"`
	ClientSecretID string `json:"clientSecretId"`
}

var keyMap map[string]string

func init() {
	keyMap = make(map[string]string)
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

		oauthConf, err := createOauthConf(&configEntries, clientID)
		if err != nil {
			respondErr(w, r, http.StatusBadRequest, err)
			return
		}

		randomKey := randomString(20)
		keyMap[randomKey] = clientID

		url := oauthConf.AuthCodeURL(randomKey)
		log.Printf("Redirect to %s", url)
		http.Redirect(w, r, url, http.StatusTemporaryRedirect)
	})

	http.HandleFunc("/auth/callback", func(w http.ResponseWriter, r *http.Request) {
		code := r.FormValue("code")
		randomKey := r.FormValue("state")

		clientID, exists := keyMap[randomKey]
		if !exists {
			respondErr(w, r, http.StatusBadRequest, "Validation of state failed")
			return
		}
		delete(keyMap, randomKey)

		config, _ := getConfigByClientID(&configEntries, clientID)
		oauthConf, _ := createOauthConf(&configEntries, clientID)

		token, err := oauthConf.Exchange(oauth2.NoContext, code)
		if err != nil {
			respondErr(w, r, http.StatusNotAcceptable)
			return
		}

		redirectURLWithToken := fmt.Sprintf("%s?token=%s", config.RedirectURL, token.AccessToken)

		w.Header().Set("Location", redirectURLWithToken)
		w.WriteHeader(http.StatusTemporaryRedirect)
	})

	log.Println("Listening for requests on port: " + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func createOauthConf(configEntries *[]Config, clientID string) (*oauth2.Config, error) {
	config, err := getConfigByClientID(configEntries, clientID)
	if err != nil {
		return nil, err
	}

	return &oauth2.Config{
		ClientID:     config.ClientID,
		ClientSecret: config.ClientSecretID,
		Scopes:       []string{"public_repo"},
		Endpoint:     github.Endpoint,
	}, nil
}

func getConfigByClientID(configEntries *[]Config, clientID string) (Config, error) {
	for _, config := range *configEntries {
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

func randomString(keyLength int) string {
	letter := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	b := make([]rune, keyLength)
	for i := range b {
		b[i] = letter[rand.Intn(len(letter))]
	}
	return string(b)
}
