package protocol

import (
	_ "log"
)

type SrsAmf0EcmaArray struct {
	properties []SrsValuePair
	eof        *SrsAmf0ObjectEOF
}

func NewSrsAmf0EcmaArray() *SrsAmf0EcmaArray {
	s := &SrsAmf0EcmaArray{eof: &SrsAmf0ObjectEOF{}}
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

func (this *SrsAmf0EcmaArray) ValueAt(i int) interface{} {
	if i < len(this.properties) {
		return this.properties[i].val
	}
	return nil
}

func (this *SrsAmf0EcmaArray) Set(key string, v interface{}) {
	p := SrsValuePair{
		name:key,
		val:v,
	}
	this.properties = append(this.properties, p)
	return
}

func (this *SrsAmf0EcmaArray) read(stream *SrsStream) (err error) {
	return
}

func (this *SrsAmf0EcmaArray) write(stream *SrsStream) error {
	stream.write_1byte(RTMP_AMF0_EcmaArray)
	stream.write_int32(0)
	// value
	for i := 0; i < len(this.properties); i++ {
		name := this.KeyAt(i)
		val := this.ValueAt(i)
		srs_amf0_write_utf8(stream, name)
		switch val.(type) {
			case string:{
				srs_amf0_write_string(stream, val.(string))
			}
			case float64:{
				srs_amf0_write_number(stream, val.(float64))
			}
		}
	}

	_ = this.eof.write(stream)
	return nil
}