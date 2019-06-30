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

import (
	"bytes"
	_ "context"
	"encoding/binary"
	"errors"
	_ "log"
	"reflect"
	_ "bufio"
	"time"
	"go_srs/srs/protocol/skt"
	"go_srs/srs/protocol/packet"
	"go_srs/srs/protocol/amf0"
	"go_srs/srs/utils"
	"go_srs/srs/global"
)

const SRS_PERF_CHUNK_STREAM_CACHE = 16

type AckWindowSize struct {
	Window uint32
	RecvBytes int64
	SequenceNumber uint32
}

type SrsProtocol struct {
	io 				*skt.SrsIOReadWriter
	chunkCache 		[]*SrsChunkStream
	chunkStreams 	map[int32]*SrsChunkStream
	inChunkSize 	int32
	OutChunkSize 	int32
	OutAckSize 		AckWindowSize
	Requests 		map[float64]string
}

func NewSrsProtocol(io_ *skt.SrsIOReadWriter) *SrsProtocol {
	cache := make([]*SrsChunkStream, SRS_PERF_CHUNK_STREAM_CACHE)
	var cid int32
	for cid = 0; cid < SRS_PERF_CHUNK_STREAM_CACHE; cid++ {
		cache[cid] = NewSrsChunkStream(cid)
		cache[cid].Header.perferCid = cid
	}

	return &SrsProtocol{
		chunkStreams:make(map[int32]*SrsChunkStream),
		chunkCache:     cache,
		io:io_,
		inChunkSize:  global.SRS_CONSTS_RTMP_PROTOCOL_CHUNK_SIZE,
		OutChunkSize: global.SRS_CONSTS_RTMP_PROTOCOL_CHUNK_SIZE,
	}
}

var mhSizes = [4]int{11, 7, 3, 0}

func (this *SrsProtocol) ReadBasicHeader() (fmt byte, cid int32, err error) {
	var buffer1 []byte
	var buffer2 []byte
	var buffer3 []byte
	if buffer1, err = this.ReadNByte(1); err != nil {
		return
	}

	cid = (int32)(buffer1[0] & 0x3f)
	fmt = (buffer1[0] >> 6) & 0x3
	// 2-63, 1B chunk header
	if cid > 1 {
		return
	}
	// 64-319, 2B chunk header
	if cid == 0 {
		if buffer2, err = this.ReadNByte(1); err != nil {
			return
		}

		cid = 64
		cid += (int32)(buffer2[0])
	} else if cid == 1 { // 64-65599, 3B chunk header
		if buffer3, err = this.ReadNByte(2); err != nil {
			return
		}

		cid = 64
		cid += (int32)(buffer3[0])
		cid += (int32)(buffer3[1])
		return
	}
	return
}

func (this *SrsProtocol) ReadNByte(count int) (b []byte, err error) {
	b = make([]byte, count)
	_, err = this.io.ReadFully(b, 1000)
	return
}

