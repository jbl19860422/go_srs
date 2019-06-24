package amf0

import (
	"errors"
	"go_srs/srs/utils"
	"reflect"
	// "fmt"
)

type SrsAmf0Object struct {
	Properties []SrsValuePair
	eof        *SrsAmf0ObjectEOF
}

func NewSrsAmf0Object() *SrsAmf0Object {
	s := &SrsAmf0Object{eof: &SrsAmf0ObjectEOF{}}
	s.Properties = make([]SrsValuePair, 0)
	return s
}

func (this *SrsAmf0Object) Decode(stream *utils.SrsStream) error {
	marker, err := stream.ReadByte()
	if err != nil {
		return err
	}

	if marker != RTMP_AMF0_Object {
		err = errors.New("amf0 check object marker failed. ")
		return err
	}

	for {
		var is_eof bool
		if is_eof, err = this.eof.IsMyType(stream); err != nil {
			return err
		}

		if is_eof {
			this.eof.Decode(stream)
			return nil
		}
		//读取属性名称
		var pname SrsAmf0Utf8 = SrsAmf0Utf8{}
		err = pname.Decode(stream)
		if err != nil {
			return err
		}

		marker, err := stream.PeekByte()
		if err != nil {
			return err
		}

		var v SrsAmf0Any
		switch marker {
		case RTMP_AMF0_Number:
			{
				v = &SrsAmf0Number{}
				err = v.Decode(stream)
			}
		case RTMP_AMF0_Boolean:
			{
				v = &SrsAmf0Boolean{}
				err = v.Decode(stream)
			}
		case RTMP_AMF0_String:
			{
				v = &SrsAmf0String{}
				err = v.Decode(stream)
			}
		case RTMP_AMF0_Object:
			{
				v = &SrsAmf0Object{}
				err = v.Decode(stream)
			}
		case RTMP_AMF0_Null:
			{
				v = &SrsAmf0Null{}
				err = v.Decode(stream)
			}
		case RTMP_AMF0_Undefined:
			{
				v = &SrsAmf0Undefined{}
				err = v.Decode(stream)
			}
		}

		if err != nil {
			return err
		}

		pair := SrsValuePair{
			Name:  pname,
			Value: v,
		}
		this.Properties = append(this.Properties, pair)
	}
	return nil
}

func (this *SrsAmf0Object) Encode(stream *utils.SrsStream) error {
	stream.WriteByte(RTMP_AMF0_Object)
	for i := 0; i < len(this.Properties); i++ {
		_ = this.Properties[i].Name.Encode(stream)
		_ = this.Properties[i].Value.Encode(stream)
	}
	_ = this.eof.Encode(stream)
	return nil
}

func (this *SrsAmf0Object) IsMyType(stream *utils.SrsStream) (bool, error) {
	marker, err := stream.PeekByte()
	if err != nil {
		return false, err
	}

	if marker != RTMP_AMF0_Object {
		return false, nil
	}
	return true, nil
}

func (this *SrsAmf0Object) Set(name string, value interface{}) {
	this.Remove(name)
	var p *SrsValuePair
	switch value.(type) {
	case string:
		p = &SrsValuePair{
			Name:  SrsAmf0Utf8{Value: name},
			Value: &SrsAmf0String{Value: SrsAmf0Utf8{Value: value.(string)}},
		}
	case bool:
		p = &SrsValuePair{
			Name:  SrsAmf0Utf8{Value: name},
			Value: &SrsAmf0Boolean{Value: value.(bool)},
		}
	case float64:
		p = &SrsValuePair{
			Name:  SrsAmf0Utf8{Value: name},
			Value: &SrsAmf0Number{Value: value.(float64)},
		}
	case *SrsAmf0Object:
		p = &SrsValuePair{
			Name:  SrsAmf0Utf8{Value: name},
			Value: value.(*SrsAmf0Object),
		}
	case *SrsAmf0EcmaArray:
		p = &SrsValuePair{
			Name:  SrsAmf0Utf8{Value: name},
			Value: value.(*SrsAmf0EcmaArray),
		}
	}

	this.Properties = append(this.Properties, *p)
}

func (this *SrsAmf0Object) Remove(name string) {
	for i := 0; i < len(this.Properties); i++ {
		if this.Properties[i].Name.Value == name {
			this.Properties = append(this.Properties[0:i], this.Properties[i+1:]...)
		}
	}
}

func (this *SrsAmf0Object) Get(name string, pval interface{}) error {
	if reflect.TypeOf(pval).Kind() != reflect.Ptr {
		return errors.New("need pointer to get value")
	}
	for i := 0; i < len(this.Properties); i++ {
		if this.Properties[i].Name.Value == name {
			if reflect.TypeOf(pval).Elem() == reflect.TypeOf(this.Properties[i].Value.GetValue()) {
				reflect.ValueOf(pval).Elem().Set(reflect.ValueOf(this.Properties[i].Value.GetValue()))
				return nil
			} else {
				return errors.New("type not match")
			}
		}
	}
	return errors.New("could not find key:" + name)
}

func (this *SrsAmf0Object) GetValue() interface{} {
	return this.Properties
}


