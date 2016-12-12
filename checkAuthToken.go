package main // github.com/gitinsky/vnc-go-web

import (
	"encoding/base64"
	"regexp"
)

var checkAuthTokenRegexp = regexp.MustCompile("^peer\\s+([^\\s]+)\\s+real\\s+([^\\s]+)\\s+((?:dest)|(?:retry))\\s+([^\\s]+)$")

func (p *Responder) CheckAuthToken() (string, string) {
	if p.r.URL.RawQuery != "" {
		authStr, _ := base64.URLEncoding.DecodeString(p.r.URL.RawQuery)
		if authStr != nil {
			authParams := slidingPassword.Decrypt(authStr, checkAuthTokenRegexp)
			if authParams != nil {
				if authParams[1] != p.r.RemoteAddr || authParams[2] != p.getHeader("X-Real-IP", "UNKNOWN") {
					return authParams[3], authParams[4]
				}
			}
		}
	}
	return "", ""
}
