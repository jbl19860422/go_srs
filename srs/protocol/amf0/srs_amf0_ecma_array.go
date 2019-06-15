package amf0

import (
	_ "log"
)

type SrsAmf0EcmaArray struct {
	properties 	[]SrsValuePair
	eof        	*SrsAmf0ObjectEOF
	count		int32
}

func NewSrsAmf0EcmaArray() *SrsAmf0EcmaArray {
	s := &SrsAmf0EcmaArray{eof: &SrsAmf0ObjectEOF{}, count:0}
	s.properties = make([]SrsValuePair, 0)
	return s
}

func (this *SrsAmf0EcmaArray) Count() int {
	return len(this.properties)
}

func (this *SrsAmf0EcmaArray) Clear() {
	this.properties = this.properties[0:0]
}

func (this *SrsAmf0EcmaArray) KeyAt(i int) string {
	if i < len(this.properties) {
		return this.properties[i].name
	}
	return ""
}

func (this *SrsAmf0EcmaArray) ValueAt(i int) SrsAmf0Any {
	if i < len(this.properties) {
		return this.properties[i].value
	}
	return nil
}

func (this *SrsAmf0EcmaArray) Set(key string, v SrsAmf0Any) {
	pair := SrsValuePair{
		name:key,
		value:v,
	}
	this.properties = append(this.properties, pair)
	return
}

func (this *SrsAmf0EcmaArray) Decode(stream *utils.SrsStream) error {
	marker, err := stream.ReadByte();
	if err != nil {
		return err
	}

	if marker != RTMP_AMF0_EcmaArray {
		err = errors.New("amf0 check ecma array marker failed. ")
		return err
	}

	count, err := stream.ReadInt32()
	if err != nil {
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
		this.properties = append(this.properties, pair)
	}
	return nil
}

func (this *SrsAmf0EcmaArray) Encode(stream *utils.SrsStream) error {
	stream.WriteByte(RTMP_AMF0_EcmaArray)
	_ = stream.WriteInt32(0, binary.BigEndian)
	for i := 0; i < len(this.properties); i++ {
		_ = this.properties[i].name.Encode(stream)
		_ = this.properties[i].value.Encode(stream)
	}
	_ = this.eof.Encode(stream)
	return nil
}