package packet

import (
	_ "log"
	"errors"
	"go_srs/srs/protocol/amf0"
	"go_srs/srs/utils"
	"go_srs/srs/global"
)

type SrsConnectAppResPacket struct {
	CommandName   amf0.SrsAmf0String
	TransactionId amf0.SrsAmf0Number
	Props         *amf0.SrsAmf0Object
	Info          *amf0.SrsAmf0Object
}

func NewSrsConnectAppResPacket() *SrsConnectAppResPacket {
	return &SrsConnectAppResPacket{
		CommandName:    amf0.SrsAmf0String{Value:amf0.SrsAmf0Utf8{Value:amf0.RTMP_AMF0_COMMAND_RESULT}},
		TransactionId:  amf0.SrsAmf0Number{Value:1},
		Props:          amf0.NewSrsAmf0Object(),
		Info:           amf0.NewSrsAmf0Object(),
	}
}

func (s *SrsConnectAppResPacket) GetMessageType() int8 {
	return global.RTMP_MSG_AMF0CommandMessage
}

func (s *SrsConnectAppResPacket) GetPreferCid() int32 {
	return global.RTMP_CID_OverConnection
}

func (this *SrsConnectAppResPacket) Decode(stream *utils.SrsStream) error {
	err := this.TransactionId.Decode(stream)
	if err != nil {
		return err
	}

	if this.TransactionId.Value != 1.0 {
		err := errors.New("amf0 decode connect transaction_id failed. ")
		return err
	}

	if !stream.Empty() {
		// marker, err := stream.PeekByte()
		// if err != nil {
		// 	return err
		// }

		// any := amf0.GenerateSrsAmf0Any(marker)
		// if any != nil {
		// 	any.Decode(stream)
		// }

		// switch any.(type) {
		// 	case amf0.SrsAmf0Object:
		// 		this.Props = any
		// 	default:
		// 		break
		// }
	}

	if err = this.Info.Decode(stream); err != nil {
		return err
	}
	return nil
}

func (this *SrsConnectAppResPacket) Encode(stream *utils.SrsStream) error {
	_ = this.CommandName.Encode(stream)
	_ = this.TransactionId.Encode(stream)
	_ = this.Props.Encode(stream)
	_ = this.Info.Encode(stream)
	return nil
}
