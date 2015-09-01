package main // github.com/gitinsky/vnc-go-web

import (
	"net/http"
    "time"
)

type LogonHandler struct {
}

func (h LogonHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    p := Responder{w, r, time.Now()}
    p.serveLogin()
}

