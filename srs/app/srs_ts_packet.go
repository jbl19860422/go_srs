package app

import (
	"go_srs/srs/utils"
)

type SrsTsPacket struct {
	tsHeader        *SrsTsHeader
	adaptationField *SrsTsAdapationField
	payload         SrsTsPayload
	payload1		[]byte
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
	this.tsHeader.Encode(stream)
	if this.tsHeader.adaptationFieldControl == SrsTsAdapationControlFieldOnly || this.tsHeader.adaptationFieldControl == SrsTsAdapationControlBoth {
		this.adaptationField.Encode(stream)
	}

	this.payload.Encode(stream)
}
