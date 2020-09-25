package main

import (
	"crypto/rand"
	"encoding/base64"
	"io"
	"io/ioutil"

	"github.com/gorilla/sessions"
)

func getParam(sess *sessions.Session, key string) string {
	v, ok := sess.Values[key]

	if !ok {
		return ""
	}

	str, ok := v.(string)

	if !ok {
		return ""
	}

	return str
}

func randomState() string {
	b, err := ioutil.ReadAll(io.LimitReader(rand.Reader, 4))

	if err != nil {
		panic(err)
	}

	return base64.StdEncoding.EncodeToString(b)
}
