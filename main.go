package main

import(
	"time"
	"os"
	"flag"
	log "github.com/sirupsen/logrus"
	"go_srs/srs"
)

var (
	port = flag.Int("p", 1935, "set port `port`")
)

func init() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
}

func main() {
	flag.Parse()
	//init server
	l := &srs.SrsStreamListener{}
	server := &srs.SrsServer{Listener:l}
	l.Svr = server
	server.StartProcess(*port)
	time.Sleep(1*time.Second)
}