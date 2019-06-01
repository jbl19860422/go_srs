package srs

import (
	log "github.com/sirupsen/logrus"
	"go_srs/srs/protocol"
	// "fmt"
)

const (
	RTMP_PORT = 1935
)

type SrsServer struct {
	streams []SrsStream	
	srsServers []*SrsRtmpServer
	Listener *SrsStreamListener
}

func (this *SrsServer) StartProcess() {
	this.Listener.ListenAndAccept()
}

func (this *SrsServer) AcceptConnection(c *SrsRtmpConn) {
	rtmpServer := &SrsRtmpServer{
		Conn:c,
		HandShaker:protocol.SrsHandshakeBytes{},
	}
	this.srsServers = append(this.srsServers, rtmpServer)
	log.Info("star a new server")
	go rtmpServer.Start()
}


