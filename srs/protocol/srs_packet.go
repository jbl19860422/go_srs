package protocol

type SrsPacket interface {
	decode(s *SrsStream) error
}
