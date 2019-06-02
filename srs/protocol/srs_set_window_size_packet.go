package protocol

type SrsSetWindowAckSizePacket struct {
	Ackowledgement_window_size int32
}

func NewSrsSetWindowAckSizePacket() *SrsSetWindowAckSizePacket {
	return &SrsSetWindowAckSizePacket{}
}

func (s *SrsSetWindowAckSizePacket) get_message_type() int8 {
	return RTMP_MSG_WindowAcknowledgementSize
}

func (s *SrsSetWindowAckSizePacket) get_prefer_cid() int32 {
	return RTMP_CID_ProtocolControl
}

func (p *SrsSetWindowAckSizePacket) decode(s *SrsStream) error {
	var err error
	p.Ackowledgement_window_size, err = s.read_int32()
	return err
}

func (s *SrsSetWindowAckSizePacket) encode() ([]byte, error) {
	b := IntToBytes(int(s.Ackowledgement_window_size))
	return b, nil
}
