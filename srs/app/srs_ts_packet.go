package app

import (
	"go_srs/srs/utils"
)

type SrsTsPacket struct {
	tsHeader        *SrsTsHeader
	adaptationField *SrsTsAdapationField
	payload		[]byte
}

func NewSrsTsPacket() *SrsTsPacket {
	return &SrsTsPacket{
		tsHeader: NewSrsTsHeader(),
	}
}

func (this *SrsTsPacket) Decode(stream *utils.SrsStream) error {
	return nil
}

func (this *SrsTsPacket) Encode(stream *utils.SrsStream) {
	this.tsHeader.Encode(stream)//4
	if this.tsHeader.adaptationFieldControl == SrsTsAdapationControlFieldOnly || this.tsHeader.adaptationFieldControl == SrsTsAdapationControlBoth {
		this.adaptationField.Encode(stream)
	}
	stream.WriteBytes(this.payload)
}
