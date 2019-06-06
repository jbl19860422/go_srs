package srs

import (
	// "fmt"
	_ "context"
	"go_srs/srs/protocol"
	"log"
	"time"
	// log "github.com/sirupsen/logrus"
)

type SrsRtmpServer struct {
	Conn       *SrsRtmpConn
	Protocol   *protocol.SrsProtocol
	HandShaker protocol.SrsHandshakeBytes
}

func NewSrsRtmpServer(conn *SrsRtmpConn) *SrsRtmpServer {
	return &SrsRtmpServer{Conn: conn, Protocol: protocol.NewSrsProtocol(), HandShaker: protocol.SrsHandshakeBytes{}}
}

func (this *SrsRtmpServer) HandShake() int {
	ret := this.HandShaker.ReadC0C1(&(this.Conn.Conn))
	if 0 != ret {
		log.Printf("HandShake ReadC0C1 failed")
		return -1
	}

	if this.HandShaker.C0C1[0] != 0x03 {
		log.Printf("only support rtmp plain text.")
		return -2
	}

	if 0 != this.HandShaker.CreateS0S1S2() {
		return -2
	}

	n, err := this.Conn.Conn.Write(this.HandShaker.S0S1S2)
	if err != nil {
		log.Printf("write s0s1s2 failed")
	} else {
		log.Printf("write s0s1s2 succeed, count=", len(this.HandShaker.S0S1S2))
	}

	if 0 != this.HandShaker.ReadC2(&(this.Conn.Conn)) {
		log.Printf("HandShake ReadC2 failed")
		return -3
	}

	if !this.HandShaker.CheckC2() {
		log.Printf("HandShake CheckC2 failed")
	}

	log.Printf("HandShake Succeed")
	_ = n

	return 0
}

func (this *SrsRtmpServer) Start() int {
	log.Printf("start rtmp server")
	ret := this.HandShake()
	if ret != 0 {
		log.Printf("HandShake failed")
		return -1
	}

	// ctx, cancel := context.WithCancel(context.Background())
	// go this.Protocol.LoopMessage(ctx, &(this.Conn.Conn))
	connPacket := protocol.NewSrsConnectAppPacket()
	this.Protocol.ExpectMessage(&(this.Conn.Conn), connPacket)
	for {
		time.Sleep(10*time.Second)
	}

	// _ = cancel
	
	// msg, err := this.Protocol.RecvInterlacedMessage(&(this.Conn.Conn))
	// if err != nil {
	// 	log.Print("RecvInterlacedMessage err=", err)
	// 	return -2
	// }

	// _ = msg

	return 0
}
