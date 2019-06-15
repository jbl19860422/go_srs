package amf0

import (
	"errors"
	"log"
)

type SrsAmf0Object struct {
	properties []SrsValuePair
	eof        *SrsAmf0ObjectEOF
}

func (this *SrsAmf0Object) Decode(stream *utils.SrsStream) error {
	marker, err := stream.ReadByte();
	if err != nil {
		return err
	}

	if marker != RTMP_AMF0_Object {
		err = errors.New("amf0 check object marker failed. ")
		return err
	}

	for {
		if is_eof, err := this.eof.IsMyType(stream); err != nil {
			return err
		}

		if is_eof {
			this.eof.Decode(stream)
			return nil
		}
		//读取属性名称
		var pname SrsAmf0Utf8
		err = pname.Decode(stream)
		if err != nil {
			return err
		}
		
		marker, err := stream.PeekByte(1)
		if err != nil {
			return err
		}

		var v SrsAmf0Any
		switch marker {
		case RTMP_AMF0_Number:{
			v := SrsAmf0Number{}
			v.Decode(stream)
		}
		case RTMP_AMF0_Boolean:{
			
		}
		}

	}
}

func (this *SrsAmf0Object) GetStringProperty(key string) (s string, err error) {
	for i := 0; i < len(this.properties); i++ {
		if this.properties[i].name == key {
			s = this.properties[i].val.(string)
		}
	}
	return
}

func (this *SrsAmf0Object) GetNumberProperty(key string) (s float64, err error) {
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
