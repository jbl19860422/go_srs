package app

import (
	"go_srs/srs/codec"
)

type SrsTsContext struct {
	ready      bool
	pids       map[int]*SrsTsChannel
	pure_audio bool
	vcodec     codec.SrsCodecVideo
	acodec     codec.SrsCodecAudio
}

func NewSrsTsContext() *SrsTsContext {
	return &SrsTsContext{}
}
