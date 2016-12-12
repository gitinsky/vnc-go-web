package main // github.com/gitinsky/vnc-go-web

import (
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"net/http"
	"os"
	"time"
)

var slidingPassword = NewSlidingPassword()

func main() {
	cfg.Parse()
	spew.Fprintf(os.Stderr, "started: %#v\n", cfg)

	if *cfg.ssl {
		panic(fmt.Errorf("SSL not implemented yet"))
	}

	go slidingPassword.UpdateLoop(time.Duration(*cfg.authTTL/2*1000) * time.Millisecond)

	http.Handle(*cfg.baseuri, RootGetPost{http.StripPrefix(*cfg.baseuri, http.FileServer(http.Dir(*cfg.root)))})
	http.Handle(*cfg.baseuri+"logon", LogonHandler{})
	http.Handle(*cfg.baseuri+"retry", RetryHandler{})
	http.Handle(*cfg.baseuri+"websockify", NewWssHandler())

	panic(http.ListenAndServe(*cfg.listen, nil))
}
