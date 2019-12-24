package app

import (
	"go_srs/srs/protocol/rtmp"
)

type SrsRtmpJitterAlgorithm int

const (
	_                          SrsRtmpJitterAlgorithm = iota
	SrsRtmpJitterAlgorithmFULL                        = 0x01
	SrsRtmpJitterAlgorithmZERO                        = 0x02
	SrsRtmpJitterAlgorithmOFF                         = 0x03
)

const CONST_MAX_JITTER_MS_NEG = -250
const CONST_MAX_JITTER_MS = 250
const DEFAULT_FRAME_TIME_MS = 10

type SrsRtmpJitter struct {
	lastPktTime        int64
	lastPktCorrectTime int64
}

func NewSrsRtmpJitter() *SrsRtmpJitter {
	return &SrsRtmpJitter{
		lastPktCorrectTime: -1,
		lastPktTime:        0,
	}
}

func (this *SrsRtmpJitter) Correct(msg *rtmp.SrsRtmpMessage, ag SrsRtmpJitterAlgorithm) error {
	if ag != SrsRtmpJitterAlgorithmFULL {
		if ag == SrsRtmpJitterAlgorithmOFF {
			return nil
		}
		// start at zero, but donot ensure monotonically increasing.
		if ag == SrsRtmpJitterAlgorithmZERO {
			if this.lastPktCorrectTime == -1 {
				this.lastPktCorrectTime = msg.GetHeader().GetTimestamp()
			}
			msg.GetHeader().SetTimestamp(msg.GetHeader().GetTimestamp() - this.lastPktCorrectTime)
			return nil
		}
		return nil
	}

	if msg.GetHeader().IsAV() {
		msg.GetHeader().SetTimestamp(0)
		return nil
	}

	/**
	* we use a very simple time jitter detect/correct algorithm:
	* 1. delta: ensure the delta is positive and valid,
	*     we set the delta to DEFAULT_FRAME_TIME_MS,
	*     if the delta of time is nagative or greater than CONST_MAX_JITTER_MS.
	* 2. last_pkt_time: specifies the original packet time,
	*     is used to detect next jitter.
	* 3. last_pkt_correct_time: simply add the positive delta,
	*     and enforce the time monotonically.
	 */

	timestamp := msg.GetHeader().GetTimestamp()
	delta := timestamp - this.lastPktTime
	if delta < CONST_MAX_JITTER_MS_NEG || delta > CONST_MAX_JITTER_MS {
		delta = DEFAULT_FRAME_TIME_MS
	}

	if this.lastPktCorrectTime+delta > 0 {
		this.lastPktCorrectTime = this.lastPktCorrectTime + delta
	}
	msg.GetHeader().SetTimestamp(this.lastPktCorrectTime)
	this.lastPktTime = timestamp
	return nil
}

func (this *SrsRtmpJitter) GetTime() int64 {
	return this.lastPktCorrectTime
}
