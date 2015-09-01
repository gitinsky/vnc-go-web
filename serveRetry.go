package main // github.com/gitinsky/vnc-go-web

import (
	"net/http"
    "time"
)

type RetryHandler struct {
}

func (h RetryHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    p := Responder{w, r, time.Now()}
    
    dest, serverID := p.CheckAuthToken()
    if dest != "retry" || serverID == "" {
        http.Redirect(w, r, *cfg.baseuri+"error.html", http.StatusFound)
        p.errorLog(http.StatusFound, "auth token invalid")
        return
    }

    serverIP, err := getServerIP(serverID)
    
    if err != nil {
        http.Redirect(p.w, p.r, *cfg.baseuri+"panic.html", http.StatusFound)
        p.errorLog(http.StatusFound, "error resolving '%s' ('%s', '%s'): %s", "-", serverID, serverIP, err.Error())
        return
    }
    
    if serverIP == "" {
        http.Redirect(p.w, p.r, *cfg.baseuri+"retry.html?"+p.GetAuthToken("retry", serverID), http.StatusFound)
        p.errorLog(http.StatusFound, "error resolving '%s' ('%s', '%s'): server offline", "-", serverID, serverIP)
        return
    }
    
    http.Redirect(p.w, p.r, p.GetVncUrl(serverIP), http.StatusFound)
    p.errorLog(http.StatusFound, "resolved '%s' ('%s', '%s')", "-", serverID, serverIP)
}
