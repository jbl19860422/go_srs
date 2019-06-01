package protocol

import (
	"bytes"
	"encoding/binary"
	"errors"
	"log"
	"net"
)

const SRS_PERF_CHUNK_STREAM_CACHE = 16

type SrsProtocol struct {
	chunkCache []*SrsChunkStream
	/**
	 * chunk stream to decode RTMP messages.
	 */
	chunkStreams map[int32]*SrsChunkStream

	/**
	 * input chunk size, default to 128, set by peer packet.
	 */
	in_chunk_size int32
}

func NewSrsProtocol() *SrsProtocol {
	cache := make([]*SrsChunkStream, SRS_PERF_CHUNK_STREAM_CACHE)
	var cid int32
	for cid = 0; cid < SRS_PERF_CHUNK_STREAM_CACHE; cid++ {
		cache[cid] = NewSrsChunkStream(cid)
	}
	return &SrsProtocol{chunkCache: cache}
}

var mh_sizes = [4]int{11, 7, 3, 0}

func (s *SrsProtocol) ReadBasicHeader(conn *net.Conn) (fmt byte, cid int32, err error) {
	var buffer1 []byte
	var buffer2 []byte
	var buffer3 []byte
	if buffer1, err = s.ReadNByte(conn, 1); err != nil {
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
		if buffer2, err = s.ReadNByte(conn, 1); err != nil {
			return
		}

		cid = 64
		cid += (int32)(buffer2[0])
	} else if cid == 1 { // 64-65599, 3B chunk header
		if buffer3, err = s.ReadNByte(conn, 2); err != nil {
			return
		}

		cid = 64
		cid += (int32)(buffer3[0])
		cid += (int32)(buffer3[1])
		return
	}
	return
}

func (s *SrsProtocol) ReadNByte(conn *net.Conn, count int) (b []byte, err error) {
	b = make([]byte, count)
	_, err = (*conn).Read(b)
	return
}

