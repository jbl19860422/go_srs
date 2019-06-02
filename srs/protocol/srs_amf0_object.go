package protocol

import (
	"errors"
	"log"
)

/**
* 2.5 Object Type
* anonymous-object-type = object-marker *(object-property)
* object-property = (UTF-8 value-type) | (UTF-8-empty object-end-marker)
 */

type SrsAmf0Object struct {
	properties map[string]interface{}
	eof        *SrsAmf0ObjectEOF
}

func (this *SrsAmf0Object) GetStringProperty(key string) (s string, err error) {
	log.Print(this.properties)
	v, ok := this.properties[key]
	if !ok {
		log.Print("no property ", key)
		err = errors.New("property " + key + " not exist")
		return
	}

	s, ok1 := v.(string)
	if !ok1 {
		err = errors.New("not string type")
		return
	} else {
		log.Print("xxxxxxxxxxxxxxxxxxxxxxxi string")
	}
	return
}

func (this *SrsAmf0Object) GetNumberProperty(key string) (s float64, err error) {
	log.Print(this.properties)
	v, ok := this.properties[key]
	if !ok {
		log.Print("no property ", key)
		err = errors.New("property " + key + " not exist")
		return
	}

	s, ok1 := v.(float64)
	if !ok1 {
		err = errors.New("not string type")
		return
	} else {
		log.Print("xxxxxxxxxxxxxxxxxxxxxxxi string")
	}
	return
}

func NewSrsAmf0Object() *SrsAmf0Object {
	s := &SrsAmf0Object{}
	s.properties = make(map[string]interface{})
	return s
}

func (this *SrsAmf0Object) read(stream *SrsStream) (err error) {
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
		log.Print("propername=", property_name)
		// property-value: any
		val, err1 := decodeAmf0(stream)
		if err1 != nil {
			err = err1
			return
		}
		this.properties[property_name] = val
		log.Print("properties len=", len(this.properties))
	}
	return
}
