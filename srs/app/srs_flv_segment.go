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
	"go_srs/srs/app/config"
	"strings"
	"strconv"
	"fmt"
)

type SrsFlvSegment struct {
	path 			string
	req 			*SrsRequest
	flvEncoder		*flvcodec.SrsFlvEncoder
	durationOffset	int64
	filesizeOffset	int64
	startTime		int64
	previousPktTime int64
	duration		int64
	streamDuration	int64
	tmpFlvFile		string
	hasKeyFrame		bool
	jitter 			*SrsRtmpJitter
	file			*os.File
}

func NewSrsFlvSegment(r *SrsRequest) *SrsFlvSegment {
	return &SrsFlvSegment{
		req:r,
		startTime:-1,
		previousPktTime:-1,
		duration:0,
		streamDuration:0,
	}
}

func (this *SrsFlvSegment) Open(useTmpFile bool) error {
	if this.file != nil {
		return nil
	}

	this.path = this.generatePath()
	fmt.Println("*******************dvr file=", this.path, "******************")
	var freshFlvFile bool = false
	if _, err := os.Stat(this.path); os.IsExist(err) {
		freshFlvFile = false
	} else {
		freshFlvFile = true
	}

	if err := this.createJitter(!freshFlvFile); err != nil {
		return err
	}

	if !freshFlvFile || !useTmpFile {
		this.tmpFlvFile = this.path
	} else {
		this.tmpFlvFile = this.path + ".tmp"
	}

	var err error
	if !freshFlvFile {
		if this.file, err = os.OpenFile(this.tmpFlvFile, os.O_APPEND, 0755); err != nil {
			return err
		}
	} else {
		if this.file, err = os.OpenFile(this.tmpFlvFile, os.O_CREATE | os.O_RDWR, 0755); err != nil {
			return err
		}
	}

	this.flvEncoder = flvcodec.NewSrsFlvEncoder(this.file)
	if freshFlvFile {
		if err = this.flvEncoder.WriteHeader(); err != nil {
			return err
		}
	}
	return nil
}

func (this *SrsFlvSegment) generatePath() string {
	dvrPath := config.GetDvrPath(this.req.vhost)
	if !strings.Contains(dvrPath, ".flv") {
		dvrPath += "/[stream].[timestamp].flv"
	}

	flvPath := dvrPath
	flvPath = utils.Srs_path_build_stream(flvPath, this.req.vhost, this.req.app, this.req.stream)
	//todo build timestamp path
	flvPath = strings.Replace(flvPath, "[timestamp]", strconv.Itoa(int(utils.GetCurrentMs())), -1)
	return flvPath
}

func (this *SrsFlvSegment) createJitter(loadFromFlv bool) error {
	if !loadFromFlv {
		this.jitter = NewSrsRtmpJitter()

		this.startTime = -1
		this.previousPktTime = -1
		this.streamDuration = 0

		this.hasKeyFrame = false
		this.duration = 0
		return nil
	}
	// when jitter ok, do nothing.
	if this.jitter != nil {
		return nil
	}

	this.jitter = NewSrsRtmpJitter()
	return nil
}

func (this *SrsFlvSegment) Close() error {
	if this.file == nil {
		return nil
	}

	var err error
	if err = this.updateFlvMetaData(); err != nil {
		return err
	}

	if err = this.file.Close(); err != nil {
		return err
	}

	if this.tmpFlvFile != this.path {
		if err = os.Rename(this.tmpFlvFile, this.path); err != nil {
			return err
		}
	}

	return nil
}

func (this *SrsFlvSegment) WriteMetaData(msg *rtmp.SrsRtmpMessage) error {
	stream := utils.NewSrsStream(msg.GetPayload())

	var command amf0.SrsAmf0String
	if err := command.Decode(stream); err != nil {
		fmt.Println("111111111111111111")
		return err
	}

	var name amf0.SrsAmf0String
	if err := name.Decode(stream); err != nil {
		fmt.Println("2222222222222222222")
		return err
	}
	
	marker, err := stream.PeekByte()
	if err != nil {
		fmt.Println("33333333333333333333333")
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
			fmt.Println("4444444444444444444444444444")
			return errors.New("error marker")
		}
	}

	if metaData != nil {
		if err = metaData.Decode(stream); err != nil {
			fmt.Println("555555555555555555555555")
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
			fmt.Println("66666666666666666666666666")
			return errors.New("error marker")
		}
	}
	
	writeStream := utils.NewSrsStream([]byte{})
	if err = name.Encode(writeStream); err != nil {
		fmt.Println("777777777777777777777777")
		return err
	}

	if err = metaData.Encode(writeStream); err != nil {
		fmt.Println("888888888888888888888888")
		return err
	}

	size := len(writeStream.Data())
	off, err := this.file.Seek(0, 1)//SEEK_CUR
	if err != nil {
		return err
	}
	fmt.Println("******************88write metadata done*****************")
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
	this.onUpdateDuration(msg)
	return nil
}

func (this *SrsFlvSegment) WriteVideo(msg *rtmp.SrsRtmpMessage) error {
	this.flvEncoder.WriteVideo(uint32(msg.GetHeader().GetTimestamp()), msg.GetPayload())
	this.onUpdateDuration(msg)
	return nil
}

func (this *SrsFlvSegment) IsOverflow(maxDuration int64) bool {
	return this.duration > maxDuration
}

func (this *SrsFlvSegment) updateFlvMetaData() error {
	off, err1 := this.file.Seek(0, 2)
	if err1 != nil {
		return err1
	}
	c := utils.Float64ToBytes(float64(off), binary.BigEndian)

	_, err2 := this.file.WriteAt(c, this.filesizeOffset)
	if err2 != nil {
		return err2
	}

	b := utils.Float64ToBytes(float64(this.duration)/1000, binary.BigEndian)
	_, err3 := this.file.WriteAt(b, this.durationOffset)
	if err3 != nil {
		return err3
	}
	return nil
}