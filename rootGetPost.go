package main // github.com/gitinsky/vnc-go-web

import (
	"net/http"
    "time"
)

type Responder struct {
    w     http.ResponseWriter
    r     *http.Request
    sTime time.Time
}

func (p *Responder) getHeader(name string, def string) string {
	h := p.r.Header.Get(name)
	if h == "" {
		return def
	}
	return h
}

type RootGetPost struct {
    fs http.Handler
}

func (h RootGetPost) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    resp := Responder{w, r, time.Now()}
    switch r.Method {
        case "GET":
            h.fs.ServeHTTP(w, r)
            resp.accessLog(http.StatusOK)
        case "POST":
            resp.serveLogin()
        default:
            http.Error(w, "GET or POST!", http.StatusMethodNotAllowed)
            resp.errorLog(http.StatusMethodNotAllowed, "Method not allowed")
    }
}

