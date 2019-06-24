package rtmp

import (
	"go_srs/srs/protocol/skt"
	"go_srs/srs/protocol/packet"
	"go_srs/srs/protocol/amf0"
	"go_srs/srs/utils"
	"go_srs/srs/global"
	"bytes"
	_ "context"
	"encoding/binary"
	"errors"
	"log"
	"reflect"
	_ "bufio"
	"fmt"
	"time"
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
	in_chunk_size 	int32
	OutChunkSize 	int32
	OutAckSize 		AckWindowSize
	Requests 		map[float64]string
}

func NewSrsProtocol(io_ *skt.SrsIOReadWriter) *SrsProtocol {
	cache := make([]*SrsChunkStream, SRS_PERF_CHUNK_STREAM_CACHE)
	var cid int32
	for cid = 0; cid < SRS_PERF_CHUNK_STREAM_CACHE; cid++ {
		cache[cid] = NewSrsChunkStream(cid)
		cache[cid].Header.perfer_cid = cid
	}

	return &SrsProtocol{
		chunkStreams:make(map[int32]*SrsChunkStream),
		chunkCache:     cache,
		io:io_,
		in_chunk_size:  global.SRS_CONSTS_RTMP_PROTOCOL_CHUNK_SIZE,
		OutChunkSize: global.SRS_CONSTS_RTMP_PROTOCOL_CHUNK_SIZE,
	}
}

var mh_sizes = [4]int{11, 7, 3, 0}

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
	// fmt.Println("ReadMessageHeader")
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
	var is_first_chunk_of_msg bool = chunk.RtmpMessage == nil

	// but, we can ensure that when a chunk stream is fresh,
	// the fmt must be 0, a new stream.
	// fmt.Println("chunk.MsgCount=", chunk.MsgCount)
	if chunk.MsgCount == 0 && format != RTMP_FMT_TYPE0 {
		if chunk.Cid == global.RTMP_CID_ProtocolControl && format == RTMP_FMT_TYPE1 {
			// log.Print("accept cid=2, fmt=1 to make librtmp happy.")
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
	var buf1 []byte
	var mh_size = mh_sizes[format]
	if mh_size > 0 {
		if buf1, err = s.ReadNByte(mh_size); err != nil {
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
		buf_timestamp := make([]byte, 4)
		buf_timestamp[2] = buf1[pos]
		pos++
		buf_timestamp[1] = buf1[pos]
		pos++
		buf_timestamp[0] = buf1[pos]
		pos++
		buf_timestamp[3] = 0
		// fmt.Println("buf_timestamp=", buf_timestamp)
		//trans to int32
		buf_reader := bytes.NewBuffer(buf_timestamp)
		binary.Read(buf_reader, binary.LittleEndian, &chunk.Header.timestamp_delta)
		// fmt.Println("chunk.Header.timestamp_delta=", chunk.Header.timestamp_delta)
		chunk.ExtendedTimestamp = chunk.Header.timestamp_delta >= global.RTMP_EXTENDED_TIMESTAMP
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
				chunk.Header.timestamp = (int64)(chunk.Header.timestamp_delta)
			} else {
				chunk.Header.timestamp += (int64)(chunk.Header.timestamp_delta)
			}
			// fmt.Println("chunk.Header.timestamp=", chunk.Header.timestamp)
			// fmt.Println("chunk.Header.timestamp=", chunk.Header.timestamp)
		}

		if format <= RTMP_FMT_TYPE1 {
			length_buf := make([]byte, 4)
			length_buf[2] = buf1[pos]
			pos += 1
			length_buf[1] = buf1[pos]
			pos += 1
			length_buf[0] = buf1[pos]
			pos += 1
			length_buf[3] = 0
			var payload_length int32
			//trans to int32
			buf_reader := bytes.NewBuffer(length_buf)
			binary.Read(buf_reader, binary.LittleEndian, &payload_length)
			// for a message, if msg exists in cache, the size must not changed.
			// always use the actual msg size to compare, for the cache payload length can changed,
			// for the fmt type1(stream_id not changed), user can change the payload
			// length(it's not allowed in the continue chunks).
			if !is_first_chunk_of_msg && chunk.Header.payload_length != payload_length {
				err = errors.New("error rtmp packet size")
				return
			}
			// log.Printf("read payload length=%d", payload_length)
			chunk.Header.payload_length = payload_length
			chunk.Header.message_type = int8(buf1[pos])
			pos += 1

			if format == RTMP_FMT_TYPE0 {
				stream_id_buf := make([]byte, 4)
				stream_id_buf[0] = buf1[pos]
				pos += 1
				stream_id_buf[1] = buf1[pos]
				pos += 1
				stream_id_buf[2] = buf1[pos]
				pos += 1
				stream_id_buf[3] = buf1[pos]
				buf_reader := bytes.NewBuffer(length_buf)
				binary.Read(buf_reader, binary.LittleEndian, &chunk.Header.stream_id)
			}
		}
	} else {
		// update the timestamp even fmt=3 for first chunk packet
		if is_first_chunk_of_msg && !chunk.ExtendedTimestamp {
			chunk.Header.timestamp += (int64)(chunk.Header.timestamp_delta)
		}
	}

	if chunk.ExtendedTimestamp {
		mh_size += 4
		var buf2 []byte
		if buf2, err = s.ReadNByte(4); err != nil {
			return
		}

		buf_timestamp := make([]byte, 4)
		buf_timestamp[3] = buf2[0]
		buf_timestamp[2] = buf2[1]
		buf_timestamp[1] = buf2[2]
		buf_timestamp[0] = buf2[3]

		var ts int32
		buf_reader := bytes.NewBuffer(buf_timestamp)
		binary.Read(buf_reader, binary.LittleEndian, &ts)
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
		var chunk_timestamp int32 = (int32)(chunk.Header.timestamp)
		/**
		 * if chunk_timestamp<=0, the chunk previous packet has no extended-timestamp,
		 * always use the extended timestamp.
		 */
		/**
		 * about the is_first_chunk_of_msg.
		 * @remark, for the first chunk of message, always use the extended timestamp.
		 */
		if !is_first_chunk_of_msg && chunk_timestamp > 0 && chunk_timestamp != ts {
			mh_size -= 4
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
	// log.Print("start ReadBasicHeader")
	fmt, cid, err := this.ReadBasicHeader()
	if nil != err {
		log.Print("read basic header failed, err=", err)
		return nil, err
	}
	// log.Print("ReadBasicHeader done, fmt=", fmt, "&cid=", cid)
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
			chunk.Header.perfer_cid = cid
		}
	}

	err = this.ReadMessageHeader(chunk, fmt)
	if err != nil {
		// log.Print("read message header ", err)
		return nil, err
	}

	// log.Print("read message header succeed")
	var msg *SrsRtmpMessage = nil
	if msg, err = this.RecvMessagePayload(chunk); err != nil {
		// log.Print("RecvMessagePayload failed")
		return nil, err
	}
	// log.Print("read payload succeed, len=", msg.size)

	return msg, nil
}

