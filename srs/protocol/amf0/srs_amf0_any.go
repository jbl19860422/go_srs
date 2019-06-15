package amf0

import (
	"utils"
)

type SrsAmf0Any {
	Decode(stream *utils.SrsStream) error
	Encode(stream *utils.SrsStream) error
	IsMyType(stream *utils.SrsStream) (bool, error)
}