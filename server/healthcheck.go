package main

import (
	"io"
	"io/ioutil"
	"net/http"
)

type HealthcheckHandler struct {
}

func (h HealthcheckHandler) ServeHTTP(w http.ResponseWriter, request *http.Request) {
	defer request.Body.Close()

	io.Copy(ioutil.Discard, request.Body)

	w.WriteHeader(200)
}