func (s *SrsProtocol) RecvMessagePayload(chunk *SrsChunkStream) (msg *SrsRtmpMessage, err error) {
	if chunk.Header.payload_length <= 0 {
		return chunk.RtmpMessage, nil
	}

	// the chunk payload size.
	//期望的剩余数据长度=总长度-已经接收的长度
	payload_size := chunk.Header.payload_length - chunk.RtmpMessage.size

	if s.in_chunk_size < payload_size {
		payload_size = s.in_chunk_size
	}
	// fmt.Println("chunk.RtmpMessage.size=", chunk.RtmpMessage.size, "&expect payload_size=", payload_size)
	// create msg payload if not initialized
	if chunk.RtmpMessage.payload == nil {
		chunk.RtmpMessage.payload = make([]byte, 0)
	}

	// read payload to buffer
	var buffer1 []byte
	if buffer1, err = s.ReadNByte(int(payload_size)); err != nil {
		return nil, err
	}

	chunk.RtmpMessage.payload = append(chunk.RtmpMessage.payload, buffer1...)
	chunk.RtmpMessage.size += payload_size
	// log.Printf("recv payload=%x %x %x %x", chunk.RtmpMessage.payload[0], chunk.RtmpMessage.payload[1], chunk.RtmpMessage.payload[2], chunk.RtmpMessage.payload[3])

	// log.Print("recv payload_length=", chunk.Header.payload_length)

	if chunk.Header.payload_length == chunk.RtmpMessage.size {
		// log.Print("recv new message")
		new_msg := chunk.RtmpMessage
		chunk.RtmpMessage = nil
		return new_msg, nil
	}

	// log.Print("not a message payload_length=", chunk.Header.payload_length, "&size=", chunk.RtmpMessage.size)

	_ = payload_size
	_ = buffer1
	return nil, nil
}

func (s *SrsProtocol) RecvMessage() (*SrsRtmpMessage, error) {
	for {
		rtmp_msg, err := s.RecvInterlacedMessage()
		if err != nil {
			return nil, err
		}

		if rtmp_msg == nil {
			continue
		}

		if rtmp_msg.size <= 0 || rtmp_msg.header.payload_length <= 0 {
			continue
		}

		if err = s.on_recv_message(rtmp_msg); err != nil {
			return nil, err
		}

		return rtmp_msg, nil
	}
	return nil, nil
}

func (this *SrsProtocol) do_decode_message(msg *SrsRtmpMessage, stream *utils.SrsStream) (pkt packet.SrsPacket, err error) {
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
		// log.Print("srs_amf0_read_string command=", command)
		// decode command object.
		//todo other message
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
			fmt.Println("RTMP_AMF0_COMMAND_UNPUBLISH****************")
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
		p := packet.NewSrsSetChunkSizePacket()
		err = p.Decode(stream)
		// log.Print("NewSrsSetChunkSizePacket ", p.ChunkSize)
		//
		pkt = p
		return
	}
	return
}

