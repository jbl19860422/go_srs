package protocol

import "errors"

/**
* 2.11 Object End Type
* object-end-type = UTF-8-empty object-end-marker
* 0x00 0x00 0x09
 */

type SrsAmf0ObjectEOF struct {
	SrsAmf0Any
}

func (this *SrsAmf0ObjectEOF) read(s *SrsStream) (err error) {
	var temp int16
	if temp, err = s.read_int16(); err != nil {
		return
	}

	if temp != 0x00 {
		err = errors.New("amf0 read object eof value check failed.")
		return
	}

	// marker
	var marker int8
	if marker, err = s.read_int8(); err != nil {
		return
	}

	if marker != RTMP_AMF0_ObjectEnd {
		err = errors.New("amf0 check object eof marker failed. ")
		return
	}
	return
}

func (this *SrsAmf0ObjectEOF) write(stream *SrsStream) error {
	b := make([]byte, 3)
	b[0] = 0
	b[1] = 0
	b[2] = 9
	stream.write_bytes(b)
	return nil
}
