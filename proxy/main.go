package main

import (
	"context"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"

	"github.com/coreos/go-oidc"
	"github.com/rs/cors"
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
	StaticFilesDir     string
	ReverseProxy       string
	Local              bool
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
	config.StaticFilesDir = os.Getenv("STATIC_FILE_SERVING")
	config.ReverseProxy = os.Getenv("REVERSE_PROXY")
	config.Local = os.Getenv("LOCAL") != ""

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

	h, err := newHandler(config)

	if err != nil {
		panic(err)
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/proxy", h.proxyHandler)
	mux.HandleFunc("/login", h.loginHandler)
	mux.HandleFunc("/callback", h.callbackHandler)

	if len(config.StaticFilesDir) != 0 {
		mux.Handle("/", http.FileServer(http.Dir(config.StaticFilesDir)))
	} else if len(config.ReverseProxy) != 0 {
		u, err := url.Parse(config.ReverseProxy)

		if err != nil {
			panic(err)
		}

		mux.Handle("/", httputil.NewSingleHostReverseProxy(u))
	}

	port := os.Getenv("PORT")

	if len(port) == 0 {
		port = "80"
	}

	var handler http.Handler = mux
	if config.Local {
		handler = cors.AllowAll().Handler(mux)
	}

	http.ListenAndServe(":"+port, handler)
}
