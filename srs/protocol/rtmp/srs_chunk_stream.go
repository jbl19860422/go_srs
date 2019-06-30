package rtmp

const (
	RTMP_FMT_TYPE0 = 0
	RTMP_FMT_TYPE1 = 1
	RTMP_FMT_TYPE2 = 2
	RTMP_FMT_TYPE3 = 3
)

type SrsChunkStream struct {
	/**
	 * represents the basic header fmt,
	 * which used to identify the variant message header type.
	 */
	Format byte
	/**
	 * represents the basic header cid,
	 * which is the chunk stream id.
	 */
	Cid int32
	/**
	 * cached message header
	 */
	Header SrsMessageHeader
	/**
	 * whether the chunk message header has extended timestamp.
	 */
	ExtendedTimestamp bool

	MsgCount int32

	RtmpMessage *SrsRtmpMessage
}

func NewSrsChunkStream(cid_ int32) *SrsChunkStream {
	s := &SrsChunkStream{
		Format:            0,
		Cid:               cid_,
		ExtendedTimestamp: false,
		RtmpMessage:       nil,
		MsgCount:          0,
	}
	return s
}
