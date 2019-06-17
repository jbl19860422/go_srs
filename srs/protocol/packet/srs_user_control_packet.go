package packet

import "encoding/binary"

// import "log"
import (
	"go_srs/srs/global"
	"go_srs/srs/utils"
)

type SrcPCUCEventType int

// 3.7. User Control message
const (
	// generally, 4bytes event-data
	_ SrcPCUCEventType = iota
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

/**
* 5.4. User Control Message (4)
*
* for the EventData is 4bytes.
* Stream Begin(=0)              4-bytes stream ID
* Stream EOF(=1)                4-bytes stream ID
* StreamDry(=2)                 4-bytes stream ID
* SetBufferLength(=3)           8-bytes 4bytes stream ID, 4bytes buffer length.
* StreamIsRecorded(=4)          4-bytes stream ID
* PingRequest(=6)               4-bytes timestamp local server time
* PingResponse(=7)              4-bytes timestamp received ping request.
*
* 3.7. User Control message
* +------------------------------+-------------------------
* | Event Type ( 2- bytes ) | Event Data
* +------------------------------+-------------------------
* Figure 5 Pay load for the 'User Control Message'.
 */

type SrsUserControlPacket struct {
	/**
	 * Event type is followed by Event data.
	 * @see: SrcPCUCEventType
	 */
	EventType int16

	/**
	 * the event data generally in 4bytes.
	 * @remark for event type is 0x001a, only 1bytes.
	 * @see SrsPCUCFmsEvent0
	 */
	EventData int32

	/**
	 * 4bytes if event_type is SetBufferLength; otherwise 0.
	 */
	ExtraData int32
}

func NewSrsUserControlPacket() *SrsUserControlPacket {
	return &SrsUserControlPacket{}
}

func (this *SrsUserControlPacket) Decode(stream *utils.SrsStream) (err error) {
	if this.EventType, err = stream.ReadInt16(binary.BigEndian); err != nil {
		return
	}

	if this.EventType == SrsPCUCFmsEvent0 {
		var d int8
		d, err = stream.ReadInt8()
		if err != nil {
			return
		}
		this.EventData = int32(d)
	} else {
		if this.EventData, err = stream.ReadInt32(binary.BigEndian); err != nil {
			return
		}
	}

	if this.EventType == SrcPCUCSetBufferLength {
		if this.ExtraData, err = stream.ReadInt32(binary.BigEndian); err != nil {
			return err
		}
	}
	err = nil
	return
}

func (this *SrsUserControlPacket) Encode(stream *utils.SrsStream) error {
	stream.WriteInt16(this.EventType, binary.BigEndian)
	if this.EventType == SrsPCUCFmsEvent0 {
		stream.WriteByte(byte(this.EventData))
	} else {
		stream.WriteInt32(this.EventData, binary.BigEndian)
	}

	if this.EventType == SrcPCUCSetBufferLength {
		stream.WriteInt32(this.ExtraData, binary.BigEndian)
	}
	return nil
}


func (this *SrsUserControlPacket) GetPreferCid() int32 {
    return global.RTMP_CID_ProtocolControl
}

func (this *SrsUserControlPacket) GetMessageType() int8 {
    return global.RTMP_MSG_UserControlMessage
}
