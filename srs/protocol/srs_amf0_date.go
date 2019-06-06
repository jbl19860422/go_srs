package protocol

/**
* 2.13 Date Type
* time-zone = S16 ; reserved, not supported should be set to 0x0000
* date-type = date-marker DOUBLE time-zone
* @see: https://github.com/ossrs/srs/issues/185
*/
type SrsAmf0Date struct {
	SrsAmf0Any
	data_value 	int64
    time_zone	int64
}

func (this *SrsAmf0Date) read(s *SrsStream) (err error) {
	err = nil
	return
}