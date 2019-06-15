package packet

type SrsPacket interface {
	Decode(stream *SrsStream) error
	Encode(stream *SrsStream) error
	GetMessageType() int8
	GetPreferCid() int32
}
