package packet
import (
	"go_srs/srs/protocol/amf0"
	"go_srs/srs/utils"
	"go_srs/srs/global"
)

type SrsCloseStreamPacket struct {
	CommandName		amf0.SrsAmf0String
	TransactionId	amf0.SrsAmf0Number
	NullObj			*amf0.SrsAmf0Object
}

func NewSrsCloseStreamPacket() *SrsCloseStreamPacket {
	return &SrsCloseStreamPacket{
		CommandName: amf0.SrsAmf0String{Value:amf0.SrsAmf0Utf8{Value:amf0.RTMP_AMF0_COMMAND_CLOSE_STREAM}},
		TransactionId: amf0.SrsAmf0Number{Value:0},
	}
}

func (s *SrsCloseStreamPacket) GetMessageType() int8 {
	return global.RTMP_MSG_AMF0CommandMessage
}

func (s *SrsCloseStreamPacket) GetPreferCid() int32 {
	return global.RTMP_CID_OverStream
}

func (this *SrsCloseStreamPacket) Decode(stream *utils.SrsStream) error {
	var err error
	err = this.TransactionId.Decode(stream)
	if err != nil {
		return err
	}

	err = this.NullObj.Decode(stream)
	if err != nil {
		return err
	}

	return nil
}

func (this *SrsCloseStreamPacket) Encode(stream *utils.SrsStream) error {
	return nil
}