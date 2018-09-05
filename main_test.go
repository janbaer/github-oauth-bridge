package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"golang.org/x/oauth2"

	"github.com/stretchr/testify/assert"
)

type MockOauth2Handler struct {
	ExchangeFunc    func(ctx context.Context, code string, opts ...oauth2.AuthCodeOption) (*oauth2.Token, error)
	AuthCodeURLFunc func(state string, opts ...oauth2.AuthCodeOption) string
}

func (m MockOauth2Handler) Exchange(ctx context.Context, code string, opts ...oauth2.AuthCodeOption) (*oauth2.Token, error) {
	return m.ExchangeFunc(ctx, code, opts...)
}

func (m MockOauth2Handler) AuthCodeURL(state string, opts ...oauth2.AuthCodeOption) string {
	return m.AuthCodeURLFunc(state, opts...)
}

func TestLoginWithValidClientId(t *testing.T) {
	oauth2HandlerMock := MockOauth2Handler{
		ExchangeFunc: func(ctx context.Context, code string, opts ...oauth2.AuthCodeOption) (*oauth2.Token, error) {
			return &oauth2.Token{}, nil
		},
		AuthCodeURLFunc: func(state string, opts ...oauth2.AuthCodeOption) string {
			return ""
		},
	}

	oauth2HandlerFunc := func(clientID string) (Oauth2Handler, error) {
		return oauth2HandlerMock, nil
	}

	keys := make(map[string]string)

	request, _ := http.NewRequest(http.MethodGet, "/login?clientId=3d99a5ff8d7eedf0bb99", nil)

	recorder := httptest.NewRecorder()

	loginHandler := handleLogin(keys, oauth2HandlerFunc)

	loginHandler.ServeHTTP(recorder, request)

	assert.Equal(t, http.StatusTemporaryRedirect, recorder.Code)
}

func TestLoginWithInvalidClientId(t *testing.T) {
	oauth2HandlerMock := MockOauth2Handler{
		ExchangeFunc: func(ctx context.Context, code string, opts ...oauth2.AuthCodeOption) (*oauth2.Token, error) {
			return &oauth2.Token{}, nil
		},
		AuthCodeURLFunc: func(state string, opts ...oauth2.AuthCodeOption) string {
			return ""
		},
	}

	oauth2HandlerFunc := func(clientID string) (Oauth2Handler, error) {
		return oauth2HandlerMock, nil
	}

	keys := make(map[string]string)

	request, _ := http.NewRequest(http.MethodGet, "/login?clientId=123456", nil)

	recorder := httptest.NewRecorder()

	loginHandler := handleLogin(keys, oauth2HandlerFunc)

	loginHandler.ServeHTTP(recorder, request)

	assert.Equal(t, http.StatusBadRequest, recorder.Code, "We expected here a BadRequest since the clientId was wrong")
}
