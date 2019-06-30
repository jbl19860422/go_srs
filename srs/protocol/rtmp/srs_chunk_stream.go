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
package rtmp

const (
	RTMP_FMT_TYPE0 = 0
	RTMP_FMT_TYPE1 = 1
	RTMP_FMT_TYPE2 = 2
	RTMP_FMT_TYPE3 = 3
)

type SrsChunkStream struct {
	/**
	 * represents the basic header fmt,
	 * which used to identify the variant message header type.
	 */
	Format byte
	/**
	 * represents the basic header cid,
	 * which is the chunk stream id.
	 */
	Cid int32
	/**
	 * cached message header
	 */
	Header SrsMessageHeader
	/**
	 * whether the chunk message header has extended timestamp.
	 */
	ExtendedTimestamp bool

	MsgCount int32

	RtmpMessage *SrsRtmpMessage
}

func NewSrsChunkStream(cid_ int32) *SrsChunkStream {
	s := &SrsChunkStream{
		Format:            0,
		Cid:               cid_,
		ExtendedTimestamp: false,
		RtmpMessage:       nil,
		MsgCount:          0,
	}
	return s
}
