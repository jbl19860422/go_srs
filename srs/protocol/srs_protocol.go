package protocol

type SrsProtocol struct {
	chunkCache []*SrsChunkMessage
}

func (s *SrsProtocol) RecvInterlacedMessage(conn *net.Conn) int {
	chunk := NewSrsChunkMessage()
	ret := chunk.ReadBasicHeader(conn)
	if 0 != ret {
		return -1
	}


}

func (s *SrsProtocol) RecvMessage() int {
	
}


type SrsChunkMessage struct {
	Buffer 	[]byte 
	Pos		int
}

func NewSrsChunkMessage() *SrsChunkMessage {
	return &SrsChunkMessage{Buffer:make([]byte, 148),Pos:0}
}

func (s *SrsChunkMessage) ReadNByte(conn *net.Conn, count int) int {
	b := s.Buffer[s.Pos:s.Pos+count+1]
	n, err := conn.Read(b)
	if err != nil {
		return -1
	}
	s.Pos += n
	return 0
}

func (s *SrsChunkMessage) ReadBasicHeader(conn *net.Conn) (fmt byte, cid int32, ret int) {
	if 0 != s.ReadNByte(conn, 1) {
		ret = -1
		return
	}

	cid = this.Buffer[0]&0x3f
	fmt = (this.Buffer[0]>>6) & 0x3
	// 2-63, 1B chunk header
	if cid > 1 {
		ret = -2
		return
	}
	// 64-319, 2B chunk header
	if cid == 0 {
		if 0 != s.ReadNByte(conn, 1) {
			ret = -3
			return
		}

		cid = 64
		cid += s.Buffer[1]
	} else if cid == 1 {// 64-65599, 3B chunk header
		if 0 != s.ReadNByte(conn, 2) {
			ret = -4
			return
		}

		cid = 64
		cid += s.Buffer[1]
		cid += s.Buffer[2]
		return
	}
	ret = -5
	return
}




