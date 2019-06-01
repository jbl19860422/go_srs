package protocol

type SrsRtmpMessage struct {
	// 4.1. Message Header
	header SrsMessageHeader
	// 4.2. Message Payload
	/**
	 * current message parsed size,
	 *       size <= header.payload_length
	 * for the payload maybe sent in multiple chunks.
	 */
	size int32
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
