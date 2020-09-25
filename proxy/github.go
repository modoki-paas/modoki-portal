package main

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"

	"github.com/google/go-github/v32/github"
	"golang.org/x/oauth2"
)

var (
	client http.Client
)

func (h *handler) parseGitHubURL(rw http.ResponseWriter, req *http.Request) (*url.URL, bool) {
	ghurl, err := url.Parse(h.cfg.GitHub.URL)

	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)

		return nil, false
	}

	return ghurl, true
}

func (h *handler) loginGitHub(rw http.ResponseWriter, req *http.Request) {
	state := randomState()

	session := h.loadSession(rw, req)
	if session == nil {
		return
	}

	session.Values["github_state"] = state

	if !h.saveSession(session, rw, req) {
		return
	}

	http.Redirect(rw, req, h.github.AuthCodeURL(state), http.StatusSeeOther)
}

func (h *handler) callbackGitHub(rw http.ResponseWriter, req *http.Request) {
	session := h.loadSession(rw, req)
	if session == nil {
		return
	}

	istate, ok := session.Values["github_state"]
	if !ok {
		http.Redirect(rw, req, "/github/login", http.StatusSeeOther)

		return
	}

	expectedState := istate.(string)

	code := req.URL.Query().Get("code")
	state := req.URL.Query().Get("state")

	if expectedState != state {
		http.Error(rw, `{"message": "invalid state"}`, http.StatusBadRequest)

		return
	}

	token, err := h.github.Exchange(req.Context(), code)

	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)

		return
	}

	session.Values["github_token"] = token

	if !h.saveSession(session, rw, req) {
		return
	}

	http.Redirect(rw, req, "/", http.StatusSeeOther)
}

func (h *handler) githubAPI(fn func(ctx context.Context, gh *github.Client) (interface{}, error)) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		ctx := req.Context()

		session := h.loadSession(rw, req)
		if session == nil {
			return
		}
		gh, ok := session.Values["github_token"]
		if !ok {
			http.Redirect(rw, req, "/github/login", http.StatusSeeOther)

			return
		}
		token := gh.(*oauth2.Token)

		ts := h.github.TokenSource(ctx, token)

		client := github.NewClient(oauth2.NewClient(ctx, ts))

		resp, err := fn(ctx, client)

		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)

			return
		}

		token, err = ts.Token()

		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)

			return
		}

		session.Values["github_token"] = token

		if !h.saveSession(session, rw, req) {
			return
		}

		rw.Header().Set("Content-Type", "application/json")
		json.NewEncoder(rw).Encode(resp)
	})
}

func (h *handler) listInstallations(rw http.ResponseWriter, req *http.Request) {
	h.githubAPI(func(ctx context.Context, gh *github.Client) (interface{}, error) {
		resp, _, err := gh.Apps.ListUserInstallations(ctx, nil)

		if err != nil {
			return nil, err
		}

		return resp, nil
	}).ServeHTTP(rw, req)
}

func (h *handler) listRepositories(rw http.ResponseWriter, req *http.Request) {
	h.githubAPI(func(ctx context.Context, gh *github.Client) (interface{}, error) {
		ins := req.URL.Query().Get("installation_id")

		insID, err := strconv.ParseInt(ins, 10, 64)

		if err != nil {
			return nil, err
		}

		resp, _, err := gh.Apps.ListUserRepos(ctx, insID, nil)

		if err != nil {
			return nil, err
		}

		return resp, nil
	}).ServeHTTP(rw, req)
}
