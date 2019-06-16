package rtmp

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
