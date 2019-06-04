package protocol

import (
	"net"
	log "github.com/sirupsen/logrus"
)

type SrsProtocol struct {
	chunkCache []*SrsChunkMessage
}

func (s *SrsProtocol) RecvInterlacedMessage(conn *net.Conn) int {
	chunk := NewSrsChunkMessage()
	fmt, cid, err := chunk.ReadBasicHeader(conn)
	if nil != err {
		log.Errorf("read basic header failed, err=%v", err)
		return -1
	}

	log.Info("read Basic header succeed, cid=", cid, ", fmt=", fmt)
	err = chunk.ReadMessageHeader(conn, fmt, cid)
	if err != nil {
		
	}

	return 0
}

func (s *SrsProtocol) RecvMessage() int {
	return 0
}




