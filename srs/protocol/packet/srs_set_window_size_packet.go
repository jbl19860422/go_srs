package packet

type SrsSetWindowAckSizePacket struct {
	AckowledgementWindowSize int32
}

func NewSrsSetWindowAckSizePacket() *SrsSetWindowAckSizePacket {
	return &SrsSetWindowAckSizePacket{}
}

func (this *SrsSetWindowAckSizePacket) GetMessageType() int8 {
	return RTMP_MSG_WindowAcknowledgementSize
}

func (this *SrsSetWindowAckSizePacket) GetPreferCid() int32 {
	return RTMP_CID_ProtocolControl
}

func (p *SrsSetWindowAckSizePacket) Decode(stream *SrsStream) error {
	var err error
	this.AckowledgementWindowSize, err = stream.ReadInt32(binary.LittleEndian)
	return err
}

func (s *SrsSetWindowAckSizePacket) Encode(stream *SrsStream) error {
	stream.WriteInt32(this.AckowledgementWindowSize, binary.LittleEndian)
	return nil
}
