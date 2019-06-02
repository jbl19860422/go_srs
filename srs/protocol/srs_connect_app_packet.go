package protocol

import (
	"errors"
	"log"
)

type SrsConnectAppPacket struct {
	/**
	 * Name of the command. Set to "connect".
	 */
	command_name string
	/**
	 * Always set to 1.
	 */
	transaction_id float64

	CommandObj *SrsAmf0Object
}

func NewSrsConnectAppPacket() *SrsConnectAppPacket {
	return &SrsConnectAppPacket{
		command_name: "connect",
		CommandObj:   NewSrsAmf0Object(),
	}
}

func (s *SrsConnectAppPacket) get_message_type() int8 {
	return RTMP_MSG_AMF0CommandMessage
}

func (s *SrsConnectAppPacket) get_prefer_cid() int32 {
	return RTMP_CID_OverConnection
}

func (this *SrsConnectAppPacket) decode(s *SrsStream) error {
	var err error
	this.transaction_id, err = srs_amf0_read_number(s)
	if err != nil {
		return err
	}

	if this.transaction_id != 1.0 {
		log.Printf("amf0 decode connect transaction_id failed.%.1f", this.transaction_id)
		err = errors.New("amf0 decode connect transaction_id failed.")
		return err
	}

	err = this.CommandObj.read(s)
	if err != nil {
		log.Print("command read failed")
		return err
	} else {
		log.Print("command_obj read succeed")
	}

	log.Print("properties len = ", len(this.CommandObj.properties))
	return nil
}

func (s *SrsConnectAppPacket) encode() ([]byte, error) {
	stream := NewSrsStream([]byte{}, 0)
	err := srs_amf0_write_string(stream, s.command_name)
	if err != nil {
		return nil, err
	}

	err = srs_amf0_write_number(stream, s.transaction_id)
	if err != nil {
		return nil, err
	}

	// if ((ret = props->write(stream)) != ERROR_SUCCESS) {
	//     srs_error("encode props failed. ret=%d", ret);
	//     return ret;
	// }

	// srs_verbose("encode props success.");

	// if ((ret = info->write(stream)) != ERROR_SUCCESS) {
	//     srs_error("encode info failed. ret=%d", ret);
	//     return ret;
	// }
	return nil, nil
}
