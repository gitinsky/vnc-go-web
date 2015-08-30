package main // github.com/gitinsky/vnc-go-web

import (
	"net/http"
    "net/url"
    "strings"
    "crypto/tls"
    "io"
    "io/ioutil"
)

func (p *Responder) serveLogin() {
	p.r.ParseForm()

    serverNum := strings.Trim(p.r.PostFormValue("servernum"), " \r\n")
    if serverNum == "" {
        http.Redirect(p.w, p.r, "/error.html", http.StatusFound)
        p.errorLog(http.StatusFound, "empty input")
        return
    }
    
    serverID, serverIP, err := getServerIP(serverNum)
    
    if err != nil {
        http.Redirect(p.w, p.r, "/panic.html", http.StatusFound)
        p.errorLog(http.StatusFound, "error resolving '%s' ('%s', '%s'): %s", serverNum, serverID, serverIP, err.Error())
        return
    }
    
    if serverIP == "" {
        http.Redirect(p.w, p.r, "/error.html", http.StatusFound)
        p.errorLog(http.StatusFound, "error resolving '%s' ('%s', '%s'): server offline", serverNum, serverID, serverIP)
        return
    }
    
    http.Redirect(p.w, p.r, "/ok.html", http.StatusFound)
    p.errorLog(http.StatusFound, "resolved '%s' ('%s', '%s')", serverNum, serverID, serverIP)
}


func getServerID(serverNum string) (string, error) {
    client := &http.Client{Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}}
    
    resp, err := client.Get(*cfg.auth+url.QueryEscape(serverNum))
    if err != nil {
        return "", err
    }
    defer resp.Body.Close()

    buf, err := ioutil.ReadAll(io.LimitReader(resp.Body, 256))
    if err != nil && err != io.EOF {
        return "", err
    }

    return strings.Trim(string(buf), " \r\n"), nil
}

func getServerIP(serverNum string) (string, string, error) {
    serverID, err := getServerID(serverNum)
    if serverID == "" {
        return "", "", err
    }
    
    resp, err := http.Get(*cfg.resolv+url.QueryEscape(serverID))
    if err != nil {
        return serverID, "", err
    }
    defer resp.Body.Close()

    buf, err := ioutil.ReadAll(io.LimitReader(resp.Body, 256))
    if err != nil && err != io.EOF {
        return serverID, "", err
    }

    return serverID,  strings.Trim(string(buf), " \r\n"), nil
}




