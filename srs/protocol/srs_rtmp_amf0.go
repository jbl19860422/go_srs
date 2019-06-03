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

func srs_amf0_write_string(s *SrsStream, val string) error {
	marker := byte(RTMP_AMF0_String)
	len := int16(len(val))
	len_buf := make([]byte, 2)
	len_buf[0] = byte((len >> 8) & 0xFF)
	len_buf[1] = byte((len) & 0xFF)
	b := []byte(val)
	data := make([]byte, 0)
	data = append(data, marker)
	data = append(data, len_buf...)
	data = append(data, b...)
	s.write_bytes(data)
	return nil
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

func srs_amf0_write_utf8(s *SrsStream, val string) error {
	len := int16(len(val))
	len_buf := make([]byte, 2)
	len_buf[0] = byte((len >> 8) & 0x0F)
	len_buf[1] = byte((len) & 0x0F)
	b := []byte(val)
	data := make([]byte, 0)
	data = append(data, len_buf...)
	data = append(data, b...)
	s.write_bytes(data)
	return nil
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

func srs_amf0_write_number(s *SrsStream, val float64) error {
	marker := byte(RTMP_AMF0_Number)
	s.write_1byte(marker)
	s.write_float64(val)
	return nil
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
		return false, err
	}
	return
}

func srs_amf0_read_null(s *SrsStream) (err error) {
	var marker int8
	if marker, err = s.read_int8(); err != nil {
		return
	}

	if marker != RTMP_AMF0_Null {
		err = errors.New("amf0 check undefined marker failed")
	}

	return
}

func srs_amf0_write_null(stream *SrsStream) error {
	marker := byte(RTMP_AMF0_Null)
	stream.write_1byte(marker)
	return nil
}

func srs_amf0_read_undefined(stream *SrsStream) (err error) {
	var marker int8
	if marker, err = stream.read_int8(); err != nil {
		return
	}

	if marker != RTMP_AMF0_Undefined {
		err = errors.New("amf0 check undefined marker failed")
	}
	return
}

func srs_amf0_write_undefined(stream *SrsStream) error {
	marker := byte(RTMP_AMF0_Undefined)
	stream.write_1byte(marker)
	return nil
}

func srs_amf0_is_object_eof(stream *SrsStream) bool {
	// detect the object-eof specially
	if stream.require(3) { //marker = 9（RTMP_AMF0_ObjectEnd），后面带两个0表示结束
		flag_buf, err := stream.read_nbytes(3)
		if err != nil {
			return false
		}

		stream.skip(-3)
		log.Printf("flag_buf=%x %x %x", flag_buf[0], flag_buf[1], flag_buf[2])
		if flag_buf[0] == 0 && flag_buf[1] == 0 && flag_buf[2] == 9 {
			return true
		}
		return false
	}

	return false
}

func decodeAmf0(stream *SrsStream) (v interface{}, err error) {
	var marker int8
	if marker, err = stream.read_int8(); err != nil {
		return
	}
	stream.skip(-1)

	switch marker {
	case RTMP_AMF0_Number:
		{
			v, err = srs_amf0_read_number(stream)
			break
		}
	case RTMP_AMF0_Boolean:
		{
			v, err = srs_amf0_read_number(stream)
			break
		}
	case RTMP_AMF0_String:
		{
			v, err = srs_amf0_read_string(stream)
			log.Print("value=", v)
			break
		}
	case RTMP_AMF0_Object:
		{
			v, err = decodeAmf0(stream)
			break
		}
	case RTMP_AMF0_MovieClip:
		{

		}
	case RTMP_AMF0_Null:
		{
			err = srs_amf0_read_null(stream)
			break
		}
	case RTMP_AMF0_Undefined:
		{
			err = srs_amf0_read_undefined(stream)
			break
		}
	case RTMP_AMF0_Reference:
		{

		}
	case RTMP_AMF0_EcmaArray:
		{

		}
	case RTMP_AMF0_ObjectEnd:
		{

		}
	case RTMP_AMF0_StrictArray:
		{

		}
	case RTMP_AMF0_Date:
		{

		}
	case RTMP_AMF0_LongString:
		{

		}
	case RTMP_AMF0_UnSupported:
		{

		}
	case RTMP_AMF0_RecordSet:
		{

		}
	case RTMP_AMF0_XmlDocument:
		{

		}
	case RTMP_AMF0_TypedObject:
		{

		}
	// AVM+ object is the AMF3 object.
	case RTMP_AMF0_AVMplusObject:
		{

		}
	// origin array whos data takes the same form as LengthValueBytes
	case RTMP_AMF0_OriginStrictArray:
		{

		}
	// User defined
	case RTMP_AMF0_Invalid:
		{

		}
	}
	return
}

func encodeAmf0(stream *SrsStream, v interface{}) (err error) {
	switch v.(type) {
	case float64:
		{
			err = srs_amf0_write_number(stream, v.(float64))
			break
		}
	case bool:
		{
			// v, err = srs_amf0_read_number(stream)
			break
		}
	case string:
		{
			err = srs_amf0_write_string(stream, v.(string))
			break
		}
	case SrsAmf0Object:
		{
			err = encodeAmf0(stream, v)
			break
		}
	case nil:
		{

		}
	}
	return
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
	return s.marker == RTMP_AMF0_String
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
