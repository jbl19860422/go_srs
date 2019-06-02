package protocol

import (
	"errors"
	"log"
)

type SrsFMLEStartPacket struct {
	command_name   string
	transaction_id float64
	Stream_name    string
}

func NewSrsFMLEStartPacket() *SrsFMLEStartPacket {
	return &SrsFMLEStartPacket{}
}

func (s *SrsFMLEStartPacket) get_message_type() int8 {
	return RTMP_MSG_AMF0CommandMessage
}

func (s *SrsFMLEStartPacket) get_prefer_cid() int32 {
	return RTMP_CID_OverConnection
}

func (this *SrsFMLEStartPacket) decode(stream *SrsStream) error {
	var err error
	this.command_name, err = srs_amf0_read_string(stream)
	if err != nil {
		return err
	}

	if len(this.command_name) <= 0 || (this.command_name != RTMP_AMF0_COMMAND_RELEASE_STREAM && this.command_name != RTMP_AMF0_COMMAND_FC_PUBLISH && this.command_name != RTMP_AMF0_COMMAND_UNPUBLISH) {
		return errors.New("amf0 decode FMLE start command_name failed.")
	}

	this.transaction_id, err = srs_amf0_read_number(stream)
	if err != nil {
		return err
	}

	err = srs_amf0_read_null(stream)
	if err != nil {
		return err
	}

	this.Stream_name, err = srs_amf0_read_string(stream)
	if err != nil {
		return errors.New("amf0 decode FMLE start stream_name failed")
	}
	log.Print("StreamName=", this.Stream_name)
	return nil
}

func (s *SrsFMLEStartPacket) encode() ([]byte, error) {
	stream := NewSrsStream([]byte{}, 0)
	err := srs_amf0_write_string(stream, s.command_name)
	if err != nil {
		return nil, err
	}

	err = srs_amf0_write_number(stream, s.transaction_id)
	if err != nil {
		return nil, err
	}

	return nil, nil
}
