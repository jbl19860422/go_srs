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
package rtmp

type SrsRtmpMessage struct {
	// 4.1. Message Header
	header SrsMessageHeader
	// 4.2. Message Payload
	/**
	 * current message parsed size,
	 *       size <= header.payload_length
	 * for the payload maybe sent in multiple chunks.
	 */
	recvedSize int32
	/**
	 * the payload of message, the SrsCommonMessage never know about the detail of payload,
	 * user must use SrsProtocol.decode_message to get concrete packet.
	 * @remark, not all message payload can be decoded to packet. for example,
	 *       video/audio packet use raw bytes, no video/audio packet.
	 */
	payload []byte
}

func NewSrsRtmpMessage() *SrsRtmpMessage {
	return &SrsRtmpMessage{}
}

func (this *SrsRtmpMessage) DeepCopy() *SrsRtmpMessage {
	msg := &SrsRtmpMessage{
		header:     this.header,
		recvedSize: this.recvedSize,
	}

	msg.payload = make([]byte, len(this.payload))
	copy(msg.payload, this.payload)
	return msg
}

func (this *SrsRtmpMessage) GetHeader() *SrsMessageHeader {
	return &(this.header)
}

func (this *SrsRtmpMessage) SetHeader(h SrsMessageHeader) {
	this.header = h
}

func (this *SrsRtmpMessage) SetPayload(p []byte) {
	this.payload = make([]byte, len(p))
	copy(this.payload, p)
}

func (this *SrsRtmpMessage) GetPayload() []byte {
	return this.payload
}

func (this *SrsRtmpMessage) ChunkHeader(c0 bool) ([]byte, error) {
	if c0 {
		d, err := srs_chunk_header_c0(this.header.perferCid, int32(this.header.timestamp), this.header.payloadLength, this.header.messageType, this.header.streamId)
		return d, err
	} else {
		d, err := srs_chunk_header_c3(this.header.perferCid, int32(this.header.timestamp))
		return d, err
	}
}
