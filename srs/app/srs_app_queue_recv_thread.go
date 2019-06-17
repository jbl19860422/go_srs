package app

import (
	"go_srs/srs/protocol/rtmp"
)

type SrsQueueRecvThread struct {
	queue []*rtmp.SrsRtmpMessage
}
