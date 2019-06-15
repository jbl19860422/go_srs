package packet

type SrsSetChunkSizePacket struct {
	/**
	 * The maximum chunk size can be 65536 bytes. The chunk size is
	 * maintained independently for each direction.
	 */
	ChunkSize int32
}

func NewSrsSetChunkSizePacket() *SrsSetChunkSizePacket {
	return &SrsSetChunkSizePacket{
		ChunkSize: SRS_CONSTS_RTMP_PROTOCOL_CHUNK_SIZE,
	}
}

func (this *SrsSetChunkSizePacket) GetMessageType() int8 {
	return RTMP_MSG_SetChunkSize
}

func (this *SrsSetChunkSizePacket) GetPreferCid() int32 {
	return RTMP_CID_ProtocolControl
}

func (this *SrsSetChunkSizePacket) Decode(stream *SrsStream) error {
	var err error
	this.ChunkSize, err = stream.ReadInt32(binary.BigEndian)
	return err
}

func (this *SrsSetChunkSizePacket) Encode(stream *SrsStream) error {
	stream.WriteInt32(this.ChunkSize, binary.BigEndian)
	return nil
}
