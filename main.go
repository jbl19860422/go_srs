package main

import(
	"time"
	"os"
	log "github.com/sirupsen/logrus"
	"go_srs/srs"
)

func init() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
}

func main() {
	//init server
	l := &srs.SrsStreamListener{}
	server := &srs.SrsServer{Listener:l}
	l.Svr = server
	server.StartProcess()
	time.Sleep(1*time.Second)
}