package protocol

import (
	"errors"
	"log"
)

const (
	RTMP_AMF0_Number      = 0x00
	RTMP_AMF0_Boolean     = 0x01
	RTMP_AMF0_String      = 0x02
	RTMP_AMF0_Object      = 0x03
	RTMP_AMF0_MovieClip   = 0x04 // reserved, not supported
	RTMP_AMF0_Null        = 0x05
	RTMP_AMF0_Undefined   = 0x06
	RTMP_AMF0_Reference   = 0x07
	RTMP_AMF0_EcmaArray   = 0x08
	RTMP_AMF0_ObjectEnd   = 0x09
	RTMP_AMF0_StrictArray = 0x0A
	RTMP_AMF0_Date        = 0x0B
	RTMP_AMF0_LongString  = 0x0C
	RTMP_AMF0_UnSupported = 0x0D
	RTMP_AMF0_RecordSet   = 0x0E
	RTMP_AMF0_XmlDocument = 0x0F
	RTMP_AMF0_TypedObject = 0x10
	// AVM+ object is the AMF3 object.
	RTMP_AMF0_AVMplusObject = 0x11
	// origin array whos data takes the same form as LengthValueBytes
	RTMP_AMF0_OriginStrictArray = 0x20
	// User defined
	RTMP_AMF0_Invalid = 0x3F
)

func srs_amf0_read_string(s *SrsStream) (val string, err error) {
	log.Print("srs_amf0_read_string start")
	if !s.require(1) {
		err = errors.New("amf0 read string marker failed")
		return
	}

	marker, err := s.read_nbytes(1)
	bmarker := marker[0]
	if int(bmarker) != RTMP_AMF0_String {
		err = errors.New("amf0 check string marker failed.")
	}

	log.Print("marker is string")

	val, err = srs_amf0_read_utf8(s)
	return
}

func srs_amf0_read_utf8(s *SrsStream) (val string, err error) {
	if !s.require(2) {
		err = errors.New("amf0 read string length failed")
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
}

func srs_amf0_read_number(s *SrsStream) (val float64, err error) {
	marker, err := s.read_int8()
	if err != nil {
		return
	}

	if marker != RTMP_AMF0_Number {
		err = errors.New("amf0 check number marker failed.")
		return
	}

	val, err = s.read_float64()
	if err != nil {
		return 0, err
	}
	return
}

func srs_amf0_read_boolean(s *SrsStream) (val bool, err error) {
	marker, err := s.read_int8()
	if err != nil {
		return
	}

	if marker != RTMP_AMF0_Number {
		err = errors.New("amf0 check number marker failed.")
		return
	}

	val, err = s.read_bool()
	if err != nil {
		return 0, err
	}
	return
}

func srs_amf0_read_null(s *SrsStream) (err error) {
	if marker, err := s.read_int8(); err != nil {
		return
	}

	if marker != RTMP_AMF0_Null {
		err = errors.New("amf0 check undefined marker failed")
	}

	return
}

func srs_amf0_read_undefined(SrsStream* stream) (err error) {
	if marker, err := s.read_int8(); err != nil {
		return
	}

	if marker != RTMP_AMF0_Undefined {
		err = errors.New("amf0 check undefined marker failed")
	}
	return 
}

func srs_amf0_is_object_eof(s *SrsStream) bool {
	// detect the object-eof specially
	if (stream->require(3)) {
		flag_buf := stream->read_nbytes(3);
		stream->skip(-3);
		bin_buf := bytes.NewBuffer(flag_buf)
		var flag int32
		binary.Read(bin_buf, binary.BigEndian, &flag)
		return 0x09 == flag
	}
	
	return false;
}

type SrsAmf0AnyInterface interface {
	is_string() bool
	is_boolean() bool
	is_number() bool
	is_null() bool
	is_undefined() bool
	is_object() bool
	is_object_eof() bool
	is_ecma_array() bool
	is_strict_array() bool
	is_date() bool
	is_complex_object() bool
	total_size() int32
	read(s *SrsStream) error
	write(s *SrsStream) error
	copy(s *SrsAmf0Any) error
}
/**
* any amf0 value.
* 2.1 Types Overview
* value-type = number-type | boolean-type | string-type | object-type 
*         | null-marker | undefined-marker | reference-type | ecma-array-type 
*         | strict-array-type | date-type | long-string-type | xml-document-type 
*         | typed-object-type
*/
type SrsAmf0Any struct {
	marker int8
}

func (s *SrsAmf0Any) is_string() bool {
	return s.marker == RTMP_AMF0_String;
}

func (s *SrsAmf0Any) is_boolean() bool {
	return s.marker == RTMP_AMF0_Boolean
}

func (s *SrsAmf0Any) is_number() bool {
	return s.marker == RTMP_AMF0_Number
}

func (s *SrsAmf0Any) is_null() bool {
	return s.marker == RTMP_AMF0_Null
}

func (s *SrsAmf0Any) is_undefined() bool {
	return s.marker == RTMP_AMF0_Undefined
}

func (s *SrsAmf0Any) is_object() bool {
	return s.marker == RTMP_AMF0_Object
}

func (s *SrsAmf0Any) is_object_eof() bool {
    return s.marker == RTMP_AMF0_ObjectEnd
}

func (s *SrsAmf0Any) is_ecma_array() bool {
	return s.marker == RTMP_AMF0_EcmaArray
}

func (s *SrsAmf0Any) is_strict_array() bool {
	return s.marker == RTMP_AMF0_StrictArray
}

func (s *SrsAmf0Any) is_date() bool {
	return s.marker == RTMP_AMF0_Date
}

func (s *SrsAmf0Any) is_complex_object() bool {
	return s.is_object() || s.is_object_eof() || s.is_ecma_array() || s.is_strict_array()
}


/**
* 2.5 Object Type
* anonymous-object-type = object-marker *(object-property)
* object-property = (UTF-8 value-type) | (UTF-8-empty object-end-marker)
*/
type SrsAmf0Object struct {

}












