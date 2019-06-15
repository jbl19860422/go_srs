package packet

import (
	"log"
)

type SrsConnectAppResPacket struct {
	CommandName   SrsAmf0String
	TransactionId SrsAmf0Number
	Props         *SrsAmf0Object
	Info          *SrsAmf0Object
}

func NewSrsConnectAppResPacket() *SrsConnectAppResPacket {
	return &SrsConnectAppResPacket{
		CommandName:   SrsAmfpString{Value:RTMP_AMF0_COMMAND_RESULT},
		transaction_id: 1,
		Props:          NewSrsAmf0Object(),
		Info:           NewSrsAmf0Object(),
	}
}

func (s *SrsConnectAppResPacket) GetMessageType() int8 {
	return RTMP_MSG_AMF0CommandMessage
}

func (s *SrsConnectAppResPacket) GetPreferCid() int32 {
	return RTMP_CID_OverConnection
}

func (this *SrsConnectAppResPacket) Decode(stream *SrsStream) error {
	err := this.TransactionId.Decode(stream)
	if err != nil {
		return err
	}

	if this.TransactionId.Value != 1.0 {
		err := errors.New("amf0 decode connect transaction_id failed. ")
		return err
	}

	if !stream.Empty() {
		marker, err := stream.PeekByte()
		if err != nil {
			return err
		}

		any := GenerateSrsAmf0Any(marker)
		if any != nil {
			any.Decode(stream)
		}

		switch any.(type) {
			case SrsAmf0Object:
				this.Props = any
			default:
				break
		}
	}

	err := this.Info.Decode(stream)
	if err != nil {
		return err
	}
	return nil
}

func (s *SrsConnectAppResPacket) Encode(stream *SrsStream) error {
	_ = this.CommandName.Encode(stream)
	_ = this.TransactionId.Encode(stream)
	_ = this.Props.Encode(stream)
	_ = this.Info.Encode(stream)
	return nil
}
