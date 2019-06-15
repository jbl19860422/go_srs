package packet

import (
	"errors"
	"log"
)

type SrsConnectAppPacket struct {
	/**
	 * Name of the command. Set to "connect".
	 */
	CommandName SrsAmf0String
	/**
	 * Always set to 1.
	 */
	TransactionId SrsAmf0Number
	/**
    * Command information object which has the name-value pairs.
    * @remark: alloc in packet constructor, user can directly use it, 
    *       user should never alloc it again which will cause memory leak.
    * @remark, never be NULL.
    */
	CommandObj *SrsAmf0Object
	/**
    * Any optional information
    * @remark, optional, init to and maybe NULL.
    */
	Args *SrsAmf0Object
}

func NewSrsConnectAppPacket() *SrsConnectAppPacket {
	return &SrsConnectAppPacket{
		CommandName: SrsAmf0String{value:RTMP_AMF0_COMMAND_CONNECT},
		TransactionId: SrsAmf0Number{value:1},
		CommandObj: NewSrsAmf0Object(),
		Args: nil,
	}
}

func (s *SrsConnectAppPacket) GetMessageType() int8 {
	return RTMP_MSG_AMF0CommandMessage
}

func (s *SrsConnectAppPacket) GetPreferCid() int32 {
	return RTMP_CID_OverConnection
}

func (this *SrsConnectAppPacket) Decode(stream *SrsStream) error {
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

func (s *SrsConnectAppPacket) Encode(stream *SrsStream) error {
	_ = this.CommandName.Encode(stream)
	_ = this.TransactionId.Encode(stream)
	_ = this.CommandObj.Encode(stream)
	if this.Args != nil {
		_ = this.Args.Encode(stream)
	}
	return nil
}
