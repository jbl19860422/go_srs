package packet

import (
	"encoding/binary"
	"go_srs/srs/global"
	"go_srs/srs/utils"
	_ "log"
)

type SrsSetChunkSizePacket struct {
	/**
	 * The maximum chunk size can be 65536 bytes. The chunk size is
	 * maintained independently for each direction.
	 */
	ChunkSize int32
}

func NewSrsSetChunkSizePacket() *SrsSetChunkSizePacket {
	return &SrsSetChunkSizePacket{
		ChunkSize: global.SRS_CONSTS_RTMP_PROTOCOL_CHUNK_SIZE,
	}
}

func (this *SrsSetChunkSizePacket) GetMessageType() int8 {
	return global.RTMP_MSG_SetChunkSize
}

func (this *SrsSetChunkSizePacket) GetPreferCid() int32 {
	return global.RTMP_CID_ProtocolControl
}

func (this *SrsSetChunkSizePacket) Decode(stream *utils.SrsStream) error {
	var err error
	this.ChunkSize, err = stream.ReadInt32(binary.BigEndian)
	return err
}

func (this *SrsSetChunkSizePacket) Encode(stream *utils.SrsStream) error {
	stream.WriteInt32(this.ChunkSize, binary.BigEndian)
	return nil
}
