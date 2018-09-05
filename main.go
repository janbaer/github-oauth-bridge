package main

import (
	"context"
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
var isProd = false
var configEntries []Config

// Oauth2Handler - defines the interface to get the token from the OauthProvider
type Oauth2Handler interface {
	Exchange(ctx context.Context, code string, opts ...oauth2.AuthCodeOption) (*oauth2.Token, error)
	AuthCodeURL(state string, opts ...oauth2.AuthCodeOption) string
}

func init() {
	keyMap = make(map[string]string)

	if os.Getenv("ENV") == "PROD" {
		isProd = true
		log.Println("Server is running on Production environment...")
	}

	configEntries = readConfig(isProd)
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	mux := newRouter()

	log.Println("Listening for requests on port: " + port)
	log.Fatal(http.ListenAndServe(":"+port, mux))
}

func newRouter() *http.ServeMux {
	mux := http.NewServeMux()

	auth2HandlerFunc := func(clientID string) (Oauth2Handler, error) {
		return createOauthConf(&configEntries, clientID)
	}

	mux.Handle("/login", handleLogin(keyMap, auth2HandlerFunc))
	mux.Handle("/auth/callback", handleAuthCallback(keyMap, auth2HandlerFunc))

	return mux
}

func handleLogin(keys map[string]string, oauth2HandlerFunc func(clientID string) (Oauth2Handler, error)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, fmt.Sprintf("The %s methos is not allowed", r.Method), http.StatusBadRequest)
			return
		}

		clientID := r.URL.Query().Get("clientId")
		if len(clientID) == 0 {
			http.Error(w, "The query parameter clientId is required", http.StatusBadRequest)
			return
		}

		oauthConf, err := oauth2HandlerFunc(clientID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		randomKey := randomString(20)
		keys[randomKey] = clientID

		url := oauthConf.AuthCodeURL(randomKey)
		log.Printf("Redirect to %s", url)
		http.Redirect(w, r, url, http.StatusTemporaryRedirect)
	})
}

func handleAuthCallback(keys map[string]string, oauth2HandlerFunc func(clientID string) (Oauth2Handler, error)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		code := r.FormValue("code")
		randomKey := r.FormValue("state")

		clientID, exists := keys[randomKey]
		if !exists {
			http.Error(w, "Validation of state failed", http.StatusBadRequest)
			return
		}
		delete(keys, randomKey)

		config, _ := getConfigByClientID(&configEntries, clientID)
		oauthConf, _ := oauth2HandlerFunc(clientID)

		token, err := oauthConf.Exchange(oauth2.NoContext, code)
		if err != nil {
			http.Error(w, "Code was not accepted by the Oauth provider", http.StatusBadRequest)
			return
		}

		redirectURLWithToken := fmt.Sprintf("%s?token=%s", config.RedirectURL, token.AccessToken)

		w.Header().Set("Location", redirectURLWithToken)
		w.WriteHeader(http.StatusTemporaryRedirect)
	})
}

func createOauthConf(configEntries *[]Config, clientID string) (Oauth2Handler, error) {
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

func randomString(keyLength int) string {
	letter := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	b := make([]rune, keyLength)
	for i := range b {
		b[i] = letter[rand.Intn(len(letter))]
	}
	return string(b)
}
