package protocol

// "fmt"

const (
	RTMP_FMT_TYPE0 = 0
	RTMP_FMT_TYPE1 = 1
	RTMP_FMT_TYPE2 = 2
	RTMP_FMT_TYPE3 = 3
)

const RTMP_EXTENDED_TIMESTAMP = 0xFFFFFF

type SrsChunkStream struct {
	/**
	 * represents the basic header fmt,
	 * which used to identify the variant message header type.
	 */
	fmt byte
	/**
	 * represents the basic header cid,
	 * which is the chunk stream id.
	 */
	cid int32
	/**
	 * cached message header
	 */
	Header SrsMessageHeader
	/**
	 * whether the chunk message header has extended timestamp.
	 */
	extendedTimestamp bool

	msgCount int32

	RtmpMessage *SrsRtmpMessage
}

func NewSrsChunkStream(cid_ int32) *SrsChunkStream {
	s := &SrsChunkStream{
		fmt:               0,
		cid:               cid_,
		extendedTimestamp: false,
		RtmpMessage:       nil,
		msgCount:          0,
	}
	return s
}
