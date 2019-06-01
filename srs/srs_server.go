package srs

import (
	// log "github.com/sirupsen/logrus"
	// "go_srs/srs/protocol"
	// "fmt"
	"log"
)

type SrsServer struct {
	streams    []SrsStream
	srsServers []*SrsRtmpServer
	Listener   *SrsStreamListener
}

func (this *SrsServer) StartProcess(port int) {
	this.Listener.ListenAndAccept(port)
}

func (this *SrsServer) AcceptConnection(c *SrsRtmpConn) {
	rtmpServer := NewSrsRtmpServer(c)
	this.srsServers = append(this.srsServers, rtmpServer)
	log.Printf("star a new server")
	go rtmpServer.Start()
}
