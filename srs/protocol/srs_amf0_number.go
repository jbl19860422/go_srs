package protocol

/**
* read amf0 string from stream.
* 2.4 String Type
* string-type = string-marker UTF-8
* @return default value is empty string.
* @remark: use SrsAmf0Any::str() to create it.
 */
type SrsAmf0Number struct {
	SrsAmf0Any
	value float64
}

func (this *SrsAmf0Number) read(s *SrsStream) (err error) {
	this.value, err = srs_amf0_read_number(s)
	return
}
