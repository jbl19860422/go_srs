package app

type SrsTsContext struct {
	ready      bool
	pids       map[int]*SrsTsChannel
	pure_audio bool
	vcodec     SrsCodecVideo
	acodec     SrsCodecAudio
}

func NewSrsTsContext() *SrsTsContext {
	return &SrsTsContext{}
}
