package main

import "github.com/gorilla/sessions"

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
