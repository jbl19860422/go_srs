package main

import (
	"time"
	// "os"
	"flag"
	// log "github.com/sirupsen/logrus"
	"go_srs/srs"
	// "log"
)

var (
	port = flag.Int("p", 1935, "set port `port`")
)

func main() {
	flag.Parse()
	//init server
	listener := &srs.SrsStreamListener{}
	server := &srs.SrsServer{Listener: l}
	l.Svr = server
	server.StartProcess(*port)
}
