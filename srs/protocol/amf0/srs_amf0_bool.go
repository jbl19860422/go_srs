package amf0
type SrsAmf0Boolean struct {
	value bool
}

func NewSrsAmf0Boolean(data bool) *SrsAmf0Boolean {
	return &SrsAmf0Boolean{
		value:data
	}
}

func (this *SrsAmf0Boolean) Decode(stream *utils.SrsStream) error {
	marker, err := stream.ReadByte()
	if err != nil {
		return err
	}

	if marker != RTMP_AMF0_Boolean {
		err := errors.New("amf0 check bool marker failed.")
		return err
	}

	this.value, err = stream.ReadBool()
	if err != nil {
		return err
	}
	return nil
}

func (this *SrsAmf0Boolean) Encode(stream *utils.SrsStream) error {
	stream.WriteByte(RTMP_AMF0_Boolean)
	stream.WriteByte(byte(this.value))
	return nil
}

func (this *SrsAmf0Boolean) IsMyType(stream *utils.SrsStream) (bool, error) {
	marker, err := stream.PeekByte(1)
	if err != nil {
		return false, err
	}

	if marker != RTMP_AMF0_Boolean {
		return false, nil
	}
	return true, nil
}


