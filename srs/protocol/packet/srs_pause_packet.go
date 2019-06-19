package packet
import (
	"go_srs/srs/protocol/amf0"
	"go_srs/srs/utils"
	"go_srs/srs/global"
)
type SrsPausePacket struct {
	CommandName		amf0.SrsAmf0String
	TransactionId	amf0.SrsAmf0Number
	NullObj			amf0.SrsAmf0Null
	IsPause			amf0.SrsAmf0Boolean
	TimeMs			amf0.SrsAmf0Number
}

func (s *SrsPausePacket) GetMessageType() int8 {
	return global.RTMP_MSG_AMF0CommandMessage
}

func (s *SrsPausePacket) GetPreferCid() int32 {
	return global.RTMP_CID_OverStream
}

func (this *SrsPausePacket) Decode(stream *utils.SrsStream) error {
	var err error
	if err = this.TransactionId.Decode(stream); err != nil {
		return err
	}

	if err = this.NullObj.Decode(stream); err != nil {
		return err
	}

	if err = this.IsPause.Decode(stream); err != nil {
		return err
	}

	if err = this.TimeMs.Decode(stream); err != nil {
		return err
	}

	return nil
}

func (this *SrsPausePacket) Encode(stream *utils.SrsStream) error {
	return nil
}