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
	// properties map[string]interface{}
	properties []SrsValuePair
	eof        *SrsAmf0ObjectEOF
}

func (this *SrsAmf0Object) GetStringProperty(key string) (s string, err error) {
	// v, ok := this.properties[key]
	// if !ok {
	// 	log.Print("no property ", key)
	// 	err = errors.New("property " + key + " not exist")
	// 	return
	// }

	// s, ok1 := v.(string)
	// if !ok1 {
	// 	err = errors.New("not string type")
	// 	return
	// }
	// return
	for i := 0; i < len(this.properties); i++ {
		if this.properties[i].name == key {
			s = this.properties[i].val.(string)
		}
	}
	return
}

func (this *SrsAmf0Object) GetNumberProperty(key string) (s float64, err error) {
	// v, ok := this.properties[key]
	// if !ok {
	// 	log.Print("no property ", key)
	// 	err = errors.New("property " + key + " not exist")
	// 	return
	// }

	// s, ok1 := v.(float64)
	// if !ok1 {
	// 	err = errors.New("not string type")
	// 	return
	// }
	// return
	for i := 0; i < len(this.properties); i++ {
		if this.properties[i].name == key {
			s = this.properties[i].val.(float64)
		}
	}
	return
}

func (this *SrsAmf0Object) Set(key string, v interface{}) {
	p := SrsValuePair{
		name:key,
		val:v,
	}
	this.properties = append(this.properties, p)
	return
}

func (this *SrsAmf0Object) SetStringProperty(key string, v string) {
	p := SrsValuePair{
		name:key,
		val:v,
	}
	this.properties = append(this.properties, p)
	return
}

func (this *SrsAmf0Object) SetNumberProperty(key string, v float64) {
	p := SrsValuePair{
		name:key,
		val:v,
	}
	this.properties = append(this.properties, p)
	// this.properties[key] = v
	return
}

func NewSrsAmf0Object() *SrsAmf0Object {
	s := &SrsAmf0Object{eof: &SrsAmf0ObjectEOF{}}
	s.properties = make([]SrsValuePair, 0)
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
		val1, err1 := decodeAmf0(stream)
		if err1 != nil {
			err = err1
			return
		}

		p := SrsValuePair{
			name:property_name,
			val:val1,
		}
		this.properties = append(this.properties, p)

		// this.properties[property_name] = val
		log.Print("properties len=", len(this.properties))
	}
	return
}

func (this *SrsAmf0Object) write(stream *SrsStream) error {
	stream.write_1byte(RTMP_AMF0_Object)
	// value
	for i := 0; i < len(this.properties); i++ {
		srs_amf0_write_utf8(stream, this.properties[i].name)
		encodeAmf0(stream, this.properties[i].val)
	}
	// for k, v := range this.properties {
	// 	srs_amf0_write_utf8(stream, k)
	// 	log.Print("encodeAmf0 k=", k)
	// 	encodeAmf0(stream, v)
	// }

	_ = this.eof.write(stream)
	return nil
}

func (s *SrsAmf0Object) total_size() int {
	// var size int = 1
	// var sz SrsAmf0Size = SrsAmf0Size{}
	// for 
	// for k, v := range s.properties {
	// 	size += sz.utf8(k)
	// 	size += sz.any(v)
	// }

	// size += sz.object_eof()
	// return size
	return 0
}
