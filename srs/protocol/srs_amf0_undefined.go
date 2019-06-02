package protocol

/**
* read amf0 undefined from stream.
* 2.8 undefined Type
* undefined-type = undefined-marker
 */
type SrsAmf0Undefined struct {
	SrsAmf0Any
	value string
}

func (this *SrsAmf0Undefined) read(s *SrsStream) (err error) {
	this.value, err = srs_amf0_read_string(s)
	return
}
