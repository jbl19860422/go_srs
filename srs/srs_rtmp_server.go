package srs

import (
	// "fmt"
	"go_srs/srs/protocol"
	log "github.com/sirupsen/logrus"
)


type SrsRtmpServer struct {
	Conn *SrsRtmpConn
	HandShaker protocol.SrsHandshakeBytes
}

func (this *SrsRtmpServer) HandShake() int {
	this.HandShaker.ReadC0C1(&(this.Conn.Conn))
	return 0
}

func (this *SrsRtmpServer) Start() int {
	log.Info("start rtmp server")
	ret := this.HandShake()
	if ret != 0 {

	}
	return 0
}