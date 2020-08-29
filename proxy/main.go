package main

import "net/http"

type Request struct {
	Method  string            `json:"method"`
	URL     string            `json:"url"`
	Headers map[string]string `json:"headers"`
	Body    []byte            `json:"body"`
}

func main() {
	http.HandleFunc("/proxy", func(rw http.ResponseWriter, req *http.Request) {

	})

	http.ListenAndServe(":80", nil)
}
