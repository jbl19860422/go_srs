package app

import (
	"go_srs/srs/protocol/rtmp"
)
type Consumer interface {
	PlayCycle() error
	StopPlay() error
	OnRecvError(err error)
	Enqueue(msg *rtmp.SrsRtmpMessage, atc bool, jitterAlgorithm *SrsRtmpJitterAlgorithm)
}
