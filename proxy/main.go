package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/coreos/go-oidc"
	"golang.org/x/oauth2"
)

type Request struct {
	Method  string            `json:"method"`
	URL     string            `json:"url"`
	Headers map[string]string `json:"headers"`
	Body    []byte            `json:"body"`
}

var config struct {
	ClientID           string
	ClientSecret       string
	Scopes             []string
	RedirectURL        string
	ProviderURL        string
	SessionStoreSecret string
}

func init() {
	config.ClientID = os.Getenv("OIDC_CLIENT_ID")
	config.ClientSecret = os.Getenv("OIDC_CLIENT_SECRET")
	config.Scopes = strings.Split(os.Getenv("OIDC_SCOPES"), ",")
	config.RedirectURL = os.Getenv("OIDC_REDIRECT_URL")
	config.ProviderURL = os.Getenv("OIDC_PROVIDER_URL")
	config.SessionStoreSecret = os.Getenv("SESSION_STORE_SECRET")

	if config.SessionStoreSecret == "" {
		config.SessionStoreSecret = "very-secure-secret"
	}
}

type Authenticator struct {
	Provider *oidc.Provider
	Config   oauth2.Config
	Ctx      context.Context
}

func NewAuthenticator() (*Authenticator, error) {
	ctx := context.Background()

	provider, err := oidc.NewProvider(ctx, config.ProviderURL)
	if err != nil {
		log.Printf("failed to get provider: %v", err)
		return nil, err
	}

	conf := oauth2.Config{
		ClientID:     config.ClientID,
		ClientSecret: config.ClientSecret,
		RedirectURL:  config.RedirectURL,
		Endpoint:     provider.Endpoint(),
		Scopes:       config.Scopes,
	}

	return &Authenticator{
		Provider: provider,
		Config:   conf,
		Ctx:      ctx,
	}, nil
}

func main() {
	handler := newHandler(config.SessionStoreSecret)

	http.HandleFunc("/proxy", handler.proxyHandler)

	http.ListenAndServe(":80", nil)
}
