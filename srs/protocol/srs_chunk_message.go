package protocol

import (
	"net"
	log "github.com/sirupsen/logrus"
	"encoding/binary"
	"bytes"
)

const (
	RTMP_FMT_TYPE0 = 0
	RTMP_FMT_TYPE1 = 1
	RTMP_FMT_TYPE2 = 2
	RTMP_FMT_TYPE3 = 3
)

const RTMP_EXTENDED_TIMESTAMP = 0xFFFFFF

type SrsChunkMessage struct {
	Buffer 	[]byte 
	Pos		int
	RtmpMessage *SrsRtmpMessage
}

var mh_sizes = [4]int{11, 7, 3, 0}

func NewSrsChunkMessage() *SrsChunkMessage {
	return &SrsChunkMessage{Buffer:make([]byte, 148),Pos:0}
}

func (s *SrsChunkMessage) ReadNByte(conn *net.Conn, count int) (b []byte, err error) {
	b = make([]byte, count)
	_, err = (*conn).Read(b)
	return
}

func (s *SrsChunkMessage) ReadBasicHeader(conn *net.Conn) (fmt byte, cid int32, err error) {
	var buffer1 []byte
	var buffer2 []byte
	var buffer3 []byte
	if buffer1, err = s.ReadNByte(conn, 1); err != nil {
		return
	}

	cid = (int32)(buffer1[0]&0x3f)
	fmt = (buffer1[0]>>6) & 0x3
	// 2-63, 1B chunk header
	if cid > 1 {
		return
	}
	// 64-319, 2B chunk header
	if cid == 0 {
		if buffer2, err = s.ReadNByte(conn, 1); err != nil {
			return
		}

		cid = 64
		cid += (int32)(buffer2[0])
	} else if cid == 1 {// 64-65599, 3B chunk header
		if buffer3, err = s.ReadNByte(conn, 2); err != nil {
			return
		}

		cid = 64
		cid += (int32)(buffer3[0])
		cid += (int32)(buffer3[1])
		return
	}
	return
}

func (s *SrsChunkMessage) ReadMessageHeader(conn *net.Conn, fmt byte, cid int32) (err error) {
	var buf1 []byte
	if fmt == RTMP_FMT_TYPE0 {
		if s.RtmpMessage == nil {
			s.RtmpMessage = NewSrsRtmpMessage()
		}
	}

	var mh_size = mh_sizes[fmt]
	if mh_size > 0 {
		if buf1, err = s.ReadNByte(conn, mh_size); err != nil {
			log.Error("read message header failed")
		}
	}

	if fmt <= RTMP_FMT_TYPE2 {
		buf_timestamp := make([]byte, 4)
		buf_timestamp[2] = buf1[0]
		buf_timestamp[1] = buf1[1]
		buf_timestamp[0] = buf1[2]
		buf_timestamp[3] = 0

		var timestamp int32
		buf_reader := bytes.NewBuffer(buf_timestamp)
		binary.Read(buf_reader, binary.LittleEndian, &timestamp)

		extend_timestamp := false
		if timestamp > RTMP_EXTENDED_TIMESTAMP {
			extend_timestamp = true
		}
		//这里需用用chunk记录一个chunk开头的时间
		// if !extend_timestamp {
		// 	if fmt == RTMP_FMT_TYPE0 {
		// 		timestamp1 := timestamp
		// 	} else {
		// 		timestamp2 := timestamp1 +  timestamp
		// 	}
		// }
		log.Info("timestamp=", timestamp)
	}
	_ = buf1
	return
}
