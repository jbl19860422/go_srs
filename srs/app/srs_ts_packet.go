package app

import (
	"fmt"
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
	this.tsHeader.Encode(stream)//4
	fmt.Println("len_ts_header=", len(stream.Data()))
	if this.tsHeader.adaptationFieldControl == SrsTsAdapationControlFieldOnly || this.tsHeader.adaptationFieldControl == SrsTsAdapationControlBoth {
		fmt.Println("rrrrrrthis.adaptationField encode")
		this.adaptationField.Encode(stream)
	}
	fmt.Println("len_adaptation=", len(stream.Data()))
	//this.payload.Encode(stream)
	stream.WriteBytes(this.payload1)
	fmt.Println("len_total=", len(stream.Data()))
}
