package packet

import(
	"go_srs/srs/utils"
	"go_srs/srs/protocol/amf0"
	"go_srs/srs/global"
)
type SrsCreateStreamResPacket struct {
	CommandName   	amf0.SrsAmf0String
	TransactionId 	amf0.SrsAmf0Number
	NullObj			amf0.SrsAmf0Null
	CommandObj     	*amf0.SrsAmf0Object
	StreamId      	amf0.SrsAmf0Number
}

func NewSrsCreateStreamResPacket(tid float64, sid float64) *SrsCreateStreamResPacket {
	return &SrsCreateStreamResPacket{
		CommandName:   	amf0.SrsAmf0String{Value:amf0.SrsAmf0Utf8{Value:amf0.RTMP_AMF0_COMMAND_RESULT}},
		TransactionId: 	amf0.SrsAmf0Number{Value:tid},
		CommandObj:     amf0.NewSrsAmf0Object(),
		StreamId:      	amf0.SrsAmf0Number{Value:sid},
	}
}
func (s *SrsCreateStreamResPacket) GetMessageType() int8 {
	return global.RTMP_MSG_AMF0CommandMessage
}

func (s *SrsCreateStreamResPacket) GetPreferCid() int32 {
	return global.RTMP_CID_OverConnection
}

func (this *SrsCreateStreamResPacket) Decode(stream *utils.SrsStream) error {
	if err := this.TransactionId.Decode(stream); err != nil {
		return err
	}
	
	if err := this.NullObj.Decode(stream); err != nil {
		return err
	}

	if err := this.StreamId.Decode(stream); err != nil {
		return err
	}

	return nil
}

func (this *SrsCreateStreamResPacket) Encode(stream *utils.SrsStream) error {
	_ = this.CommandName.Encode(stream)
	_ = this.TransactionId.Encode(stream)
	_ = this.NullObj.Encode(stream)
	_ = this.StreamId.Encode(stream)
	return nil
}
