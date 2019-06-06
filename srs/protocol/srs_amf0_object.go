package protocol

/**
* 2.5 Object Type
* anonymous-object-type = object-marker *(object-property)
* object-property = (UTF-8 value-type) | (UTF-8-empty object-end-marker)
*/

type SrsAmf0Object struct {
	properties map[string]SrsAmf0Any
	eof *SrsAmf0ObjectEOF
}

func (this *SrsAmf0Object)read(s *SrsStream) (err error) {
	if marker, err := s.read_int8(); err != nil {
		return
	}

	if marker != RTMP_AMF0_Object {
		err = errors.New("amf0 check object marker failed. ")
		return
	}

	for !s.empty() {

	}
}
