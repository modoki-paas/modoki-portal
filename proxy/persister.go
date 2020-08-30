package main

import (
	"net/http"

	"github.com/gorilla/sessions"
	restclient "k8s.io/client-go/rest"
)

type CookieConfigPersister struct {
	session *sessions.Session
	rw      http.ResponseWriter
	req     *http.Request
}

func NewCookieConfigPersister(session *sessions.Session, rw http.ResponseWriter, req *http.Request) restclient.AuthProviderConfigPersister {
	return &CookieConfigPersister{
		session: session,
		rw:      rw,
		req:     req,
	}
}

var _ restclient.AuthProviderConfigPersister = &CookieConfigPersister{}

func (p *CookieConfigPersister) Persist(param map[string]string) error {
	session := p.session

	session.Values["id_token"] = param["id-token"]
	session.Values["refresh_token"] = param["refresh-token"]

	session.Save(p.req, p.rw)

	return nil
}
