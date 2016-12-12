package main // github.com/gitinsky/vnc-go-web

import (
	"fmt"
	"golang.org/x/net/websocket"
	"io"
	"net"
	"net/http"
	"sync"
	"time"
)

type WssVncRequest struct {
}

func NewWssHandler() websocket.Server {
	return websocket.Server{Handshake: bootHandshake, Handler: handleWss}
}

func (p *Responder) CopyR2W(t *copySyncer, dst io.Writer, src io.Reader, descr string) {
	defer t.KillIt()

	buf := make([]byte, 1000)

	for t.IsAlive() {
		rn, err := src.Read(buf)
		if err != nil {
			return
		}
		_, err = dst.Write(buf[:rn])
		if err != nil {
			return
		}
	}
}

type copySyncer struct {
	alive bool
	lock  sync.RWMutex
}

func (t *copySyncer) IsAlive() bool {
	defer t.lock.RUnlock()
	t.lock.RLock()
	return t.alive
}

func (t *copySyncer) KillIt() {
	defer t.lock.Unlock()
	t.lock.Lock()
	t.alive = false
}

func handleWss(wsconn *websocket.Conn) {
	p := Responder{nil, wsconn.Request(), time.Now()}
	serverIP := wsconn.Request().Header.Get("X-Server-IP")

	conn, err := net.Dial("tcp", serverIP)
	if err != nil {
		p.errorLog(http.StatusInternalServerError, "Error connecting to '%s': %s", serverIP, err.Error())
		wsconn.Close()
		return
	}
	defer conn.Close()
	defer wsconn.Close()

	wsconn.PayloadType = websocket.BinaryFrame

	t := &copySyncer{alive: true}

	go p.CopyR2W(t, conn, wsconn, serverIP+" ws2vnc")
	go p.CopyR2W(t, wsconn, conn, serverIP+" vnc2ws")

	p.errorLog(http.StatusOK, "websocket started: '%s'", serverIP)

	for t.IsAlive() {
		time.Sleep(100 * time.Millisecond)
	}

	p.errorLog(http.StatusOK, "websocket closed: '%s'", serverIP)
}

func bootHandshake(config *websocket.Config, r *http.Request) error {
	p := Responder{nil, r, time.Now()}

	dest, serverIP := p.CheckAuthToken()
	if dest != "dest" || serverIP == "" {
		p.errorLog(http.StatusForbidden, "auth token invalid")
		return fmt.Errorf("auth token invalid")
	}

	config.Protocol = []string{"binary"}

	r.Header.Set("X-Server-IP", serverIP)
	r.Header.Set("Access-Control-Allow-Origin", "*")
	r.Header.Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE")

	p.accessLog(http.StatusSwitchingProtocols)

	return nil
}
