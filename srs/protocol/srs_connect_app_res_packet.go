package protocol

import (
	"log"
)

type SrsConnectAppResPacket struct {
	command_name   string
	transaction_id float64
	Props          *SrsAmf0Object
	Info           *SrsAmf0Object
}

func NewSrsConnectAppResPacket() *SrsConnectAppResPacket {
	return &SrsConnectAppResPacket{
		command_name:   RTMP_AMF0_COMMAND_RESULT,
		transaction_id: 1,
		Props:          NewSrsAmf0Object(),
		Info:           NewSrsAmf0Object(),
	}
}

func (s *SrsConnectAppResPacket) get_message_type() int8 {
	return RTMP_MSG_AMF0CommandMessage
}

func (s *SrsConnectAppResPacket) get_prefer_cid() int32 {
	return RTMP_CID_OverConnection
}

func (this *SrsConnectAppResPacket) decode(s *SrsStream) error {
	return nil
}

func (s *SrsConnectAppResPacket) encode() ([]byte, error) {
	stream := NewSrsStream([]byte{}, 0)
	err := srs_amf0_write_string(stream, s.command_name)
	if err != nil {
		return nil, err
	}

	err = srs_amf0_write_number(stream, s.transaction_id)
	if err != nil {
		return nil, err
	}

	err = s.Props.write(stream)
	if err != nil {
		return nil, err
	}

	err = s.Info.write(stream)
	if err != nil {
		return nil, err
	}

	log.Print("xxxxxxxxxxxstream encode len=", len(stream.p), "   xxxxxxxxxxxxxxxxx")
	return stream.p, nil
}
