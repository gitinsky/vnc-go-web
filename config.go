package main // github.com/gitinsky/vnc-go-web

import (
	"flag"
	"fmt"
	"os"
)

type Config struct {
	// http server specific parameters
	listen *string
	ssl    *bool
    root   *string
    auth   *string
    resolv *string
}

var cfg = Config{
    listen: flag.String("listen", ":8080", "Address to HTTP(S) listen. [ADDR]:PORT"),
	ssl:    flag.Bool("ssl", false, "Use SSL"),
	root:   flag.String("root", "./", "Document root for static pages"),
    auth:   flag.String("auth", "https://contentdelivery.mf-master.ru:8443/api/?method=get_uuid&server_name=", "External authentication URL"),
    resolv: flag.String("resolv", "http://172.23.0.1:8080/cgi-bin/uuid2ip.cgi?", "External ID to IP resolving URL"),
}

func (*Config) Parse() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage:\n")
		flag.PrintDefaults()
	}

	flag.Parse()
}

type strList []string

func (s *strList) String() string {
	return fmt.Sprintf("%s", *s)
}

func (s *strList) Set(v string) error {
	*s = append(*s, v)
	return nil
}
