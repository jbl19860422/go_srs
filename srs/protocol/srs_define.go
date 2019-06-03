package protocol

const (
	RTMP_MSG_SetChunkSize               = 0x01
	RTMP_MSG_AbortMessage               = 0x02
	RTMP_MSG_Acknowledgement            = 0x03
	RTMP_MSG_UserControlMessage         = 0x04
	RTMP_MSG_WindowAcknowledgementSize  = 0x05
	RTMP_MSG_SetPeerBandwidth           = 0x06
	RTMP_MSG_EdgeAndOriginServerCommand = 0x07
	/**
	3. Types of messages
	The server and the client send messages over the network to
	communicate with each other. The messages can be of any type which
	includes audio messages, video messages, command messages, shared
	object messages, data messages, and user control messages.
	3.1. Command message
	Command messages carry the AMF-encoded commands between the client
	and the server. These messages have been assigned message type value
	of 20 for AMF0 encoding and message type value of 17 for AMF3
	encoding. These messages are sent to perform some operations like
	connect, createStream, publish, play, pause on the peer. Command
	messages like onstatus, result etc. are used to inform the sender
	about the status of the requested commands. A command message
	consists of command name, transaction ID, and command object that
	contains related parameters. A client or a server can request Remote
	Procedure Calls (RPC) over streams that are communicated using the
	command messages to the peer.
	*/
	RTMP_MSG_AMF3CommandMessage = 0x11 //17	AMF3
	RTMP_MSG_AMF0CommandMessage = 0x14 //20	AFM0
	/**
	3.2. Data message
	The client or the server sends this message to send Metadata or any
	user data to the peer. Metadata includes details about the
	data(audio, video etc.) like creation time, duration, theme and so
	on. These messages have been assigned message type value of 18 for
	AMF0 and message type value of 15 for AMF3.
	*/
	RTMP_MSG_AMF0DataMessage = 0x12
	RTMP_MSG_AMF3DataMessage = 0x0f
	/**
	3.3. Shared object message
	A shared object is a Flash object (a collection of name value pairs)
	that are in synchronization across multiple clients, instances, and
	so on. The message types kMsgContainer=19 for AMF0 and
	kMsgContainerEx=16 for AMF3 are reserved for shared object events.
	Each message can contain multiple events.
	*/
	RTMP_MSG_AMF3SharedObject = 0x10
	RTMP_MSG_AMF0SharedObject = 0x13
	/**
	3.4. Audio message
	The client or the server sends this message to send audio data to the
	peer. The message type value of 8 is reserved for audio messages.
	*/
	RTMP_MSG_AudioMessage = 0x08
	/* *
	3.5. Video message
	The client or the server sends this message to send video data to the
	peer. The message type value of 9 is reserved for video messages.
	These messages are large and can delay the sending of other type of
	messages. To avoid such a situation, the video message is assigned
	the lowest priority.
	*/
	RTMP_MSG_VideoMessage = 0x09
	/**
	3.6. Aggregate message
	An aggregate message is a single message that contains a list of submessages.
	The message type value of 22 is reserved for aggregate
	messages.
	*/
	RTMP_MSG_AggregateMessage = 0x16
)

/**
 * amf0 command message, command name macros
 */
const (
	RTMP_AMF0_COMMAND_CONNECT        = "connect"
	RTMP_AMF0_COMMAND_CREATE_STREAM  = "createStream"
	RTMP_AMF0_COMMAND_CLOSE_STREAM   = "closeStream"
	RTMP_AMF0_COMMAND_PLAY           = "play"
	RTMP_AMF0_COMMAND_PAUSE          = "pause"
	RTMP_AMF0_COMMAND_ON_BW_DONE     = "onBWDone"
	RTMP_AMF0_COMMAND_ON_STATUS      = "onStatus"
	RTMP_AMF0_COMMAND_RESULT         = "_result"
	RTMP_AMF0_COMMAND_ERROR          = "_error"
	RTMP_AMF0_COMMAND_RELEASE_STREAM = "releaseStream"
	RTMP_AMF0_COMMAND_FC_PUBLISH     = "FCPublish"
	RTMP_AMF0_COMMAND_UNPUBLISH      = "FCUnpublish"
	RTMP_AMF0_COMMAND_PUBLISH        = "publish"
	RTMP_AMF0_DATA_SAMPLE_ACCESS     = "|RtmpSampleAccess"
)

/****************************************************************************
 *****************************************************************************
 ****************************************************************************/
const (
	/**
	 * the chunk stream id used for some under-layer message,
	 * for example, the PC(protocol control) message.
	 */
	RTMP_CID_ProtocolControl = 0x02
	/**
	 * the AMF0/AMF3 command message, invoke method and return the result, over NetConnection.
	 * generally use 0x03.
	 */
	RTMP_CID_OverConnection = 0x03
	/**
	 * the AMF0/AMF3 command message, invoke method and return the result, over NetConnection,
	 * the midst state(we guess).
	 * rarely used, e.g. onStatus(NetStream.Play.Reset).
	 */
	RTMP_CID_OverConnection2 = 0x04
	/**
	 * the stream message(amf0/amf3), over NetStream.
	 * generally use 0x05.
	 */
	RTMP_CID_OverStream = 0x05
	/**
	 * the stream message(amf0/amf3), over NetStream, the midst state(we guess).
	 * rarely used, e.g. play("mp4:mystram.f4v")
	 */
	RTMP_CID_OverStream2 = 0x08
	/**
	 * the stream message(video), over NetStream
	 * generally use 0x06.
	 */
	RTMP_CID_Video = 0x06
	/**
	 * the stream message(audio), over NetStream.
	 * generally use 0x07.
	 */
	RTMP_CID_Audio = 0x07
)

/**
 * the rtmp client type.
 */
type SrsRtmpConnType int

const (
	_                           SrsRtmpConnType = iota
	SrsRtmpConnPlay                             = 0
	SrsRtmpConnFMLEPublish                      = 1
	SrsRtmpConnFlashPublish                     = 2
	SrsRtmpConnHaivisionPublish                 = 3
)

/**
 * 6.1. Chunk Format
 * Extended timestamp: 0 or 4 bytes
 * This field MUST be sent when the normal timsestamp is set to
 * 0xffffff, it MUST NOT be sent if the normal timestamp is set to
 * anything else. So for values less than 0xffffff the normal
 * timestamp field SHOULD be used in which case the extended timestamp
 * MUST NOT be present. For values greater than or equal to 0xffffff
 * the normal timestamp field MUST NOT be used and MUST be set to
 * 0xffffff and the extended timestamp MUST be sent.
 */
const RTMP_EXTENDED_TIMESTAMP = 0xFFFFFF
