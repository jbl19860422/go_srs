package app
import (
	"go_srs/srs/codec"
)

type SrsTsMuxer struct {

}

func NewSrsTsMuxer(c *SrsHlsCacheWriter, ac codec.SrsCodecAudio, vc codec.SrsCodecVideo) *SrsTsMuxer {
	return &SrsTsMuxer{}
}