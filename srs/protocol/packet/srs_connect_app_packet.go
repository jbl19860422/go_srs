package packet

import (
	"errors"
	"log"
	"go_srs/srs/protocol/amf0"
	"go_srs/srs/utils"
	"go_srs/srs/global"
)

type SrsConnectAppPacket struct {
	CommandName 	amf0.SrsAmf0String
	TransactionId 	amf0.SrsAmf0Number
	CommandObj 		*amf0.SrsAmf0Object
	Args 			*amf0.SrsAmf0Object
}

func NewSrsConnectAppPacket() *SrsConnectAppPacket {
	return &SrsConnectAppPacket{
		CommandName: amf0.SrsAmf0String{Value:amf0.SrsAmf0Utf8{Value:amf0.RTMP_AMF0_COMMAND_CONNECT}},
		TransactionId: amf0.SrsAmf0Number{Value:1},
		CommandObj: amf0.NewSrsAmf0Object(),
		Args: nil,
	}
}

func (s *SrsConnectAppPacket) GetMessageType() int8 {
	return global.RTMP_MSG_AMF0CommandMessage
}

func (s *SrsConnectAppPacket) GetPreferCid() int32 {
	return global.RTMP_CID_OverConnection
}

func (this *SrsConnectAppPacket) Decode(stream *utils.SrsStream) error {
	var err error
	err = this.TransactionId.Decode(stream)
	if err != nil {
		return err
	}

	if this.TransactionId.Value != 1.0 {
		log.Printf("amf0 decode connect transaction_id failed.%.1f", this.TransactionId.Value)
		err = errors.New("amf0 decode connect transaction_id failed.")
		return err
	}

	err = this.CommandObj.Decode(stream)
	if err != nil {
		log.Print("command read failed")
		return err
	} else {
		log.Print("command_obj read succeed")
	}

	if !stream.Empty() {
		err = this.Args.Decode(stream)
		if err != nil {
			return err
		}
	}

	return nil
}

func (this *SrsConnectAppPacket) Encode(stream *utils.SrsStream) error {
	_ = this.CommandName.Encode(stream)
	_ = this.TransactionId.Encode(stream)
	_ = this.CommandObj.Encode(stream)
	if this.Args != nil {
		_ = this.Args.Encode(stream)
	}
	return nil
}
