package protocol

import "errors"

type SrsFMLEStartResPacket struct {
	command_name   string
	Transaction_id float64
	CommandObj     *SrsAmf0Object
	Args           *SrsAmf0Object
}

func NewSrsFMLEStartResPacket(trans_id float64) *SrsFMLEStartResPacket {
	return &SrsFMLEStartResPacket{
		command_name:   RTMP_AMF0_COMMAND_RESULT,
		Transaction_id: trans_id,
		CommandObj:     NewSrsAmf0Object(),
		Args:           NewSrsAmf0Object(),
	}
}

func (s *SrsFMLEStartResPacket) get_message_type() int8 {
	return RTMP_MSG_AMF0CommandMessage
}

func (s *SrsFMLEStartResPacket) get_prefer_cid() int32 {
	return RTMP_CID_OverConnection
}

func (p *SrsFMLEStartResPacket) decode(stream *SrsStream) error {
	var err error
	p.command_name, err = srs_amf0_read_string(stream)
	if err != nil {
		return err
	}

	if len(p.command_name) <= 0 || p.command_name != RTMP_AMF0_COMMAND_RESULT {
		return errors.New("amf0 decode FMLE start response command_name failed.")
	}

	p.Transaction_id, err = srs_amf0_read_number(stream)
	if err != nil {
		return err
	}

	err = srs_amf0_read_null(stream)
	if err != nil {
		return err
	}

	err = srs_amf0_read_undefined(stream)
	if err != nil {
		return err
	}

	return nil
}

func (this *SrsFMLEStartResPacket) encode() ([]byte, error) {
	stream := NewSrsStream([]byte{}, 0)

	err := srs_amf0_write_string(stream, this.command_name)
	if err != nil {
		return nil, err
	}

	err = srs_amf0_write_number(stream, this.Transaction_id)
	if err != nil {
		return nil, err
	}

	err = srs_amf0_write_null(stream)
	if err != nil {
		return nil, err
	}

	err = srs_amf0_write_undefined(stream)
	if err != nil {
		return nil, err
	}

	return stream.p, nil
}
