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
    p := Responder{w, r, time.Now()}
    switch r.Method {
        case "GET":
            h.fs.ServeHTTP(w, r)
            p.accessLog(http.StatusOK)
        default:
            http.Error(w, r.Method+" not allowed", http.StatusMethodNotAllowed)
            p.errorLog(http.StatusMethodNotAllowed, "Method '%s' not allowed", r.Method)
    }
}

