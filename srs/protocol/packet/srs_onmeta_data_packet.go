package packet

import (
	"go_srs/srs/protocol/amf0"
	"go_srs/srs/utils"
	"go_srs/srs/global"
	// "fmt"
)
type SrsOnMetaDataPacket struct {
	Name 		amf0.SrsAmf0String
	OMetaData 	*amf0.SrsAmf0Object
	AMetaData	*amf0.SrsAmf0EcmaArray
	IsObjMeta	bool
}

func NewSrsOnMetaDataPacket(command string) *SrsOnMetaDataPacket {
	return &SrsOnMetaDataPacket{
		Name:amf0.SrsAmf0String{Value:amf0.SrsAmf0Utf8{Value:command}},
		OMetaData:amf0.NewSrsAmf0Object(),
		AMetaData:amf0.NewSrsAmf0EcmaArray(),
		IsObjMeta:true,
	}
}

func (s *SrsOnMetaDataPacket) GetMessageType() int8 {
	return global.RTMP_MSG_AMF0DataMessage
}

func (s *SrsOnMetaDataPacket) GetPreferCid() int32 {
	return global.RTMP_CID_OverConnection2
}

func (this *SrsOnMetaDataPacket) Decode(stream *utils.SrsStream) error {
	var err error
	if this.Name.GetValue().(string) == amf0.SRS_CONSTS_RTMP_SET_DATAFRAME {
		if err = this.Name.Decode(stream); err != nil {
			return err
		}
	}
	
	marker, err2 := stream.PeekByte()
	if err2 != nil {
		return err2
	}

	switch marker {
	case amf0.RTMP_AMF0_Object:{
		this.IsObjMeta = true
		if err = this.OMetaData.Decode(stream); err != nil {
			return err
		}
	}
	case amf0.RTMP_AMF0_EcmaArray:{
		// fmt.Println("xxxxxxxxxxxxxxxxxxxxxxxx is arrayxxxxxxxxxxxxxxxxxxxxx")
		this.IsObjMeta = false
		if err = this.AMetaData.Decode(stream); err != nil {
			return err
		}
	}
	}
	return nil
}

func (this *SrsOnMetaDataPacket) Encode(stream *utils.SrsStream) error {
	_ = this.Name.Encode(stream)
	
	if this.IsObjMeta {
		_ = this.OMetaData.Encode(stream)
	} else {
		_ = this.AMetaData.Encode(stream)
	}
	return nil
}