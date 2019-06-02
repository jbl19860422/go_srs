package protocol

type SrsPacket interface {
	decode(s *SrsStream) error
	encode() ([]byte, error)
	get_message_type() int8
	get_prefer_cid() int32
}
