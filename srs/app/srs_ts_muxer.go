package app

import (
	"os"
	"go_srs/srs/codec"
)

type SrsTsMuxer struct {
	acodec codec.SrsCodecAudio
	vcodec codec.SrsCodecVideo

	context	*SrsTsContext
	writer	*os.File
	path	string
}

func NewSrsTsMuxer(f *os.File, c *SrsTsContext, ac codec.SrsCodecAudio, vc codec.SrsCodecVideo) *SrsTsMuxer {
	return &SrsTsMuxer{
		writer:f,
		context:c,
		acodec:ac,
		vcodec:vc,
	}
}

func (this *SrsTsMuxer) open(p string) error {
	this.path = p
	this.close()
	this.context.Reset()
	return nil
}

func (this *SrsTsMuxer) close() {
	this.writer.Close()
}

func (this *SrsTsMuxer) update_acodec(ac codec.SrsCodecAudio) error {
	this.acodec = ac
	return nil
}
