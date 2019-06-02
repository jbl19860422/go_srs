package protocol

/**
* read amf0 string from stream.
* 2.4 String Type
* string-type = string-marker UTF-8
* @return default value is empty string.
* @remark: use SrsAmf0Any::str() to create it.
 */
type SrsAmf0String struct {
	SrsAmf0Any
	value string
}

func (this *SrsAmf0String) read(s *SrsStream) (err error) {
	this.value, err = srs_amf0_read_string(s)
	return
}
