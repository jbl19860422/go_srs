/*
The MIT License (MIT)

Copyright (c) 2019 GOSRS(gosrs)

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
package flvcodec

import (
	"encoding/binary"
	"go_srs/srs/codec"
	"go_srs/srs/protocol/rtmp"
	"go_srs/srs/utils"
	"io"
)

func VideoIsKeyFrame(data []byte) bool {
	if len(data) < 1 {
		return false
	}

	frameType := (data[0] >> 4) & 0x0F
	return frameType == codec.SrsCodecVideoAVCFrameKeyFrame
}

func VideoIsSequenceHeader(data []byte) bool {
	if !VideoIsH264(data) {
		return false
	}

	if len(data) < 2 {
		return false
	}

	formatType := (data[0] >> 4) & 0x0F
	avcPacketType := data[1]
	return formatType == codec.SrsCodecVideoAVCFrameKeyFrame && avcPacketType == codec.SrsCodecVideoAVCTypeSequenceHeader
}

func AudioIsSequenceHeader(data []byte) bool {
	if !AudioIsAAC(data) {
		return false
	}

	if len(data) < 2 {
		return false
	}

	aacPacketType := data[1]
	return aacPacketType == codec.SrsCodecAudioTypeSequenceHeader
}

func VideoIsH264(data []byte) bool {
	if len(data) < 1 {
		return false
	}

	codecId := data[0] & 0x0F
	return codecId == codec.SrsCodecVideoAVC
}

func AudioIsAAC(data []byte) bool {
	if len(data) < 1 {
		return false
	}

	soundFormat := (data[0] >> 4) & 0x0F
	return soundFormat == codec.SrsCodecAudioAAC
}

func VideoIsAcceptable(data []byte) bool {
	if len(data) < 1 {
		return false
	}

	formatType := data[0]
	codecId := formatType & 0x0F
	formatType = (formatType >> 4) & 0x0F

	if formatType < 1 || formatType > 5 {
		return false
	}

	if codecId < 2 || codecId > 7 {
		return false
	}

	return true
}

const (
	AudioTagType    = 0x08
	VideoTagType    = 0x09
	MetaDataTagType = 0x18
)

const (
	SRS_FLV_TAG_HEADER_SIZE   = 11
	SRS_FLV_PREVIOUS_TAG_SIZE = 4
)

type SrsFlvHeader struct {
	signature  []byte //FLV
	version    byte
	flags      byte   //第0位和第2位,分别表示 audio 与 video 存在的情况.(1表示存在,0表示不存在)。
	headerSize []byte //即自身的总长度，一直为9, 4字节
}

func NewSrsFlvHeader(hasAudio bool, hasVideo bool) *SrsFlvHeader {
	var f byte = 0
	if hasAudio {
		f |= 1 << 0
	}

	if hasVideo {
		f |= 1 << 2
	}

	h := utils.Int32ToBytes(0x09, binary.BigEndian)
	return &SrsFlvHeader{
		signature:  []byte{'F', 'L', 'V'},
		version:    0x01,
		flags:      f,
		headerSize: h,
	}
}

func (this *SrsFlvHeader) Data() []byte {
	data := make([]byte, 0)
	data = append(data, this.signature...)
	data = append(data, this.version)
	data = append(data, this.flags)
	data = append(data, this.headerSize...)
	return data
}

type TagHeader struct {
	tagType   byte
	dataSize  []byte //3byte
	timestamp []byte //4byte
	reserved  []byte //全0
}

func NewTagHeader(typ byte, timestamp uint32, dataSize int32) *TagHeader {
	timestamp &= 0x7fffffff
	s := utils.Int32ToBytes(dataSize, binary.BigEndian)
	t := utils.UInt32ToBytes(timestamp, binary.BigEndian)[1:]
	t = append(t, byte((timestamp>>24)&0xFF))
	return &TagHeader{
		tagType:   typ,
		dataSize:  s[1:4],
		timestamp: t,
		reserved:  []byte{0, 0, 0},
	}
}

func (this *TagHeader) Data() []byte {
	data := make([]byte, 0)
	data = append(data, this.tagType)
	data = append(data, this.dataSize...)
	data = append(data, this.timestamp...)
	data = append(data, this.reserved...)
	return data
}

type SrsFlvEncoder struct {
	header   *SrsFlvHeader
	writer   io.Writer
	tagCount int32
}

func NewSrsFlvEncoder(w io.Writer) *SrsFlvEncoder {
	return &SrsFlvEncoder{
		writer:   w,
		header:   NewSrsFlvHeader(true, true),
		tagCount: 0,
	}
}

func (this *SrsFlvEncoder) WriteHeader() error {
	// 9bytes header and 4bytes first previous-tag-size
	if _, err := this.writer.Write(this.header.Data()); err != nil {
		return err
	}
	// previous tag size.
	pts := []byte{0x00, 0x00, 0x00, 0x00}
	if _, err := this.writer.Write(pts); err != nil {
		return err
	}
	return nil
}

func (this *SrsFlvEncoder) WriteMetaData(data []byte) (uint32, error) {
	header := NewTagHeader(MetaDataTagType, 0, int32(len(data)))
	return this.writeTag(header, data)
}

func (this *SrsFlvEncoder) WriteAudio(timestamp uint32, data []byte) (uint32, error) {
	header := NewTagHeader(AudioTagType, timestamp, int32(len(data)))
	return this.writeTag(header, data)
}

func (this *SrsFlvEncoder) WriteVideo(timestamp uint32, data []byte) (uint32, error) {
	header := NewTagHeader(VideoTagType, timestamp, int32(len(data)))
	return this.writeTag(header, data)
}

func (this *SrsFlvEncoder) writeTag(header *TagHeader, data []byte) (uint32, error) {
	this.tagCount++
	d := header.Data()
	d = append(d, data...)

	prevTagSize := int32(len(d))
	p := utils.Int32ToBytes(prevTagSize, binary.BigEndian)
	d = append(d, p...)
	n, err := this.writer.Write(d)
	_ = n
	return uint32(len(d)), err
}

func (this *SrsFlvEncoder) WriteTags(msgs []*rtmp.SrsRtmpMessage) error {
	for i := 0; i < len(msgs); i++ {
		if msgs[i].GetHeader().IsAudio() {
			_, _ = this.WriteAudio(uint32(msgs[i].GetHeader().GetTimestamp()), msgs[i].GetPayload())
		} else if msgs[i].GetHeader().IsVideo() {
			_, _ = this.WriteVideo(uint32(msgs[i].GetHeader().GetTimestamp()), msgs[i].GetPayload())
		} else {
			_, _ = this.WriteMetaData(msgs[i].GetPayload())
		}
	}
	return nil
}