func (s *SrsProtocol) ReadMessageHeader(conn *net.Conn, chunk *SrsChunkStream, fmt byte) (err error) {
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
	if chunk.msgCount == 0 && fmt != RTMP_FMT_TYPE0 {
		if chunk.cid == RTMP_CID_ProtocolControl && fmt == RTMP_FMT_TYPE1 {
			log.Print("accept cid=2, fmt=1 to make librtmp happy.")
		} else {
			err = errors.New("error rtmp chunk start")
			return
		}
	}

	// when exists cache msg, means got an partial message,
	// the fmt must not be type0 which means new message.
	if chunk.RtmpMessage != nil && fmt == RTMP_FMT_TYPE0 {
		err = errors.New("error rtmp chunk start")
		return
	}

	if chunk.RtmpMessage == nil {
		chunk.RtmpMessage = NewSrsRtmpMessage()
	}

	// read message header from socket to buffer.
	var buf1 []byte
	var mh_size = mh_sizes[fmt]
	if mh_size > 0 {
		if buf1, err = s.ReadNByte(conn, mh_size); err != nil {
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
	if fmt <= RTMP_FMT_TYPE2 {
		buf_timestamp := make([]byte, 4)
		buf_timestamp[2] = buf1[pos]
		pos += 1
		buf_timestamp[1] = buf1[pos]
		pos += 1
		buf_timestamp[0] = buf1[pos]
		pos += 1
		buf_timestamp[3] = 0
		//trans to int32
		buf_reader := bytes.NewBuffer(buf_timestamp)
		binary.Read(buf_reader, binary.LittleEndian, &chunk.Header.timestamp_delta)
		chunk.extendedTimestamp = chunk.Header.timestamp_delta >= RTMP_EXTENDED_TIMESTAMP
		if chunk.extendedTimestamp {
			// Extended timestamp: 0 or 4 bytes
			// This field MUST be sent when the normal timsestamp is set to
			// 0xffffff, it MUST NOT be sent if the normal timestamp is set to
			// anything else. So for values less than 0xffffff the normal
			// timestamp field SHOULD be used in which case the extended timestamp
			// MUST NOT be present. For values greater than or equal to 0xffffff
			// the normal timestamp field MUST NOT be used and MUST be set to
			// 0xffffff and the extended timestamp MUST be sent.
			if fmt == RTMP_FMT_TYPE0 {
				chunk.Header.timestamp = (int64)(chunk.Header.timestamp_delta)
			} else {
				chunk.Header.timestamp += (int64)(chunk.Header.timestamp_delta)
			}
		}

		if fmt <= RTMP_FMT_TYPE1 {
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
			log.Printf("read payload length=%d", payload_length)
			chunk.Header.payload_length = payload_length
			chunk.Header.message_type = int8(buf1[pos])
			pos += 1

			if fmt == RTMP_FMT_TYPE0 {
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
		if is_first_chunk_of_msg && !chunk.extendedTimestamp {
			chunk.Header.timestamp += (int64)(chunk.Header.timestamp_delta)
		}
	}

	if chunk.extendedTimestamp {
		mh_size += 4
		var buf2 []byte
		if buf2, err = s.ReadNByte(conn, 4); err != nil {
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
	chunk.msgCount++
	return
}

func (s *SrsProtocol) RecvInterlacedMessage(conn *net.Conn) (error, *SrsRtmpMessage) {
	fmt, cid, err := s.ReadBasicHeader(conn)
	if nil != err {
		// fmt.Println("read basic header failed, err=", err)
		return nil, nil
	}

	var chunk *SrsChunkStream

	if cid < SRS_PERF_CHUNK_STREAM_CACHE {
		chunk = s.chunkCache[cid]
	} else {
		var ok bool
		if chunk, ok = s.chunkStreams[cid]; !ok {
			s.chunkStreams[cid] = NewSrsChunkStream(cid)
			chunk = s.chunkStreams[cid]
			// set the perfer cid of chunk,
			// which will copy to the message received.
			chunk.Header.perfer_cid = cid
		}
	}

	err = s.ReadMessageHeader(conn, chunk, fmt)
	if err != nil {
		log.Print("read message header ", err)
		return err, nil
	}

	var msg *SrsRtmpMessage = nil
	if msg, err = s.RecvMessagePayload(conn, chunk); err != nil {
		log.Print("RecvMessagePayload failed")
		return err, nil
	}

	return nil, msg
}

func (s *SrsProtocol) RecvMessagePayload(conn *net.Conn, chunk *SrsChunkStream) (msg *SrsRtmpMessage, err error) {
	if chunk.Header.payload_length <= 0 {
		return chunk.RtmpMessage, nil
	}

	// the chunk payload size.
	payload_size := chunk.Header.payload_length - chunk.RtmpMessage.size

	if s.in_chunk_size < payload_size {
		payload_size = s.in_chunk_size
	}

	// create msg payload if not initialized
	if chunk.RtmpMessage.payload == nil {
		chunk.RtmpMessage.payload = make([]byte, 0)
	}

	// read payload to buffer
	var buffer1 []byte
	if buffer1, err = s.ReadNByte(conn, int(chunk.Header.payload_length)); err != nil {
		return nil, err
	}

	chunk.RtmpMessage.payload = append(chunk.RtmpMessage.payload, buffer1...)
	chunk.RtmpMessage.size += payload_size
	log.Printf("recv payload=%x %x %x %x", chunk.RtmpMessage.payload[0], chunk.RtmpMessage.payload[1], chunk.RtmpMessage.payload[2], chunk.RtmpMessage.payload[3])

	log.Print("recv payload_length=", chunk.Header.payload_length)

	if chunk.Header.payload_length == chunk.RtmpMessage.size {
		new_msg := chunk.RtmpMessage
		chunk.RtmpMessage = nil
		return new_msg, nil
	}

	_ = payload_size
	_ = buffer1
	return nil, nil
}
