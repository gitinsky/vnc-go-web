package main // github.com/gitinsky/vnc-go-web

import (
	"net/http"
    "github.com/davecgh/go-spew/spew"
    "os"
)

func main() {
    cfg.Parse()
	spew.Fprintf(os.Stderr, "started: %#v\n", cfg)

    http.Handle("/", RootGetPost{http.StripPrefix("/", http.FileServer(http.Dir(*cfg.root)))})

	panic(http.ListenAndServe(*cfg.listen, nil))
}
