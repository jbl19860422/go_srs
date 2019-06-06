package protocol

/**
* 2.11 Object End Type
* object-end-type = UTF-8-empty object-end-marker
* 0x00 0x00 0x09
*/

type SrsAmf0ObjectEOF struct {
	SrsAmf0Any
}


func (this *SrsAmf0ObjectEOF) read(s *SrsStream) (err error) {
    if temp, err := s.read_int16(); err != nil {
		return
	}

    if temp != 0x00 {
		err = errors.New("amf0 read object eof value check failed.")
		return
    }
    
	// marker
	if marker, err := s.read_int8(); err != nil {
		return
	}

	if marker != RTMP_AMF0_ObjectEnd {
		err = errors.New("amf0 check object eof marker failed. ")
		return
	}
	return
}


