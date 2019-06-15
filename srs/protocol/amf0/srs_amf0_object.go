package amf0

import (
	"errors"
	"log"
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
		var pname SrsAmf0Utf8 = SrsAmf0Utf8{}
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
				v = SrsAmf0Number{}
				err = v.Decode(stream)
			}
			case RTMP_AMF0_Boolean:{
				v = SrsAmf0Boolean{}
				err = v.Decode(stream)
			}
			case RTMP_AMF0_String:{
				v = SrsAmf0String{}
				err = v.Decode(stream)
			}
			case RTMP_AMF0_Object: {
				v = SrsAmf0Object{}
				err = v.Decode(stream)
			}
			case RTMP_AMF0_Null:{
				v = SrsAmf0Null{}
				err = v.Decode(stream)
			}
			case RTMP_AMF0_Undefined:{
				v = SrsAmf0Undefined{}
				err = v.Decode(stream)
			}
		}

		if err != nil {
			return err
		}

		pair := SrsValuePair{
			name:pname,
			value:v,
		}
		this.Properties = append(this.Properties, pair)
	}
	return nil
}

func (this *SrsAmf0Object) Encode(stream *utils.SrsStream) error {
	stream.WriteByte(RTMP_AMF0_Object)
	for i := 0; i < len(this.Properties); i++ {
		_ = this.Properties[i].name.Encode(stream)
		_ = this.Properties[i].value.Encode(stream)
	}
	_ = this.eof.Encode(stream)
	return nil
}

func (this *SrsAmf0Object) IsMyType(stream *utils.SrsStream) (bool, error) {
	marker, err := stream.PeekByte()
	if err != nil {
		return err
	}

	if marker != RTMP_AMF0_Object {
		return false, nil
	}
	return true, nil
}
