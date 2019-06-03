package protocol

//message header
type SrsMessageHeader {
	/**
     * 3bytes.
     * Three-byte field that contains a timestamp delta of the message.
     * @remark, only used for decoding message from chunk stream.
     */
	timestamp_delta int32
	/**
     * 3bytes.
     * Three-byte field that represents the size of the payload in bytes.
     * It is set in big-endian format.
     */
	payload_length int32
	 /**
     * 1byte.
     * One byte field to represent the message type. A range of type IDs
     * (1-7) are reserved for protocol control messages.
	 */
	message_byte int8

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

func (s *SrsMessageHeader) IsAudio() bool {

}

func (s *SrsMessageHeader) IsVideo() bool {

}

func (s *SrsMessageHeader) IsAmf0Command() bool {

}

func (s *SrsMessageHeader) IsAmf0Data() bool {

}

func (s *SrsMessageHeader) IsAmf3Command() bool {

}

func (s *SrsMessageHeader) IsAmf3Data() bool {

}

func (s *SrsMessageHeader) IsWindowAckledgementSize() bool {

}

func (s *SrsMessageHeader) IsAckledgement() bool {

}

func (s *SrsMessageHeader) IsSetChunkSize() bool {

}

func (s *SrsMessageHeader) IsUserControlMessage() bool {

}

func (s *SrsMessageHeader) IsSetPeerBandwidth() bool {

}

func (s *SrsMessageHeader) IsAggregate() bool {
	
}