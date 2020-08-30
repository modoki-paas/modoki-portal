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

type Config struct {
	ClusterAddress     string
	CAData             string
	ClientID           string
	ClientSecret       string
	Scopes             []string
	RedirectURL        string
	IssuerURL          string
	SessionStoreSecret string
}

func initConfig() *Config {
	var config Config

	config.ClusterAddress = os.Getenv("K8S_CLUSTER_ADDRESS")
	config.CAData = os.Getenv("K8S_CERTIFICATE_AUTHORITY_DATA")
	config.ClientID = os.Getenv("OIDC_CLIENT_ID")
	config.ClientSecret = os.Getenv("OIDC_CLIENT_SECRET")
	config.Scopes = strings.Split(os.Getenv("OIDC_SCOPES"), ",")
	config.RedirectURL = os.Getenv("OIDC_REDIRECT_URL")
	config.IssuerURL = os.Getenv("OIDC_ISSUER_URL")
	config.SessionStoreSecret = os.Getenv("SESSION_STORE_SECRET")

	if config.SessionStoreSecret == "" {
		config.SessionStoreSecret = "very-secure-secret"
	}

	return &config
}

type Authenticator struct {
	Provider *oidc.Provider
	Config   oauth2.Config
}

func NewAuthenticator(ctx context.Context, config *Config) (*Authenticator, error) {
	provider, err := oidc.NewProvider(ctx, config.IssuerURL)
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
	}, nil
}

func main() {
	config := initConfig()

	handler, err := newHandler(config)

	if err != nil {
		panic(err)
	}

	http.HandleFunc("/proxy", handler.proxyHandler)
	http.HandleFunc("/login", handler.loginHandler)
	http.HandleFunc("/callback", handler.callbackHandler)

	port := os.Getenv("PORT")

	if len(port) == 0 {
		port = "80"
	}

	http.ListenAndServe(":"+port, nil)
}