func (s *SrsProtocol) ReadMessageHeader(chunk *SrsChunkStream, format byte) (err error) {
	/**
	 * we should not assert anything about fmt, for the first packet.
	 * (when first packet, the chunk->msg is NULL).
	 * the fmt maybe 0/1/2/3, the FMLE will send a 0xC4 for some audio packet.
	 * the previous packet is:
	 *     04                // fmt=0, cid=4
	 *     00 00 1a          // timestamp=26
	 *     00 00 9d          // payload_length=157
	 *     08                // message_type=8(audio)
	 *     01 00 00 00       // stream_id=1
	 * the current packet maybe:
	 *     c4             // fmt=3, cid=4
	 * it's ok, for the packet is audio, and timestamp delta is 26.
	 * the current packet must be parsed as:
	 *     fmt=0, cid=4
	 *     timestamp=26+26=52
	 *     payload_length=157
	 *     message_type=8(audio)
	 *     stream_id=1
	 * so we must update the timestamp even fmt=3 for first packet.
	 */
	// fresh packet used to update the timestamp even fmt=3 for first packet.
	// fresh packet always means the chunk is the first one of message.
	var isFirstChunkOfMsg bool = chunk.RtmpMessage == nil

	// but, we can ensure that when a chunk stream is fresh,
	// the fmt must be 0, a new stream.
	if chunk.MsgCount == 0 && format != RTMP_FMT_TYPE0 {
		if chunk.Cid == global.RTMP_CID_ProtocolControl && format == RTMP_FMT_TYPE1 {
		} else {
			err = errors.New("error rtmp chunk start")
			return
		}
	}

	// when exists cache msg, means got an partial message,
	// the fmt must not be type0 which means new message.
	if chunk.RtmpMessage != nil && format == RTMP_FMT_TYPE0 {
		err = errors.New("error rtmp chunk start")
		return
	}

	if chunk.RtmpMessage == nil {
		chunk.RtmpMessage = NewSrsRtmpMessage()
	}

	// read message header from socket to buffer.
	var buf []byte
	var mhSize = mhSizes[format]
	if mhSize > 0 {
		if buf, err = s.ReadNByte(mhSize); err != nil {
			return
		}
	}

	/**
	 * parse the message header.
	 *   3bytes: timestamp delta,    fmt=0,1,2
	 *   3bytes: payload length,     fmt=0,1
	 *   1bytes: message type,       fmt=0,1
	 *   4bytes: stream id,          fmt=0
	 * where:
	 *   fmt=0, 0x0X
	 *   fmt=1, 0x4X
	 *   fmt=2, 0x8X
	 *   fmt=3, 0xCX
	 */
	var pos int32 = 0
	if format <= RTMP_FMT_TYPE2 {
		bufTimestamp := make([]byte, 4)
		bufTimestamp[2] = buf[pos]
		pos++
		bufTimestamp[1] = buf[pos]
		pos++
		bufTimestamp[0] = buf[pos]
		pos++
		bufTimestamp[3] = 0
		//trans to int32
		bufReader := bytes.NewBuffer(bufTimestamp)
		binary.Read(bufReader, binary.LittleEndian, &chunk.Header.timestampDelta)
		chunk.ExtendedTimestamp = chunk.Header.timestampDelta >= global.RTMP_EXTENDED_TIMESTAMP
		if !chunk.ExtendedTimestamp {
			// Extended timestamp: 0 or 4 bytes
			// This field MUST be sent when the normal timsestamp is set to
			// 0xffffff, it MUST NOT be sent if the normal timestamp is set to
			// anything else. So for values less than 0xffffff the normal
			// timestamp field SHOULD be used in which case the extended timestamp
			// MUST NOT be present. For values greater than or equal to 0xffffff
			// the normal timestamp field MUST NOT be used and MUST be set to
			// 0xffffff and the extended timestamp MUST be sent.
			
			if format == RTMP_FMT_TYPE0 {
				chunk.Header.timestamp = (int64)(chunk.Header.timestampDelta)
			} else {
				chunk.Header.timestamp += (int64)(chunk.Header.timestampDelta)
			}
		}

		if format <= RTMP_FMT_TYPE1 {
			lengthBuf := make([]byte, 4)
			lengthBuf[2] = buf[pos]
			pos += 1
			lengthBuf[1] = buf[pos]
			pos += 1
			lengthBuf[0] = buf[pos]
			pos += 1
			lengthBuf[3] = 0
			var payloadLength int32
			//trans to int32
			bufReader := bytes.NewBuffer(lengthBuf)
			binary.Read(bufReader, binary.LittleEndian, &payloadLength)
			// for a message, if msg exists in cache, the size must not changed.
			// always use the actual msg size to compare, for the cache payload length can changed,
			// for the fmt type1(stream_id not changed), user can change the payload
			// length(it's not allowed in the continue chunks).
			if !isFirstChunkOfMsg && chunk.Header.payloadLength != payloadLength {
				err = errors.New("error rtmp packet size")
				return
			}

			chunk.Header.payloadLength = payloadLength
			chunk.Header.messageType = int8(buf[pos])
			pos += 1

			if format == RTMP_FMT_TYPE0 {
				streamIdBuf := make([]byte, 4)
				streamIdBuf[0] = buf[pos]
				pos += 1
				streamIdBuf[1] = buf[pos]
				pos += 1
				streamIdBuf[2] = buf[pos]
				pos += 1
				streamIdBuf[3] = buf[pos]
				bufReader := bytes.NewBuffer(lengthBuf)
				binary.Read(bufReader, binary.LittleEndian, &chunk.Header.streamId)
			}
		}
	} else {
		// update the timestamp even fmt=3 for first chunk packet
		if isFirstChunkOfMsg && !chunk.ExtendedTimestamp {
			chunk.Header.timestamp += (int64)(chunk.Header.timestampDelta)
		}
	}

	if chunk.ExtendedTimestamp {
		mhSize += 4
		var buf2 []byte
		if buf2, err = s.ReadNByte(4); err != nil {
			return
		}

		bufTimestamp := make([]byte, 4)
		bufTimestamp[3] = buf2[0]
		bufTimestamp[2] = buf2[1]
		bufTimestamp[1] = buf2[2]
		bufTimestamp[0] = buf2[3]

		var ts int32
		bufReader := bytes.NewBuffer(bufTimestamp)
		binary.Read(bufReader, binary.LittleEndian, &ts)
		// always use 31bits timestamp, for some server may use 32bits extended timestamp.
		// @see https://github.com/ossrs/srs/issues/111
		ts &= 0x7ffffff

		/**
		 * RTMP specification and ffmpeg/librtmp is false,
		 * but, adobe changed the specification, so flash/FMLE/FMS always true.
		 * default to true to support flash/FMLE/FMS.
		 *
		 * ffmpeg/librtmp may donot send this filed, need to detect the value.
		 * @see also: http://blog.csdn.net/win_lin/article/details/13363699
		 * compare to the chunk timestamp, which is set by chunk message header
		 * type 0,1 or 2.
		 *
		 * @remark, nginx send the extended-timestamp in sequence-header,
		 * and timestamp delta in continue C1 chunks, and so compatible with ffmpeg,
		 * that is, there is no continue chunks and extended-timestamp in nginx-rtmp.
		 *
		 * @remark, srs always send the extended-timestamp, to keep simple,
		 * and compatible with adobe products.
		 */
		var chunkTimestamp int32 = (int32)(chunk.Header.timestamp)
		/**
		 * if chunk_timestamp<=0, the chunk previous packet has no extended-timestamp,
		 * always use the extended timestamp.
		 */
		/**
		 * about the is_first_chunk_of_msg.
		 * @remark, for the first chunk of message, always use the extended timestamp.
		 */
		if !isFirstChunkOfMsg && chunkTimestamp > 0 && chunkTimestamp != ts {
			mhSize -= 4
			//这里需要考虑下怎么处理
			//no 4bytes extended timestamp in the continued chunk
		} else {
			chunk.Header.timestamp = (int64)(ts)
		}
	}

	// the extended-timestamp must be unsigned-int,
	//         24bits timestamp: 0xffffff = 16777215ms = 16777.215s = 4.66h
	//         32bits timestamp: 0xffffffff = 4294967295ms = 4294967.295s = 1193.046h = 49.71d
	// because the rtmp protocol says the 32bits timestamp is about "50 days":
	//         3. Byte Order, Alignment, and Time Format
	//                Because timestamps are generally only 32 bits long, they will roll
	//                over after fewer than 50 days.
	//
	// but, its sample says the timestamp is 31bits:
	//         An application could assume, for example, that all
	//        adjacent timestamps are within 2^31 milliseconds of each other, so
	//        10000 comes after 4000000000, while 3000000000 comes before
	//        4000000000.
	// and flv specification says timestamp is 31bits:
	//        Extension of the Timestamp field to form a SI32 value. This
	//        field represents the upper 8 bits, while the previous
	//        Timestamp field represents the lower 24 bits of the time in
	//        milliseconds.
	// in a word, 31bits timestamp is ok.
	// convert extended timestamp to 31bits.

	chunk.Header.timestamp &= 0x7fffffff
	// copy header to msg
	chunk.RtmpMessage.header = chunk.Header

	// increase the msg count, the chunk stream can accept fmt=1/2/3 message now.
	chunk.MsgCount++
	return
}

