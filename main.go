package main // github.com/gitinsky/vnc-go-web

import (
	"net/http"
    "github.com/davecgh/go-spew/spew"
    "os"
    "time"
)

var slidingPassword = NewSlidingPassword()

func main() {
    cfg.Parse()
	spew.Fprintf(os.Stderr, "started: %#v\n", cfg)

    go slidingPassword.UpdateLoop(time.Duration(*cfg.authTTL/2*1000) * time.Millisecond)
    
    http.Handle("/", RootGetPost{http.StripPrefix("/", http.FileServer(http.Dir(*cfg.root)))})
    http.Handle("/retry", RetryHandler{})
    http.Handle("/websockify", NewWssHandler())

	panic(http.ListenAndServe(*cfg.listen, nil))
}
