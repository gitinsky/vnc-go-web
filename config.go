package main // github.com/gitinsky/vnc-go-web

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

const (
	authCookieName = "vnc-go-web"
)

type Config struct {
	// http server specific parameters
	listen  *string
	ssl     *bool
	root    *string
	passwd  *string
	auth    *string
	resolv  *string
	authTTL *float64
	baseuri *string

	// VNC server specific parameters
	vnc_port         *int
	vnc_true_color   *bool
	vnc_local_cursor *bool
	vnc_shared       *bool
	vnc_view_only    *bool
}

var cfg = Config{
	listen:  flag.String("listen", ":8080", "Address to HTTP(S) listen. [ADDR]:PORT"),
	ssl:     flag.Bool("ssl", false, "Use SSL"),
	root:    flag.String("root", "./", "Document root for static pages"),
	passwd:  flag.String("passwd", "", "Password file name"),
	auth:    flag.String("auth", "http://127.0.0.1/auth?", "External authentication URL"),
	resolv:  flag.String("resolv", "http://127.0.0.1/resolv?", "External ID to IP resolving URL"),
	authTTL: flag.Float64("authTTL", 10, "Authentication token TTL in seconds"),
	baseuri: flag.String("baseuri", "", "Base URI"),

	vnc_port:         flag.Int("vnc_port", 5900, "VNC port to connect"),
	vnc_true_color:   flag.Bool("vnc_true_color", false, "true_color noVNC param"),
	vnc_local_cursor: flag.Bool("vnc_local_cursor", true, "cursor noVNC param"),
	vnc_shared:       flag.Bool("vnc_shared", true, "shared noVNC param"),
	vnc_view_only:    flag.Bool("vnc_view_only", false, "view_only noVNC param"),
}

func (cfg *Config) Parse() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage:\n")
		flag.PrintDefaults()
	}

	flag.Parse()

	str := strings.Trim(*cfg.baseuri, " \t\r\n/")
	str = "/" + str + ternary(len(str) > 0, "/", "").(string)
	cfg.baseuri = &str
}

type strList []string

func (s *strList) String() string {
	return fmt.Sprintf("%s", *s)
}

func (s *strList) Set(v string) error {
	*s = append(*s, v)
	return nil
}
