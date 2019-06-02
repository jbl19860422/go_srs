package protocol

/**
* read amf0 boolean from stream.
* 2.4 String Type
* boolean-type = boolean-marker U8
*         0 is false, <> 0 is true
* @return default value is false.
 */
type SrsAmf0Boolean struct {
	SrsAmf0Any
	value bool
}

func (this *SrsAmf0Boolean) read(s *SrsStream) (err error) {
	this.value, err = srs_amf0_read_boolean(s)
	return
}
