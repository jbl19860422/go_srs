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
	ret := this.HandShaker.ReadC0C1(&(this.Conn.Conn))
	if 0 != ret {
		log.Error("HandShake ReadC0C1 failed")
		return -1
	}

	if this.HandShaker.C0C1[0] != 0x03 {
		log.Info("only support rtmp plain text.")
		return -2
	}

	if 0 != this.HandShaker.CreateS0S1S2() {
		return -2
	}

	n, err := this.Conn.Conn.Write(this.HandShaker.S0S1S2)
	if err != nil {
		log.Error("write s0s1s2 failed")
	} else {
		log.Info("write s0s1s2 succeed, count=", len(this.HandShaker.S0S1S2))
	}

	if 0 != this.HandShaker.ReadC2(&(this.Conn.Conn)) {
		log.Error("HandShake ReadC2 failed")
		return -3
	}

	if !this.HandShaker.CheckC2() {
		log.Error("HandShake CheckC2 failed")
	}

	log.Info("HandShake Succeed")
	_ = n

	return 0
}

func (this *SrsRtmpServer) Start() int {
	log.Info("start rtmp server")
	ret := this.HandShake()
	if ret != 0 {
		log.Info("HandShake failed")
		return -1
	}

	
	return 0
}