func (s *SrsProtocol) DecodeMessage(msg *SrsRtmpMessage) (packet packet.SrsPacket, err error) {
	stream := utils.NewSrsStream(msg.payload)
	if stream == nil {
		err = errors.New("newsrsstream failed")
		return
	}
	// log.Print("NewSrsStream size=", msg.size)
	packet, err = s.do_decode_message(msg, stream)
	if err != nil {
		return
	}

	return
}

func (s *SrsProtocol) on_recv_message(msg *SrsRtmpMessage) error {
	var pkt packet.SrsPacket
	// log.Print("message.type=", msg.header.message_type)
	if msg.header.message_type == global.RTMP_MSG_SetChunkSize || msg.header.message_type == global.RTMP_MSG_UserControlMessage || msg.header.message_type == global.RTMP_MSG_WindowAcknowledgementSize {
		var err error
		pkt, err = s.DecodeMessage(msg)
		if err != nil {
			// log.Print("decode packet from message payload failed. ")
			return errors.New("decode packet from message payload failed.")
		}
		_ = pkt
	}

	if msg.header.message_type == global.RTMP_MSG_SetChunkSize {
		//参数检查
		s.in_chunk_size = pkt.(*packet.SrsSetChunkSizePacket).ChunkSize
		// log.Print("in_chunk_size=", s.in_chunk_size)
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
				// log.Print("decode message failed, err=", err1)
				continue
			}

			if reflect.TypeOf(p) != reflect.TypeOf(pkt) {
				// log.Print("drop message, ", reflect.TypeOf(p), reflect.TypeOf(pkt).Elem())
				continue
			}
			// log.Print("recv message done")
			donePkt <- p
			break
		}
	}()
	
	var tmp_pkt packet.SrsPacket
	for {
		select {
		case tmp_pkt = <-donePkt:
			// fmt.Println(reflect.TypeOf(pkt))
			reflect.ValueOf(pkt).Elem().Set(reflect.ValueOf(tmp_pkt).Elem())
			return nil
		case <- time.After(time.Second*2): 
			return errors.New("expect message timeout, type=" + reflect.TypeOf(pkt).String())
		}
	}
}

func (this *SrsProtocol) SendPacket(packet packet.SrsPacket, streamId int32) error {
	err := this.do_send_packet(packet, streamId)
	return err
}

func (this *SrsProtocol) SendMsg(msg *SrsRtmpMessage, streamId int32) error {
	return nil
}

func (this *SrsProtocol) do_send_packet(pkt packet.SrsPacket, streamId int32) error {
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
	header.payload_length = int32(len(payload))
	header.message_type = pkt.GetMessageType()
	header.stream_id = streamId
	header.perfer_cid = pkt.GetPreferCid()

	err = this.do_simple_send(&header, payload)
	if err == nil {
		return this.on_send_packet(&header, pkt)
	}
	return err
}

func (this *SrsProtocol) do_simple_send(mh *SrsMessageHeader, payload []byte) error {
	var sendedCount int = 0
	var d []byte
	var err error
	firstPkt := true
	leftPayload := payload
	for len(leftPayload) > 0 {
		if firstPkt {
			firstPkt = false
			d, err = srs_chunk_header_c0(mh.perfer_cid, int32(mh.timestamp), mh.payload_length, mh.message_type, mh.stream_id)
			if err != nil {
				return err
			}
		} else {
			d, err = srs_chunk_header_c3(mh.perfer_cid, int32(mh.timestamp))
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
		// fmt.Println("xxxxxxxxxxxxxxxxxxxxxxSendMessages ", len(msgs[i].GetPayload()), "xxxxxxxxxxxxx")
		var sendedCount int = 0
		var d []byte
		var err error
		firstPkt := true
		for len(leftPayload) > 0 {
			if firstPkt {
				firstPkt = false
				d, err = srs_chunk_header_c0(msg.GetHeader().perfer_cid, int32(msg.GetHeader().timestamp), msg.GetHeader().payload_length, msg.GetHeader().message_type, int32(streamId))
				if err != nil {
					return err
				}
			} else {
				d, err = srs_chunk_header_c3(msg.GetHeader().perfer_cid, int32(msg.GetHeader().timestamp))
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

func (this *SrsProtocol) on_send_packet(mh *SrsMessageHeader, pkt packet.SrsPacket) error {
	// fmt.Println("on_send_packet", mh.message_type)
	if pkt == nil {
		return errors.New("send pkt is nil")
	}

	switch mh.message_type {
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
		// fmt.Println("RTMP_MSG_VideoMessage")
	case global.RTMP_MSG_AudioMessage:
		// fmt.Println("RTMP_MSG_AudioMessage")
	}
	return nil
}

const SRS_CONSTS_RTMP_MAX_FMT0_HEADER_SIZE = 16

func srs_chunk_header_c0(perfer_cid int32, timestamp int32, payload_length int32, message_type int8, stream_id int32) ([]byte, error) {
	var len int32 = 0
	// to directly set the field.
	data := make([]byte, SRS_CONSTS_RTMP_MAX_FMT0_HEADER_SIZE)
	// write new chunk stream header, fmt is 0
	data[0] = byte(0x00 | (perfer_cid & 0x3F))
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
