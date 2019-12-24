/*
The MIT License (MIT)

Copyright (c) 2019 GOSRS(gosrs)

Permission is hereby granted, free of charge, to any person obtaining a copy of
this software and associated documentation files (the "Software"), to deal in
the Software without restriction, including without limitation the rights to
use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
the Software, and to permit persons to whom the Software is furnished to do so,
subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/
package amf0

import (
	"encoding/binary"
	"errors"
	"go_srs/srs/utils"
	"reflect"
	_ "log"
	"fmt"
)

type SrsAmf0EcmaArray struct {
	Properties []SrsValuePair
	eof        *SrsAmf0ObjectEOF
	count      int32
}

func NewSrsAmf0EcmaArray() *SrsAmf0EcmaArray {
	s := &SrsAmf0EcmaArray{eof: &SrsAmf0ObjectEOF{}, count: 0}
	s.Properties = make([]SrsValuePair, 0)
	return s
}

func (this *SrsAmf0EcmaArray) Count() int {
	return len(this.Properties)
}

func (this *SrsAmf0EcmaArray) Clear() {
	this.Properties = this.Properties[0:0]
}

func (this *SrsAmf0EcmaArray) KeyAt(i int) string {
	if i < len(this.Properties) {
		return this.Properties[i].Name.Value
	}
	return ""
}

func (this *SrsAmf0EcmaArray) ValueAt(i int) SrsAmf0Any {
	if i < len(this.Properties) {
		return this.Properties[i].Value
	}
	return nil
}

func (this *SrsAmf0EcmaArray) Decode(stream *utils.SrsStream) error {
	marker, err := stream.ReadByte()
	if err != nil {
		return err
	}

	if marker != RTMP_AMF0_EcmaArray {
		err = errors.New("amf0 check ecma array marker failed. ")
		return err
	}

	this.count, err = stream.ReadInt32(binary.BigEndian)
	if err != nil {
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

func (this *SrsAmf0EcmaArray) Encode(stream *utils.SrsStream) error {
	stream.WriteByte(byte(RTMP_AMF0_EcmaArray))
	stream.WriteInt32(int32(len(this.Properties)), binary.BigEndian)
	for i := 0; i < len(this.Properties); i++ {
		_ = this.Properties[i].Name.Encode(stream)
		_ = this.Properties[i].Value.Encode(stream)
	}
	_ = this.eof.Encode(stream)
	return nil
}

func (this *SrsAmf0EcmaArray) IsMyType(stream *utils.SrsStream) (bool, error) {
	marker, err := stream.PeekByte()
	if err != nil {
		return false, err
	}

	if marker == RTMP_AMF0_EcmaArray {
		return true, nil
	}
	return false, nil
}

func (this *SrsAmf0EcmaArray) Remove(name string) {
	for i := 0; i < len(this.Properties); i++ {
		if this.Properties[i].Name.Value == name {
			this.Properties = append(this.Properties[0:i], this.Properties[i+1:]...)
		}
	}
}

func (this *SrsAmf0EcmaArray) Set(name string, value interface{}) {
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
	}
	this.Properties = append(this.Properties, *p)
}

func (this *SrsAmf0EcmaArray) Get(name string, pval interface{}) error {
	if reflect.TypeOf(pval).Kind() != reflect.Ptr {
		return errors.New("need pointer to get value")
	}

	for i := 0; i < len(this.Properties); i++ {
		fmt.Println(this.Properties[i].Name.Value, name)
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

func (this *SrsAmf0EcmaArray) GetValue() interface{} {
	return this.Properties
}