func (this *SrsProtocol) RecvInterlacedMessage() (*SrsRtmpMessage, error) {
	fmt, cid, err := this.ReadBasicHeader()
	if nil != err {
		return nil, err
	}
	var chunk *SrsChunkStream

	if cid < SRS_PERF_CHUNK_STREAM_CACHE {
		chunk = this.chunkCache[cid]
	} else {
		var ok bool
		if chunk, ok = this.chunkStreams[cid]; !ok {
			this.chunkStreams[cid] = NewSrsChunkStream(cid)
			chunk = this.chunkStreams[cid]
			// set the perfer cid of chunk,
			// which will copy to the message received.
			chunk.Header.perferCid = cid
		}
	}

	err = this.ReadMessageHeader(chunk, fmt)
	if err != nil {
		return nil, err
	}

	var msg *SrsRtmpMessage = nil
	if msg, err = this.RecvMessagePayload(chunk); err != nil {
		return nil, err
	}

	return msg, nil
}

func (s *SrsProtocol) RecvMessagePayload(chunk *SrsChunkStream) (msg *SrsRtmpMessage, err error) {
	if chunk.Header.payloadLength <= 0 {
		return chunk.RtmpMessage, nil
	}

	// the chunk payload size.
	//期望的剩余数据长度=总长度-已经接收的长度
	payloadSize := chunk.Header.payloadLength - chunk.RtmpMessage.recvedSize

	if s.inChunkSize < payloadSize {//如果长度大于in_chunk_size，则最大是in_chunk_size
		payloadSize = s.inChunkSize
	}
	// create msg payload if not initialized
	if chunk.RtmpMessage.payload == nil {
		chunk.RtmpMessage.payload = make([]byte, 0)
	}

	// read payload to buffer
	var buf []byte
	if buf, err = s.ReadNByte(int(payloadSize)); err != nil {
		return nil, err
	}

	chunk.RtmpMessage.payload = append(chunk.RtmpMessage.payload, buf...)
	chunk.RtmpMessage.recvedSize += payloadSize
	
	if chunk.Header.payloadLength == chunk.RtmpMessage.recvedSize {
		newMsg := chunk.RtmpMessage
		chunk.RtmpMessage = nil
		return newMsg, nil
	}
	return nil, nil
}

