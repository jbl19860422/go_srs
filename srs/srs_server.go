package srs

import (
	log "github.com/sirupsen/logrus"
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
		HandShaker:SrsHandshakeBytes{
			C0C1:make([]byte, 1537),
			S0S1S2:make([]byte, 3073),
			C2:make([]byte, 1536),
		},
	}
	this.srsServers = append(this.srsServers, rtmpServer)
	log.Info("star a new server")
	go rtmpServer.Start()
}


