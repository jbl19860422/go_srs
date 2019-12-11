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
package global

const SRS_CONSTS_RTMP_PROTOCOL_CHUNK_SIZE = 128

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

// default vhost of rtmp
const SRS_CONSTS_RTMP_DEFAULT_VHOST = "__defaultVhost__"

// default port of rtmp
const SRS_CONSTS_RTMP_DEFAULT_PORT = "1935"

const RTMP_SIG_FMS_VER = "3,5,3,888"
const RTMP_SIG_AMF0_VER = 0
const RTMP_SIG_CLIENT_ID = "ASAICiss"

// FMLE
const RTMP_AMF0_COMMAND_ON_FC_PUBLISH = "onFCPublish"
const RTMP_AMF0_COMMAND_ON_FC_UNPUBLISH = "onFCUnpublish"

/**
 * onStatus consts.
 */
const (
	StatusLevel       = "level"
	StatusCode        = "code"
	StatusDescription = "description"
	StatusDetails     = "details"
	StatusClientId    = "clientid"
	// status value
	StatusLevelStatus = "status"
	// status error
	StatusLevelError = "error"
	// code value
	StatusCodeConnectSuccess   = "NetConnection.Connect.Success"
	StatusCodeConnectRejected  = "NetConnection.Connect.Rejected"
	StatusCodeStreamReset      = "NetStream.Play.Reset"
	StatusCodeStreamStart      = "NetStream.Play.Start"
	StatusCodeStreamPause      = "NetStream.Pause.Notify"
	StatusCodeStreamUnpause    = "NetStream.Unpause.Notify"
	StatusCodePublishStart     = "NetStream.Publish.Start"
	StatusCodeDataStart        = "NetStream.Data.Start"
	StatusCodeUnpublishSuccess = "NetStream.Unpublish.Success"
)

// provider info.
const RTMP_SIG_SRS_KEY = "SRS"
const RTMP_SIG_SRS_CODE = "ZhouGuowen"
const RTMP_SIG_SRS_AUTHROS = "winlin,wenjie.zhao"

// contact info.
const RTMP_SIG_SRS_WEB = "http://ossrs.net"
const RTMP_SIG_SRS_EMAIL = "winlin@vip.126.com"

// debug info.
const RTMP_SIG_SRS_ROLE = "cluster"
const RTMP_SIG_SRS_NAME = "SRS(Simple RTMP Server)"
const RTMP_SIG_SRS_URL_SHORT = "github.com/ossrs/srs"
const RTMP_SIG_SRS_URL = "https://github.com/ossrs/srs"
const RTMP_SIG_SRS_LICENSE = "The MIT License (MIT)"
const RTMP_SIG_SRS_COPYRIGHT = "Copyright (c) 2019 SRS(ossrs)"
const RTMP_SIG_SRS_PRIMARY = "SRS/2.0release"
const RTMP_SIG_SRS_CONTRIBUTORS_URL = "https://github.com/ossrs/srs/blob/master/AUTHORS.txt"
const RTMP_SIG_SRS_HANDSHAKE = "SRS(2.0.263)"
const RTMP_SIG_SRS_RELEASE = "https://github.com/ossrs/srs/tree/2.0release"
const RTMP_SIG_SRS_VERSION = "2.0.263"
const RTMP_SIG_SRS_SERVER = "SRS/2.0.263(ZhouGuowen)"

// 3.7. User Control message
type SrcPCUCEventType int

const (
	// generally, 4bytes event-data

	/**
	 * The server sends this event to notify the client
	 * that a stream has become functional and can be
	 * used for communication. By default, this event
	 * is sent on ID 0 after the application connect
	 * command is successfully received from the
	 * client. The event data is 4-byte and represents
	 * the stream ID of the stream that became
	 * functional.
	 */
	SrcPCUCStreamBegin = 0x00

	/**
	 * The server sends this event to notify the client
	 * that the playback of data is over as requested
	 * on this stream. No more data is sent without
	 * issuing additional commands. The client discards
	 * the messages received for the stream. The
	 * 4 bytes of event data represent the ID of the
	 * stream on which playback has ended.
	 */
	SrcPCUCStreamEOF = 0x01

	/**
	 * The server sends this event to notify the client
	 * that there is no more data on the stream. If the
	 * server does not detect any message for a time
	 * period, it can notify the subscribed clients
	 * that the stream is dry. The 4 bytes of event
	 * data represent the stream ID of the dry stream.
	 */
	SrcPCUCStreamDry = 0x02

	/**
	 * The client sends this event to inform the server
	 * of the buffer size (in milliseconds) that is
	 * used to buffer any data coming over a stream.
	 * This event is sent before the server starts
	 * processing the stream. The first 4 bytes of the
	 * event data represent the stream ID and the next
	 * 4 bytes represent the buffer length, in
	 * milliseconds.
	 */
	SrcPCUCSetBufferLength = 0x03 // 8bytes event-data

	/**
	 * The server sends this event to notify the client
	 * that the stream is a recorded stream. The
	 * 4 bytes event data represent the stream ID of
	 * the recorded stream.
	 */
	SrcPCUCStreamIsRecorded = 0x04

	/**
	 * The server sends this event to test whether the
	 * client is reachable. Event data is a 4-byte
	 * timestamp, representing the local server time
	 * when the server dispatched the command. The
	 * client responds with kMsgPingResponse on
	 * receiving kMsgPingRequest.
	 */
	SrcPCUCPingRequest = 0x06

	/**
	 * The client sends this event to the server in
	 * response to the ping request. The event data is
	 * a 4-byte timestamp, which was received with the
	 * kMsgPingRequest request.
	 */
	SrcPCUCPingResponse = 0x07

	/**
	 * for PCUC size=3, the payload is "00 1A 01",
	 * where we think the event is 0x001a, fms defined msg,
	 * which has only 1bytes event data.
	 */
	SrsPCUCFmsEvent0 = 0x1a
)
