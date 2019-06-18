package packet

import (
	"errors"
	"go_srs/srs/protocol/amf0"
	"go_srs/srs/utils"
	"go_srs/srs/global"
	// "fmt"
)
type SrsOnMetaDataPacket struct {
	Name 		amf0.SrsAmf0String
	MetaData	amf0.SrsAmf0Any
	// OMetaData 	*amf0.SrsAmf0Object
	// AMetaData	*amf0.SrsAmf0EcmaArray

	// IsObjMeta	bool
}

func NewSrsOnMetaDataPacket(command string) *SrsOnMetaDataPacket {
	return &SrsOnMetaDataPacket{
		Name:amf0.SrsAmf0String{Value:amf0.SrsAmf0Utf8{Value:command}},
		// OMetaData:amf0.NewSrsAmf0Object(),
		// AMetaData:amf0.NewSrsAmf0EcmaArray(),
		// IsObjMeta:true,
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
			this.MetaData = amf0.GenerateSrsAmf0Any(marker)
		}
		case amf0.RTMP_AMF0_EcmaArray:{
			this.MetaData = amf0.GenerateSrsAmf0Any(marker)
		}
	}

	if this.MetaData != nil {
		if err = this.MetaData.Decode(stream); err != nil {
			return err
		}
	}
	
	return nil
}

func (this *SrsOnMetaDataPacket) Set(name string, value interface{}) error {
	if this.MetaData == nil {
		this.MetaData =  amf0.GenerateSrsAmf0Any(amf0.RTMP_AMF0_EcmaArray)
	}

	switch this.MetaData.(type) {
		case *amf0.SrsAmf0Object: {
			this.MetaData.(*amf0.SrsAmf0Object).Set(name, value)
		}
		case *amf0.SrsAmf0EcmaArray: {
			this.MetaData.(*amf0.SrsAmf0EcmaArray).Set(name, value)
		}
	}
	return nil
}

func (this *SrsOnMetaDataPacket) Get(name string, value interface{}) error {
	if this.MetaData == nil {
		return errors.New("metadata is nil")
	}

	switch this.MetaData.(type) {
		case *amf0.SrsAmf0Object: {
			return this.MetaData.(*amf0.SrsAmf0Object).Get(name, value)
		}
		case *amf0.SrsAmf0EcmaArray: {
			return this.MetaData.(*amf0.SrsAmf0EcmaArray).Get(name, value)
		}
	}
	return errors.New("metadata's type error")
}

func (this *SrsOnMetaDataPacket) Encode(stream *utils.SrsStream) error {
	_ = this.Name.Encode(stream)
	_ = this.MetaData.Encode(stream)
	return nil
}