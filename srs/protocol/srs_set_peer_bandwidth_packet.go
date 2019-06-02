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
	Bandwidth int32
	Typ       int8
}

func NewSrsSetPeerBandwidthPacket() *SrsSetPeerBandwidthPacket {
	return &SrsSetPeerBandwidthPacket{
		Bandwidth: 0,
		Typ:       SrsPeerBandwidthDynamic,
	}
}

func (s *SrsSetPeerBandwidthPacket) get_message_type() int8 {
	return RTMP_MSG_SetPeerBandwidth
}

func (s *SrsSetPeerBandwidthPacket) get_prefer_cid() int32 {
	return RTMP_CID_ProtocolControl
}

func (p *SrsSetPeerBandwidthPacket) decode(s *SrsStream) error {
	var err error
	return err
}

func (s *SrsSetPeerBandwidthPacket) encode() ([]byte, error) {
	b := IntToBytes(int(s.Bandwidth))
	b = append(b, byte(s.Typ))
	return b, nil
}
