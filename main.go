package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/janbaer/github-oauth-bridge/config"
	oauthGithub "github.com/janbaer/github-oauth-bridge/github"
	"github.com/janbaer/github-oauth-bridge/handlers"

	"golang.org/x/oauth2"
)

var keyMap map[string]string
var isProd = false
var configEntries []config.Config

// Oauth2Handler defines the interface to get the token from the OauthProvider
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

	configFile := "./config.dev.json"
	if isProd {
		configFile = "./config.prod.json"
	}
	configEntries = config.ReadConfig(configFile)
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	mux := newRouter()

	log.Printf("Listening for requests on port: %s\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), mux))
}

func newRouter() *http.ServeMux {
	mux := http.NewServeMux()

	mux.Handle("/login", handlers.HandleLogin(&configEntries, keyMap, oauthGithub.AuthCodeURL))
	mux.Handle("/auth/callback", handlers.HandleAuthCallback(&configEntries, keyMap, oauthGithub.Exchange))

	return mux
}