func (s *SrsProtocol) RecvMessage() (*SrsRtmpMessage, error) {
	for {
		rtmpMsg, err := s.RecvInterlacedMessage()
		if err != nil {
			return nil, err
		}

		if rtmpMsg == nil {
			continue
		}

		if rtmpMsg.recvedSize <= 0 || rtmpMsg.header.payloadLength <= 0 {
			continue
		}

		if err = s.OnRecvRtmpMessage(rtmpMsg); err != nil {
			return nil, err
		}

		return rtmpMsg, nil
	}
	return nil, nil
}

func (this *SrsProtocol) doDecodeMessage(msg *SrsRtmpMessage, stream *utils.SrsStream) (pkt packet.SrsPacket, err error) {
	if msg.header.IsAmf0Command() || msg.header.IsAmf3Command() || msg.header.IsAmf0Data() || msg.header.IsAmf3Data() {
		// skip 1bytes to decode the amf3 command.
		if msg.header.IsAmf3Command() && stream.Require(1) {
			stream.Skip(1)
		}
		// amf0 command message.
		// need to read the command name.
		var amf0Command amf0.SrsAmf0String
		err = amf0Command.Decode(stream)
		if err != nil {
			err = errors.New("srs_amf0_read_string error")
			return
		}
		command := amf0Command.Value.Value
		// decode command object.
		// todo other message
		if command == amf0.RTMP_AMF0_COMMAND_CONNECT {
			pkt = packet.NewSrsConnectAppPacket()
			err = pkt.Decode(stream)
			return
		} else if command == amf0.RTMP_AMF0_COMMAND_PLAY {
			pkt = packet.NewSrsPlayPacket()
			err = pkt.Decode(stream)
			return
		} else if command == amf0.RTMP_AMF0_COMMAND_RELEASE_STREAM {
			pkt = packet.NewSrsFMLEStartPacket(command)
			err = pkt.Decode(stream)
			return
		} else if command == amf0.RTMP_AMF0_COMMAND_FC_PUBLISH {
			pkt = packet.NewSrsFMLEStartPacket(command)
			err = pkt.Decode(stream)
			return
		} else if command == amf0.RTMP_AMF0_COMMAND_CREATE_STREAM {
			pkt = packet.NewSrsCreateStreamPacket()
			err = pkt.Decode(stream)
			return
		} else if command == amf0.RTMP_AMF0_COMMAND_PUBLISH {
			pkt = packet.NewSrsPublishPacket()
			err = pkt.Decode(stream)
			return
		}  else if command == amf0.RTMP_AMF0_COMMAND_UNPUBLISH {
            pkt = packet.NewSrsFMLEStartPacket(command)
			err = pkt.Decode(stream)
        } else if command == amf0.RTMP_AMF0_COMMAND_CLOSE_STREAM {
			pkt = packet.NewSrsCloseStreamPacket()
			err = pkt.Decode(stream)
			return
        } else if command == amf0.SRS_CONSTS_RTMP_SET_DATAFRAME || command == amf0.SRS_CONSTS_RTMP_ON_METADATA {
			pkt = packet.NewSrsOnMetaDataPacket(command)
			err = pkt.Decode(stream)
			return 
        } 
	} else if msg.header.IsSetChunkSize() {
		pkt = packet.NewSrsSetChunkSizePacket()
		err = pkt.Decode(stream)
		return
	}
	return
}

