package main

import (
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/coreos/go-oidc"
	"github.com/gorilla/sessions"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
	_ "k8s.io/client-go/plugin/pkg/client/auth/oidc"
	restclient "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd/api"
)

const (
	cookieKey = "auth-session"
)

type handler struct {
	cfg       *Config
	store     *sessions.CookieStore
	transport http.RoundTripper

	github *oauth2.Config
}

func newHandler(cfg *Config) (*handler, error) {
	store := sessions.NewCookieStore([]byte(cfg.SessionStoreSecret))

	h := &handler{
		cfg:   cfg,
		store: store,
		github: &oauth2.Config{
			ClientID:     cfg.GitHub.ClientID,
			ClientSecret: cfg.GitHub.ClientSecret,
			Endpoint:     github.Endpoint,
		},
	}
	gob.Register(map[string]interface{}{})
	gob.Register(&oauth2.Token{})

	var tlsConf *tls.Config
	if len(cfg.CAData) != 0 {
		caDataRaw, err := base64.StdEncoding.DecodeString(cfg.CAData)

		if err != nil {
			return nil, err
		}

		roots := x509.NewCertPool()
		ok := roots.AppendCertsFromPEM(caDataRaw)
		if !ok {
			return nil, fmt.Errorf("invalid cerificate for CA")
		}

		tlsConf = &tls.Config{RootCAs: roots}
	}

	h.transport = &http.Transport{
		TLSClientConfig: tlsConf,
	}

	return h, nil
}

func (h *handler) proxyHandler(rw http.ResponseWriter, req *http.Request) {
	session, err := h.store.Get(req, cookieKey)

	if err != nil {
		rw.Write([]byte(`{"message": "not signed in"}`))
		rw.WriteHeader(http.StatusUnauthorized)

		return
	}

	idToken := getParam(session, "id_token")
	refreshToken := getParam(session, "refresh_token")

	if len(idToken)+len(refreshToken) == 0 {
		http.Redirect(rw, req, "/login", http.StatusFound)

		return
	}

	provider, err := restclient.GetAuthProvider(h.cfg.ClusterAddress, &api.AuthProviderConfig{
		Name: "oidc",
		Config: map[string]string{
			"idp-issuer-url": h.cfg.IssuerURL,
			"client-id":      h.cfg.ClientID,
			"client-secret":  h.cfg.ClientSecret,
			"id-token":       idToken,
			"refresh-token":  refreshToken,
		},
	}, NewCookieConfigPersister(session, rw, req))

	if err != nil {
		log.Printf("GetAuthProvider failed: %+v", err)

		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	transport := provider.WrapTransport(h.transport)

	innerReq := &Request{}
	if err := json.NewDecoder(req.Body).Decode(innerReq); err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		http.Error(rw, `{"message": "invalid request"}`, http.StatusBadRequest)

		return
	}

	resp, err := NewProxy(
		h.cfg.ClusterAddress,
		&http.Client{
			Transport: transport,
		},
	).Run(innerReq)

	if err != nil {
		http.Error(rw, `{"message": "internal server error"}`, http.StatusInternalServerError)
		log.Printf("internal server error: %+v", err)

		return
	}
	defer resp.Body.Close()

	for k, vals := range resp.Header {
		for _, v := range vals {
			rw.Header().Add(k, v)
		}
	}

	rw.WriteHeader(resp.StatusCode)
	io.Copy(rw, resp.Body)
}

func (h *handler) loginHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Generate random state
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	state := base64.StdEncoding.EncodeToString(b)

	session, err := h.store.Get(r, cookieKey)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	session.Values["state"] = state
	err = session.Save(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	authenticator, err := NewAuthenticator(ctx, h.cfg)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, authenticator.Config.AuthCodeURL(state), http.StatusTemporaryRedirect)
}

func (h *handler) callbackHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	session, err := h.store.Get(r, cookieKey)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Println(getParam(session, "state"), r.URL.Query().Get("state"))

	if r.URL.Query().Get("state") != getParam(session, "state") {
		http.Error(w, "Invalid state parameter", http.StatusBadRequest)
		return
	}

	authenticator, err := NewAuthenticator(ctx, h.cfg)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	token, err := authenticator.Config.Exchange(ctx, r.URL.Query().Get("code"))
	if err != nil {
		log.Printf("no token found: %v", err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	rawIDToken, ok := token.Extra("id_token").(string)
	if !ok {
		http.Error(w, "No id_token field in oauth2 token.", http.StatusInternalServerError)
		return
	}

	oidcConfig := &oidc.Config{
		ClientID: h.cfg.ClientID,
	}

	idToken, err := authenticator.Provider.Verifier(oidcConfig).Verify(ctx, rawIDToken)

	authenticator.Provider.Verifier(oidcConfig)

	if err != nil {
		http.Error(w, "Failed to verify ID Token: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Getting now the userInfo
	var profile map[string]interface{}
	if err := idToken.Claims(&profile); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	session.Values["id_token"] = rawIDToken
	session.Values["access_token"] = token.AccessToken
	session.Values["refresh_token"] = token.RefreshToken
	session.Values["profile"] = profile
	err = session.Save(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Redirect to logged in page
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
