package amf0

type SrsAmf0String struct {
	value 	string
}

func (this *SrsAmf0String) Decode(stream *SrsStream) error {
	marker, err := stream.ReadByte()
	if err != nil {
		return err
	}

	if marker != RTMP_AMF0_String {
		err := errors.New("amf0 check string marker failed.")
		return
	}

	len, err := s.read_int16()
	if err != nil {
		return
	}
	log.Print("utf8 len=", len)
	if len <= 0 {
		err = errors.New("amf0 read empty string.")
		return
	}

	val, err = s.read_string(int32(len))
	return

	return
}
