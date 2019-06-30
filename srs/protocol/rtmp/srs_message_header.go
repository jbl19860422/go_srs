package rtmp

import (
	"go_srs/srs/global"
)

//message header
type SrsMessageHeader struct {
	/**
	 * 3bytes.
	 * Three-byte field that contains a timestamp delta of the message.
	 * @remark, only used for decoding message from chunk stream.
	 */
	timestampDelta int32
	/**
	 * 3bytes.
	 * Three-byte field that represents the size of the payload in bytes.
	 * It is set in big-endian format.
	 */
	payloadLength int32
	/**
	 * 1byte.
	 * One byte field to represent the message type. A range of type IDs
	 * (1-7) are reserved for protocol control messages.
	 */
	message_type int8

	/**
	* 4bytes.
	* Four-byte field that identifies the stream of the message. These
	* bytes are set in little-endian format.
	 */
	stream_id int32
	/**
	* Four-byte field that contains a timestamp of the message.
	* The 4 bytes are packed in the big-endian order.
	* @remark, used as calc timestamp when decode and encode time.
	* @remark, we use 64bits for large time for jitter detect and hls.
	 */
	timestamp int64
	/**
	* get the perfered cid(chunk stream id) which sendout over.
	* set at decoding, and canbe used for directly send message,
	* for example, dispatch to all connections.
	 */
	perfer_cid int32
}

func (this *SrsMessageHeader) GetTimestamp() int64 {
	return this.timestamp
}

func (this *SrsMessageHeader) SetTimestamp(t int64) {
	this.timestamp = t
}

func (this *SrsMessageHeader) Print() {
}

func (s *SrsMessageHeader) IsAudio() bool {
	return s.message_type == global.RTMP_MSG_AudioMessage
}

func (s *SrsMessageHeader) IsVideo() bool {
	return s.message_type == global.RTMP_MSG_VideoMessage
}

func (s *SrsMessageHeader) IsAmf0Command() bool {
	return s.message_type == global.RTMP_MSG_AMF0CommandMessage
}

func (s *SrsMessageHeader) IsAmf0Data() bool {
	return s.message_type == global.RTMP_MSG_AMF0DataMessage
}

func (s *SrsMessageHeader) IsAmf3Command() bool {
	return s.message_type == global.RTMP_MSG_AMF3CommandMessage
}

func (s *SrsMessageHeader) IsAmf3Data() bool {
	return s.message_type == global.RTMP_MSG_AMF3DataMessage
}

func (s *SrsMessageHeader) IsWindowAckledgementSize() bool {
	return s.message_type == global.RTMP_MSG_WindowAcknowledgementSize
}

func (s *SrsMessageHeader) IsAckledgement() bool {
	return s.message_type == global.RTMP_MSG_Acknowledgement
}

func (s *SrsMessageHeader) IsSetChunkSize() bool {
	return s.message_type == global.RTMP_MSG_SetChunkSize
}

func (s *SrsMessageHeader) IsUserControlMessage() bool {
	return s.message_type == global.RTMP_MSG_UserControlMessage
}

func (s *SrsMessageHeader) IsSetPeerBandwidth() bool {
	return s.message_type == global.RTMP_MSG_SetPeerBandwidth
}

func (s *SrsMessageHeader) IsAggregate() bool {
	return s.message_type == global.RTMP_MSG_AggregateMessage
}

func (this *SrsMessageHeader) SetLength(len int32) {
	this.payloadLength = len
}
