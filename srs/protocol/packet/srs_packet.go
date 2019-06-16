package packet

import (
	"go_srs/srs/utils"
)

type SrsPacket interface {
	Decode(stream *utils.SrsStream) error
	Encode(stream *utils.SrsStream) error
	GetMessageType() int8
	GetPreferCid() int32
}
