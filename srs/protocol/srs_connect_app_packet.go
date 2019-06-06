package protocol

import (
	"log"
	"errors"
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
}

func NewSrsConnectAppPacket() *SrsConnectAppPacket {
	return &SrsConnectAppPacket{command_name:"connect"}
}

func (this *SrsConnectAppPacket) decode(s *SrsStream) error {
	var err error
	this.transaction_id, err = srs_amf0_read_number(s)
	if err != nil {
		log.Print("srs_amf0_read_string 2222222222222222")
		return err
	}

	if this.transaction_id != 1.0 {
		log.Printf("amf0 decode connect transaction_id failed.%.1f", this.transaction_id)
		err = errors.New("amf0 decode connect transaction_id failed.")
		return err
	}

	log.Printf("transaction_id=%.1f", this.transaction_id)

	return nil
}

func (s *SrsConnectAppPacket) encode(payload []byte, size int32) int32 {
	return 0
}
