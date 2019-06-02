package protocol

import "errors"

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
	if !s.require(1) {
		err = errors.New("amf0 read string marker failed")
		return
	}

	marker, err := s.read_nbytes(1)
	bmarker := marker[0]
	if int(bmarker) != RTMP_AMF0_String {
		err = errors.New("amf0 check string marker failed.")
	}

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

	if len <= 0 {
		err = errors.New("amf0 read empty string.")
		return
	}

	val, err = s.read_string(int32(len))
	return
}
