package protocol

import (
	"log"
)
/**
* 2.5 Object Type
* anonymous-object-type = object-marker *(object-property)
* object-property = (UTF-8 value-type) | (UTF-8-empty object-end-marker)
*/

type SrsAmf0Object struct {
	properties map[string]SrsAmf0Any
	eof *SrsAmf0ObjectEOF
}

func (this *SrsAmf0Object)read(stream *SrsStream) (err error) {
	var marker int8
	if marker, err = stream.read_int8(); err != nil {
		return
	}

	if marker != RTMP_AMF0_Object {
		err = errors.New("amf0 check object marker failed. ")
		return
	}

	for !stream.empty() {
		// detect whether is eof.
		if srs_amf0_is_object_eof(stream) {
			var pbj_eof SrsAmf0ObjectEOF
			if err = pbj_eof.read(stream); err != nil {
				return
			}
			break
		}
		// property-name: utf8 string
		var property_name string
		if property_name, err = srs_amf0_read_utf8(stream); err != nil {
			log.Print("amf0 object read property name failed")
			return
		}
		// property-value: any

	}
	return
}

func (this *SrsAmf0Object)decode(stream *SrsStream) (err error) {
	var marker int8
	if marker, err = stream.read_int8(); err != nil {
		return
	}

	switch marker {
	case RTMP_AMF0_Number: {
		
	}
	case RTMP_AMF0_Boolean: {
		
	}
	case RTMP_AMF0_String: {
		
	}
	case RTMP_AMF0_Object: {

	}
	case RTMP_AMF0_MovieClip: {

	}
	case RTMP_AMF0_Null: {

	}
	case RTMP_AMF0_Undefined: {

	}
	case RTMP_AMF0_Reference: {

	}
	case RTMP_AMF0_EcmaArray: {

	}
	case RTMP_AMF0_ObjectEnd: {

	}
	case RTMP_AMF0_StrictArray: {

	}
	case RTMP_AMF0_Date: {

	}

	case RTMP_AMF0_LongString: {

	}

	case RTMP_AMF0_UnSupported: {

	}

	case RTMP_AMF0_RecordSet: {

	}
	case RTMP_AMF0_XmlDocument: {

	}

	case RTMP_AMF0_TypedObject: {

	}
	// AVM+ object is the AMF3 object.
	case RTMP_AMF0_AVMplusObject: {

	}
		// origin array whos data takes the same form as LengthValueBytes
	case RTMP_AMF0_OriginStrictArray: {

	}
		// User defined
	case RTMP_AMF0_Invalid: {
		
	}
	}
}