func (s *SrsProtocol) DecodeMessage(msg *SrsRtmpMessage) (packet packet.SrsPacket, err error) {
	stream := utils.NewSrsStream(msg.payload)
	if stream == nil {
		err = errors.New("NewSrsStream failed")
		return
	}
	packet, err = s.doDecodeMessage(msg, stream)
	if err != nil {
		return
	}

	return
}

func (s *SrsProtocol) OnRecvRtmpMessage(msg *SrsRtmpMessage) error {
	var pkt packet.SrsPacket
	if msg.header.messageType == global.RTMP_MSG_SetChunkSize || msg.header.messageType == global.RTMP_MSG_UserControlMessage || msg.header.messageType == global.RTMP_MSG_WindowAcknowledgementSize {
		var err error
		pkt, err = s.DecodeMessage(msg)
		if err != nil {
			return errors.New("decode packet from message payload failed.")
		}
		_ = pkt
	}

	if msg.header.messageType == global.RTMP_MSG_SetChunkSize {
		//参数检查
		s.inChunkSize = pkt.(*packet.SrsSetChunkSizePacket).ChunkSize
	}

	return nil
}

func (this *SrsProtocol) ExpectMessage(pkt packet.SrsPacket) error {
	if reflect.TypeOf(pkt).Kind() != reflect.Ptr {
		return errors.New("need ptr to store result")
	}

	donePkt := make(chan packet.SrsPacket)
	//todo 这里需要修改为cancelctx
	go func() {
		for {
			msg, err := this.RecvMessage()
			if err != nil {
				continue
			}
	
			if msg == nil {
				continue
			}
	
			p, err1 := this.DecodeMessage(msg)
			if err1 != nil {
				continue
			}

			if reflect.TypeOf(p) != reflect.TypeOf(pkt) {
				continue
			}
			donePkt <- p
			break
		}
	}()
	
	var tmp_pkt packet.SrsPacket
	for {
		select {
		case tmp_pkt = <-donePkt:
			reflect.ValueOf(pkt).Elem().Set(reflect.ValueOf(tmp_pkt).Elem())
			return nil
		case <- time.After(time.Second*2): 
			return errors.New("expect message timeout, type=" + reflect.TypeOf(pkt).String())
		}
	}
}

func (this *SrsProtocol) SendPacket(packet packet.SrsPacket, streamId int32) error {
	err := this.doSendPacket(packet, streamId)
	return err
}

func (this *SrsProtocol) SendMsg(msg *SrsRtmpMessage, streamId int32) error {
	return nil
}

func (this *SrsProtocol) doSendPacket(pkt packet.SrsPacket, streamId int32) error {
	stream := utils.NewSrsStream([]byte{})
	err := pkt.Encode(stream)
	if err != nil {
		return err
	}

	payload := stream.Data()
	if len(payload) <= 0 {
		return errors.New("packet is empty, ignore empty message.")
	}
	var header SrsMessageHeader
	header.payloadLength = int32(len(payload))
	header.messageType = pkt.GetMessageType()
	header.streamId = streamId
	header.perferCid = pkt.GetPreferCid()

	err = this.doSimpleSend(&header, payload)
	if err == nil {
		return this.onSendPacket(&header, pkt)
	}
	return err
}

