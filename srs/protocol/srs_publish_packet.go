package protocol

import (
	"errors"
)

type SrsPublishPacket struct {
	command_name   string
	Transaction_id float64
	CommandObj     *SrsAmf0Object
	Stream_name    string
	typ            string
}

func NewSrsPublishPacket() *SrsPublishPacket {
	return &SrsPublishPacket{
		command_name:   RTMP_AMF0_COMMAND_PUBLISH,
		Transaction_id: 0,
		CommandObj:     NewSrsAmf0Object(),
		typ:            "live",
	}
}

func (s *SrsPublishPacket) get_message_type() int8 {
	return RTMP_MSG_AMF0CommandMessage
}

func (s *SrsPublishPacket) get_prefer_cid() int32 {
	return RTMP_CID_OverStream
}

func (this *SrsPublishPacket) decode(stream *SrsStream) error {
	var err error
	this.Transaction_id, err = srs_amf0_read_number(stream)
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

	if !stream.empty() {
		this.typ, err = srs_amf0_read_string(stream)
	}
	return nil
}

func (s *SrsPublishPacket) encode() ([]byte, error) {
	stream := NewSrsStream([]byte{}, 0)
	err := srs_amf0_write_string(stream, s.command_name)
	if err != nil {
		return nil, err
	}

	err = srs_amf0_write_number(stream, s.Transaction_id)
	if err != nil {
		return nil, err
	}

	err = srs_amf0_write_null(stream)
	if err != nil {
		return nil, err
	}

	err = srs_amf0_write_string(stream, s.Stream_name)
	if err != nil {
		return nil, err
	}

	err = srs_amf0_write_string(stream, s.typ)
	if err != nil {
		return nil, err
	}

	return nil, nil
}
