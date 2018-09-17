package github

import (
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

// AuthCodeURLFunc is an adapter for the AuthCodeURL function
type AuthCodeURLFunc func(clientID string, clientSecretID string, state string) string

// ExchangeFunc is an adapter for the Exchange function
type ExchangeFunc func(clientID string, clientSecretID string, code string) (string, error)

// AuthCodeURL creates the url to redirect to Github to start authentication
func AuthCodeURL(clientID string, clientSecretID string, state string) string {
	oauthConfig := createOauth2Config(clientID, clientSecretID)
	return oauthConfig.AuthCodeURL(state)
}

// Exchange - Sends the code to Github to return the AccessToken
func Exchange(clientID string, clientSecretID string, code string) (string, error) {
	oauthConfig := createOauth2Config(clientID, clientSecretID)
	token, err := oauthConfig.Exchange(oauth2.NoContext, code)
	if err != nil {
		return "", err
	}
	return token.AccessToken, nil
}

func createOauth2Config(clientID string, clientSecretID string) *oauth2.Config {
	return &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecretID,
		Scopes:       []string{"public_repo"},
		Endpoint:     github.Endpoint,
	}
}
