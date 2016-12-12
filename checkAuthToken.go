package main // github.com/gitinsky/vnc-go-web

import "encoding/base64"

type AuthToken struct {
	Peer  string `json:"peer,omitempty"`
	Real  string `json:"real,omitempty"`
	Dest  string `json:"dest,omitempty"`
	Retry string `json:"retry,omitempty"`
}

func (p *Responder) CheckAuthToken() *AuthToken {
	if p.r.URL.RawQuery != "" {
		authStr, _ := base64.URLEncoding.DecodeString(p.r.URL.RawQuery)
		if authStr != nil {
			authParams := slidingPassword.Decrypt(authStr)
			if authParams != nil {
				if authParams.Peer != p.r.RemoteAddr || authParams.Real != p.getHeader("X-Real-IP", "UNKNOWN") {
					return authParams
				}
			}
		}
	}
	return nil
}
