/*
The MIT License (MIT)

Copyright (c) 2013-2015 GOSRS(gosrs)

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

package app

import (
	"os"
	"errors"
	"encoding/binary"
	"go_srs/srs/protocol/amf0"
	"go_srs/srs/protocol/rtmp"
	"go_srs/srs/global"
	"go_srs/srs/utils"
	"go_srs/srs/codec/flv"
)

type SrsFlvSegment struct {
	path 			string
	flvEncoder		*flvcodec.SrsFlvEncoder
	durationOffset	int64
	filesizeOffset	int64
	startTime		int64
	previousPktTime int64
	duration		int64
	streamDuration	int64
	file			*os.File
}

func NewSrsFlvSegment(fname string) *SrsFlvSegment {
	f, err := os.OpenFile(fname, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return nil
	}
	f.Truncate(0)
	return &SrsFlvSegment{
		path:fname,
		flvEncoder:flvcodec.NewSrsFlvEncoder(f),
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
			metaData.(*amf0.SrsAmf0Object).Remove("fileSize")
			metaData.(*amf0.SrsAmf0Object).Remove("framerate")
			metaData.(*amf0.SrsAmf0Object).Set("service", global.RTMP_SIG_SRS_SERVER)
			metaData.(*amf0.SrsAmf0Object).Set("filesize", float64(0))
			metaData.(*amf0.SrsAmf0Object).Set("duration", float64(0))
		}
		case amf0.RTMP_AMF0_EcmaArray:{
			metaData.(*amf0.SrsAmf0EcmaArray).Remove("fileSize")
			metaData.(*amf0.SrsAmf0EcmaArray).Remove("framerate")
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
	this.durationOffset = off + int64(size) + 11 - 3 - 8
	this.filesizeOffset = this.durationOffset - 1 - (2 + int64(len("duration"))) - 8
	_, err = this.flvEncoder.WriteMetaData(writeStream.Data())
	return err
}

func (this *SrsFlvSegment) onUpdateDuration(msg *rtmp.SrsRtmpMessage) error {
	if this.startTime < 0 {
		this.startTime = msg.GetHeader().GetTimestamp()
	}

	if this.previousPktTime < 0 || this.previousPktTime > msg.GetHeader().GetTimestamp() {
		this.previousPktTime = msg.GetHeader().GetTimestamp()
	}
	this.duration += msg.GetHeader().GetTimestamp() - this.previousPktTime
	this.streamDuration += msg.GetHeader().GetTimestamp() - this.previousPktTime
	this.previousPktTime = msg.GetHeader().GetTimestamp()
	return nil
}

func (this *SrsFlvSegment) WriteAudio(msg *rtmp.SrsRtmpMessage) error {
	this.flvEncoder.WriteAudio(uint32(msg.GetHeader().GetTimestamp()), msg.GetPayload())
	// this.onUpdateDuration(msg)
	return nil
}

func (this *SrsFlvSegment) WriteVideo(msg *rtmp.SrsRtmpMessage) error {
	this.flvEncoder.WriteVideo(uint32(msg.GetHeader().GetTimestamp()), msg.GetPayload())
	this.onUpdateDuration(msg)
	return nil
}

func (this *SrsFlvSegment) Close() error {
	this.updateMetaData()
	return nil
}

func (this *SrsFlvSegment) updateMetaData() error {
	off, _ := this.file.Seek(0, 2)
	c := utils.Float64ToBytes(float64(off), binary.BigEndian)
	this.file.WriteAt(c, this.filesizeOffset)
	b := utils.Float64ToBytes(float64(this.duration)/1000, binary.BigEndian)
	this.file.WriteAt(b, this.durationOffset)

	this.file.Close()
	return nil
}