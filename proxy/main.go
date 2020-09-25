package main

import (
	"context"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"

	"github.com/coreos/go-oidc"
	"github.com/modoki-paas/modoki-portal/proxy/config"
	"github.com/rs/cors"
	"golang.org/x/oauth2"
)

type Config = config.Config

func initConfig() *Config {
	config, err := config.ReadConfig()

	if err != nil {
		panic(err)
	}

	return config
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
	mux.HandleFunc("/github/login", h.loginGitHub)
	mux.HandleFunc("/github/callback", h.callbackGitHub)
	mux.HandleFunc("/github/installations", h.listInstallations)
	mux.HandleFunc("/github/repositories", h.listRepositories)

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