func (this *SrsProtocol) doSimpleSend(mh *SrsMessageHeader, payload []byte) error {
	var sendedCount int = 0
	var d []byte
	var err error
	firstPkt := true
	leftPayload := payload
	for len(leftPayload) > 0 {
		if firstPkt {
			firstPkt = false
			d, err = srs_chunk_header_c0(mh.perferCid, int32(mh.timestamp), mh.payloadLength, mh.messageType, mh.streamId)
			if err != nil {
				return err
			}
		} else {
			d, err = srs_chunk_header_c3(mh.perferCid, int32(mh.timestamp))
		}

		payloadSize := utils.MinInt32(int32(len(leftPayload)), this.OutChunkSize)//int32(len(leftPayload))//
		sendPayload := leftPayload[:payloadSize]
		leftPayload = leftPayload[payloadSize:]
		d = append(d, sendPayload...)
		n2, err2 := this.io.Write(d)
		if err2 != nil {
			return err2
		}
		sendedCount += n2
	}
	return nil
}

func (this *SrsProtocol) SendMessages(msgs []*SrsRtmpMessage, streamId int) error {
	for i := 0; i < len(msgs); i++ {
		if msgs[i] == nil {
			continue
		}

		if len(msgs[i].GetPayload()) <= 0 {
			continue
		}

		msg := msgs[i]
		leftPayload := msg.GetPayload()
		var sendedCount int = 0
		var d []byte
		var err error
		firstPkt := true
		for len(leftPayload) > 0 {
			if firstPkt {
				firstPkt = false
				d, err = srs_chunk_header_c0(msg.GetHeader().perferCid, int32(msg.GetHeader().timestamp), msg.GetHeader().payloadLength, msg.GetHeader().messageType, int32(streamId))
				if err != nil {
					return err
				}
			} else {
				d, err = srs_chunk_header_c3(msg.GetHeader().perferCid, int32(msg.GetHeader().timestamp))
			}

			payloadSize := utils.MinInt32(int32(len(leftPayload)), this.OutChunkSize)//int32(len(leftPayload))//
			sendPayload := leftPayload[:payloadSize]
			leftPayload = leftPayload[payloadSize:]
			d = append(d, sendPayload...)
			n2, err2 := this.io.Write(d)
			if err2 != nil {
				return err2
			}
			sendedCount += n2
		}
		return nil

	}
	return nil
}

func (this *SrsProtocol) onSendPacket(mh *SrsMessageHeader, pkt packet.SrsPacket) error {
	if pkt == nil {
		return errors.New("send pkt is nil")
	}

	switch mh.messageType {
	case global.RTMP_MSG_SetChunkSize:
		this.OutChunkSize = pkt.(*packet.SrsSetChunkSizePacket).ChunkSize
	case global.RTMP_MSG_WindowAcknowledgementSize:
		this.OutAckSize.Window = uint32(pkt.(*packet.SrsSetWindowAckSizePacket).AckowledgementWindowSize)
	case global.RTMP_MSG_AMF0CommandMessage, global.RTMP_MSG_AMF3CommandMessage:
		switch pkt.(type) {
			case *packet.SrsConnectAppPacket:{
				p := pkt.(*packet.SrsConnectAppPacket)
				this.Requests[p.TransactionId.GetValue().(float64)] = p.CommandName.GetValue().(string)
			}
			case *packet.SrsCreateStreamPacket:{
				p := pkt.(*packet.SrsCreateStreamPacket)
				this.Requests[p.TransactionId.GetValue().(float64)] = p.CommandName.GetValue().(string)
			}
			case *packet.SrsFMLEStartPacket:{
				p := pkt.(*packet.SrsFMLEStartPacket)
				this.Requests[p.TransactionId.GetValue().(float64)] = p.CommandName.GetValue().(string)
			}
		}
	case global.RTMP_MSG_VideoMessage:
		//todo
	case global.RTMP_MSG_AudioMessage:
		//todo
	}
	return nil
}

const SRS_CONSTS_RTMP_MAX_FMT0_HEADER_SIZE = 16

