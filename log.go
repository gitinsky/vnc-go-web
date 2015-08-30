package main

import (
	"fmt"
	"os"
	"time"
)

func (p *Responder) accessLog(status int) {
	//timestamp spentTime peer x-real-ip method status 'request URI'
	fmt.Fprintf(os.Stderr, "%s %d %s %s %s %3.3d '%s'\n",
		p.sTime.Local().Format("2006-01-02-15-04-05.000"),
		int(time.Now().Sub(p.sTime).Seconds()*1000),
		p.r.RemoteAddr,
		p.getHeader("X-Real-IP", "UNKNOWN"),
		p.r.Method,
		status,
		p.r.URL.Path,
	)
}

func (p *Responder) errorLog(status int, msg string, params ...interface{}) {
	//timestamp spentTime peer x-real-ip method status 'request URI' message
	fmt.Fprintf(os.Stderr, "%s %d %s %s %s %3.3d '%s' %s\n",
		p.sTime.Local().Format("2006-01-02-15-04-05.000"),
		int(time.Now().Sub(p.sTime).Seconds()*1000),
		p.r.RemoteAddr,
		p.getHeader("X-Real-IP", "UNKNOWN"),
		p.r.Method,
		status,
		p.r.URL.Path,
		fmt.Sprintf(msg, params...),
	)
}

func infoLog(msg string, params ...interface{}) {
	//timestamp spentTime peer x-real-ip method status 'request URI' message
	fmt.Fprintf(os.Stderr, "%s %s\n",
		time.Now().Local().Format("2006-01-02-15-04-05.000"),
		fmt.Sprintf(msg, params...),
	)
}
