package handlers_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/janbaer/github-oauth-bridge/config"
	"github.com/janbaer/github-oauth-bridge/handlers"
	"github.com/stretchr/testify/assert"
)

var (
	token = "token"
	key   = "ABCDEFG"
)

var testConfig = config.Config{ClientID: "12345678", ClientSecretID: "ABCDEFG", RedirectURL: "http://any-server.de"}
var configEntries = []config.Config{testConfig}

var authCodeURLMock = func(clientID, clientSecretID, state string) string {
	return "http://github.com/..."
}

var exchangeFuncMock = func(clientID string, clientSecretID string, code string) (string, error) {
	return token, nil
}

func TestAuthCodeURL_when_clientID_is_valid(t *testing.T) {
	var authCodeURLMock = func(clientID, clientSecretID, state string) string {
		assert.Equal(t, clientID, testConfig.ClientID, "Expected clientID was not used")
		assert.Equal(t, clientSecretID, testConfig.ClientSecretID, "Expected clientSecretID was not used")
		return "http://github.com/..."
	}

	keyStore := make(map[string]string)
	requestURL := fmt.Sprintf("/login?clientId=%s", testConfig.ClientID)

	request, _ := http.NewRequest(http.MethodGet, requestURL, nil)
	recorder := httptest.NewRecorder()

	loginHandler := handlers.HandleLogin(&configEntries, keyStore, authCodeURLMock)
	loginHandler.ServeHTTP(recorder, request)

	assert.Equal(t, http.StatusTemporaryRedirect, recorder.Code)
}

func TestAuthCodeURL_when_no_clientID_was_found_in_the_request_url(t *testing.T) {
	keyStore := make(map[string]string)

	request, _ := http.NewRequest(http.MethodGet, "/login", nil)
	recorder := httptest.NewRecorder()

	loginHandler := handlers.HandleLogin(&configEntries, keyStore, authCodeURLMock)
	loginHandler.ServeHTTP(recorder, request)

	assert.Equal(t, http.StatusBadRequest, recorder.Code)
}

func TestAuthCodeURL_when_clientID_is_not_valid(t *testing.T) {
	keyStore := make(map[string]string)

	request, _ := http.NewRequest(http.MethodGet, "/login?clientId=53535235", nil)
	recorder := httptest.NewRecorder()

	loginHandler := handlers.HandleLogin(&configEntries, keyStore, authCodeURLMock)
	loginHandler.ServeHTTP(recorder, request)

	assert.Equal(t, http.StatusBadRequest, recorder.Code)
}

func TestAuthCodeURL_when_a_post_was_sent(t *testing.T) {
	keyStore := make(map[string]string)

	request, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("/login?clientId=%s", testConfig.ClientID), nil)
	recorder := httptest.NewRecorder()

	loginHandler := handlers.HandleLogin(&configEntries, keyStore, authCodeURLMock)
	loginHandler.ServeHTTP(recorder, request)

	assert.Equal(t, http.StatusBadRequest, recorder.Code)
}

func TestAuthCallback_when_a_valid_key_was_sent(t *testing.T) {
	expectedURL := fmt.Sprintf("%s?token=%s", testConfig.RedirectURL, token)

	keyStore := make(map[string]string)
	keyStore[key] = testConfig.ClientID

	request, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/callback?code=%s&state=%s", token, key), nil)
	recorder := httptest.NewRecorder()

	loginHandler := handlers.HandleAuthCallback(&configEntries, keyStore, exchangeFuncMock)
	loginHandler.ServeHTTP(recorder, request)

	assert.Equal(t, http.StatusTemporaryRedirect, recorder.Code)
	assert.Equal(t, expectedURL, recorder.Header().Get("Location"))
}

func TestAuthCallback_when_a_invalid_key_was_sent(t *testing.T) {
	keyStore := make(map[string]string)
	keyStore[key] = testConfig.ClientID

	request, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/callback?code=%s&state=invalid_key", token), nil)
	recorder := httptest.NewRecorder()

	loginHandler := handlers.HandleAuthCallback(&configEntries, keyStore, exchangeFuncMock)
	loginHandler.ServeHTTP(recorder, request)

	assert.Equal(t, http.StatusBadRequest, recorder.Code)
}

func TestAuthCallback_when_a_exchangeFunc_return_error(t *testing.T) {
	exchangeFuncMock := func(clientID string, clientSecretID string, code string) (string, error) {
		return "", fmt.Errorf("Invalid code %s was sent", code)
	}

	keyStore := make(map[string]string)
	keyStore[key] = testConfig.ClientID

	request, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/callback?code=invalid_code&state=%s", key), nil)
	recorder := httptest.NewRecorder()

	loginHandler := handlers.HandleAuthCallback(&configEntries, keyStore, exchangeFuncMock)
	loginHandler.ServeHTTP(recorder, request)

	assert.Equal(t, http.StatusBadRequest, recorder.Code)
}
