package protocol

type SrsPeerBandwidthType int

// 5.6. Set Peer Bandwidth (6)
const (
	_                       SrsPeerBandwidthType = iota
	SrsPeerBandwidthHard                         = 0x00
	SrsPeerBandwidthSoft                         = 0x01
	SrsPeerBandwidthDynamic                      = 0x02
)

type SrsSetPeerBandwidthPacket struct {
	Bandwidth 	int32
	Type       	int8
}

func NewSrsSetPeerBandwidthPacket() *SrsSetPeerBandwidthPacket {
	return &SrsSetPeerBandwidthPacket{
		Bandwidth: 	0,
		Type:       SrsPeerBandwidthDynamic,
	}
}

func (this *SrsSetPeerBandwidthPacket) GetMessageType() int8 {
	return RTMP_MSG_SetPeerBandwidth
}

func (this *SrsSetPeerBandwidthPacket) GetPreferCid() int32 {
	return RTMP_CID_ProtocolControl
}

func (this *SrsSetPeerBandwidthPacket) Decode(stream *SrsStream) error {
	return nil
}

func (this *SrsSetPeerBandwidthPacket) Encode(stream *SrsStream) error {
	stream.WriteInt32(this.Bandwidth, binary.LittleEndian)
	stream.WriteByte(this.Type)
	return nil
}
