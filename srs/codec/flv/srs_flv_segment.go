package flvcodec

import (
	"os"
	"fmt"
	"errors"
	"go_srs/srs/protocol/amf0"
	"go_srs/srs/protocol/rtmp"
	"go_srs/srs/global"
	"go_srs/srs/utils"
)

type SrsFlvSegment struct {
	path 			string
	flvEncoder		*SrsFlvEncoder
	durationOffset	int64
	filesizeOffset	int64
	startTime		int64
	previousPktTime int64
	duration		int64
	streamDuration	int64
	file			*os.File
}

func NewSrsFlvSegment(fname string) *SrsFlvSegment {
	fmt.Println("**********************open ", fname, "********************")
	f, err := os.OpenFile(fname, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		fmt.Println("**********************open ", fname, " failed********************", err)
		return nil
	}
	return &SrsFlvSegment{
		path:"./record.flv",
		flvEncoder:NewSrsFlvEncoder(f),
		startTime:-1,
		previousPktTime:-1,
		duration:0,
		streamDuration:0,
		file:f,
	}
}

func (this *SrsFlvSegment) Initialize() {
	_ = this.flvEncoder.WriteHeader()
}

func (this *SrsFlvSegment) WriteMetaData(msg *rtmp.SrsRtmpMessage) error {
	stream := utils.NewSrsStream(msg.GetPayload())

	var command amf0.SrsAmf0String
	if err := command.Decode(stream); err != nil {
		return err
	}

	var name amf0.SrsAmf0String
	if err := name.Decode(stream); err != nil {
		return err
	}
	
	marker, err := stream.PeekByte()
	if err != nil {
		return err
	}
	// fmt.Println("xxxxxxxxxxxxxxxxxxxxx command=", command.GetValue().(string), ",name=", name.GetValue().(string)," xxxxxxxxxxxxxxxx", marker)

	var metaData amf0.SrsAmf0Any
	switch marker {
		case amf0.RTMP_AMF0_Object:{
			metaData = amf0.GenerateSrsAmf0Any(marker)
		}
		case amf0.RTMP_AMF0_EcmaArray:{
			metaData = amf0.GenerateSrsAmf0Any(marker)
		}
		default:{
			return errors.New("error marker")
		}
	}

	if metaData != nil {
		if err = metaData.Decode(stream); err != nil {
			return err
		}
	}

	switch marker {
		case amf0.RTMP_AMF0_Object:{
			metaData.(*amf0.SrsAmf0Object).Set("service", global.RTMP_SIG_SRS_SERVER)
			metaData.(*amf0.SrsAmf0Object).Set("filesize", float64(0))
			metaData.(*amf0.SrsAmf0Object).Set("duration", float64(0))
		}
		case amf0.RTMP_AMF0_EcmaArray:{
			metaData.(*amf0.SrsAmf0EcmaArray).Set("service", global.RTMP_SIG_SRS_SERVER)
			metaData.(*amf0.SrsAmf0EcmaArray).Set("filesize", float64(0))
			metaData.(*amf0.SrsAmf0EcmaArray).Set("duration", float64(0))
		}
		default:{
			return errors.New("error marker")
		}
	}
	
	writeStream := utils.NewSrsStream([]byte{})
	_ = name.Encode(writeStream)
	_ = metaData.Encode(writeStream)
	size := len(writeStream.Data())
	off, err := this.file.Seek(0, 1)//SEEK_CUR
	// 11B flv tag header, 3B object EOF, 8B number value, 1B number flag.
	//todo fix me, write readable code
	this.durationOffset = off + int64(size) + 11 - 3 - 9
	this.filesizeOffset = this.durationOffset - (2 + int64(len("duration"))) - 9

	_, err = this.flvEncoder.WriteMetaData(writeStream.Data())
	return err
}

func (this *SrsFlvSegment) onUpdateDuration(msg *rtmp.SrsRtmpMessage) error {
	if this.startTime < 0 {
		this.startTime = msg.GetTimestamp()
	}

	if this.previousPktTime < 0 || this.previousPktTime > msg.GetTimestamp() {
		this.previousPktTime = msg.GetTimestamp()
	}

	this.duration += msg.GetTimestamp() - this.previousPktTime
	this.streamDuration += msg.GetTimestamp() - this.previousPktTime
	return nil
}

func (this *SrsFlvSegment) WriteAudio(msg *rtmp.SrsRtmpMessage) error {
	this.flvEncoder.WriteAudio(uint32(msg.GetTimestamp()), msg.GetPayload())
	this.onUpdateDuration(msg)
	return nil
}

func (this *SrsFlvSegment) WriteVideo(msg *rtmp.SrsRtmpMessage) error {
	this.flvEncoder.WriteVideo(uint32(msg.GetTimestamp()), msg.GetPayload())
	this.onUpdateDuration(msg)
	return nil
}

func (this *SrsFlvSegment) Close() error {
	return nil
}

func (this *SrsFlvSegment) updateMetaData() error {
	return nil
}