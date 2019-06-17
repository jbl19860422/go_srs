package packet

import (
	"go_srs/srs/global"
	"go_srs/srs/protocol/amf0"
	"go_srs/srs/utils"
	"errors"
)

//todo copy comment
type SrsPlayPacket struct {
	CommandName 	amf0.SrsAmf0String
	TransactionId 	amf0.SrsAmf0Number
	NullObj			amf0.SrsAmf0Null
	StreamName		amf0.SrsAmf0String
	Start			amf0.SrsAmf0Number
	Duration		amf0.SrsAmf0Number
	Reset			amf0.SrsAmf0Boolean
}

func NewSrsPlayPacket() *SrsPlayPacket {
	return &SrsPlayPacket{
		CommandName:amf0.SrsAmf0String{Value:amf0.SrsAmf0Utf8{Value:amf0.RTMP_AMF0_COMMAND_PLAY}},
		TransactionId:amf0.SrsAmf0Number{Value:0},
		Start:amf0.SrsAmf0Number{Value:-2.0},
		Duration:amf0.SrsAmf0Number{Value:-1.0},
		Reset:amf0.SrsAmf0Boolean{Value:true},
	}
}

func (s *SrsPlayPacket) GetMessageType() int8 {
	return global.RTMP_MSG_AMF0CommandMessage
}

func (s *SrsPlayPacket) GetPreferCid() int32 {
	return global.RTMP_CID_OverStream
}

func (this *SrsPlayPacket) Decode(stream *utils.SrsStream) error {
	if err := this.TransactionId.Decode(stream); err != nil {
		return err
	}

	if err := this.NullObj.Decode(stream); err != nil {
		return err
	}

	if err := this.StreamName.Decode(stream); err != nil {
		return err
	}

	if len(this.StreamName.GetValue().(string)) > 0 {
		if err := this.Start.Decode(stream); err != nil {
			return err
		}
	}
	//todo fix this
	return nil

	if len(this.StreamName.GetValue().(string)) > 0 {
		if err := this.Duration.Decode(stream); err != nil {
			return err
		}
	}

	if len(this.StreamName.GetValue().(string)) <= 0 {
		return nil
	}

	marker,err := stream.PeekByte()
	if err != nil {
		return errors.New("no reset field")
	}
	switch marker {
	case amf0.RTMP_AMF0_Boolean:
		if err := this.Reset.Decode(stream); err !=  nil {
			return err
		}
	case amf0.RTMP_AMF0_Number:
		n := amf0.SrsAmf0Number{}
		if err := n.Decode(stream); err != nil {
			return err
		}

		if n.GetValue() != 0 {
			this.Reset.Value = true
		} else {
			this.Reset.Value = false
		}
	default:
		return errors.New("amf0 invalid, the reset requires number or bool")
	}

	return nil
}

func (this *SrsPlayPacket) Encode(stream *utils.SrsStream) error {
	_ = this.CommandName.Encode(stream)
	_ = this.TransactionId.Encode(stream)
	_ = this.NullObj.Encode(stream)
	_ = this.StreamName.Encode(stream)
	//why
	if this.Start.GetValue().(float64) != -2 || this.Duration.GetValue().(float64) != -1 || !this.Reset.GetValue().(bool) {
		err := this.Start.Encode(stream)
		if err != nil {
			return err
		}
	}
	//why
	if this.Duration.GetValue().(float64) != -1 || !this.Reset.GetValue().(bool) {
		_ = this.Duration.Encode(stream)
	}

	if !this.Reset.GetValue().(bool) {
		_ = this.Reset.Encode(stream)
	}

	return nil
}
