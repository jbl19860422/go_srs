package srs

import (
	log "github.com/sirupsen/logrus"
)


type SrsRtmpServer struct {
	Conn *SrsRtmpConn
	HandShaker SrsHandshakeBytes
}

func (this *SrsRtmpServer) HandShake() int {
	this.HandShaker.ReadC0C1(this.Conn)
	return 0
}

func (this *SrsRtmpServer) Start() int {
	log.Info("start rtmp server")
	this.HandShake()
	return 0
}