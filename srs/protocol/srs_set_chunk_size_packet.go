package protocol

type SrsSetChunkSizePacket struct {
	SrsPacket
	/**
	 * The maximum chunk size can be 65536 bytes. The chunk size is
	 * maintained independently for each direction.
	 */
	chunk_size int32
}

func NewSrsSetChunkSizePacket() *SrsSetChunkSizePacket {
	return &SrsSetChunkSizePacket{
		chunk_size: SRS_CONSTS_RTMP_PROTOCOL_CHUNK_SIZE,
	}
}

func (s *SrsSetChunkSizePacket) get_message_type() int8 {
	return RTMP_MSG_AMF0CommandMessage
}

func (s *SrsSetChunkSizePacket) get_prefer_cid() int32 {
	return RTMP_CID_ProtocolControl
}

func (p *SrsSetChunkSizePacket) decode(s *SrsStream) error {
	var err error
	p.chunk_size, err = s.read_int32()
	return err
}

func (s *SrsSetChunkSizePacket) encode() ([]byte, error) {
	return nil, nil
}