func srs_chunk_header_c0(perferCid int32, timestamp int32, payload_length int32, message_type int8, stream_id int32) ([]byte, error) {
	var len int32 = 0
	// to directly set the field.
	data := make([]byte, SRS_CONSTS_RTMP_MAX_FMT0_HEADER_SIZE)
	// write new chunk stream header, fmt is 0
	data[0] = byte(0x00 | (perferCid & 0x3F))
	len += 1
	// chunk message header, 11 bytes
	// timestamp, 3bytes, big-endian
	if timestamp < global.RTMP_EXTENDED_TIMESTAMP {
		b := utils.Int32ToBytes(timestamp, binary.LittleEndian)
		data[1] = b[2]
		data[2] = b[1]
		data[3] = b[0]
	} else {//有扩展字段，则timestamp全f
		data[1] = 0xFF
		data[2] = 0xFF
		data[3] = 0xFF
	}
	len += 3

	// message_length, 3bytes, big-endian
	b := utils.Int32ToBytes(payload_length, binary.LittleEndian)
	data[4] = b[2]
	data[5] = b[1]
	data[6] = b[0]
	len += 3
	// message_type, 1bytes
	data[7] = byte(message_type)
	// log.Print("data[7]=", data[7])
	len += 1
	// stream_id, 4bytes, little-endian
	b = utils.Int32ToBytes(stream_id, binary.LittleEndian)
	data[8] = b[0]
	data[9] = b[1]
	data[10] = b[2]
	data[11] = b[3]
	len += 4
	// for c0
	// chunk extended timestamp header, 0 or 4 bytes, big-endian
	//
	// for c3:
	// chunk extended timestamp header, 0 or 4 bytes, big-endian
	// 6.1.3. Extended Timestamp
	// This field is transmitted only when the normal time stamp in the
	// chunk message header is set to 0x00ffffff. If normal time stamp is
	// set to any value less than 0x00ffffff, this field MUST NOT be
	// present. This field MUST NOT be present if the timestamp field is not
	// present. Type 3 chunks MUST NOT have this field.
	// adobe changed for Type3 chunk:
	//        FMLE always sendout the extended-timestamp,
	//        must send the extended-timestamp to FMS,
	//        must send the extended-timestamp to flash-player.
	// @see: ngx_rtmp_prepare_message
	// @see: http://blog.csdn.net/win_lin/article/details/13363699
	// TODO: FIXME: extract to outer.
	if timestamp >= global.RTMP_EXTENDED_TIMESTAMP {
		b = utils.Int32ToBytes(timestamp, binary.BigEndian)
		data[12] = b[3]
		data[13] = b[2]
		data[14] = b[1]
		data[15] = b[0]
		len += 4
	}
	return data[:len], nil
}

const SRS_CONSTS_RTMP_MAX_FMT3_HEADER_SIZE = 5
func srs_chunk_header_c3(prefer_cid int32, timestamp int32) ([]byte, error) {
	// to directly set the field.
	var len int32 = 0
	// to directly set the field.
	data := make([]byte, SRS_CONSTS_RTMP_MAX_FMT3_HEADER_SIZE)
    // write no message header chunk stream, fmt is 3
    // @remark, if perfer_cid > 0x3F, that is, use 2B/3B chunk header,
    // SRS will rollback to 1B chunk header.
	data[0] = byte(0xC0 | (prefer_cid & 0x3F))
    len++
    // for c0
    // chunk extended timestamp header, 0 or 4 bytes, big-endian
    //
    // for c3:
    // chunk extended timestamp header, 0 or 4 bytes, big-endian
    // 6.1.3. Extended Timestamp
    // This field is transmitted only when the normal time stamp in the
    // chunk message header is set to 0x00ffffff. If normal time stamp is
    // set to any value less than 0x00ffffff, this field MUST NOT be
    // present. This field MUST NOT be present if the timestamp field is not
    // present. Type 3 chunks MUST NOT have this field.
    // adobe changed for Type3 chunk:
    //        FMLE always sendout the extended-timestamp,
    //        must send the extended-timestamp to FMS,
    //        must send the extended-timestamp to flash-player.
    // @see: ngx_rtmp_prepare_message
    // @see: http://blog.csdn.net/win_lin/article/details/13363699
    // TODO: FIXME: extract to outer.
    if (timestamp >= global.RTMP_EXTENDED_TIMESTAMP) {
        b := utils.Int32ToBytes(timestamp, binary.BigEndian)
		data[1] = b[3]
		data[2] = b[2]
		data[3] = b[1]
		data[4] = b[0]
		len += 4
	}
	return data[:len], nil
}
