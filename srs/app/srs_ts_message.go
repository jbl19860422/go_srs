package app

import (
	"go_srs/srs/utils"
)

/**
* the pid of ts packet,
* Table 2-3 - PID table, hls-mpeg-ts-iso13818-1.pdf, page 37
* NOTE - The transport packets with PID values 0x0000, 0x0001, and 0x0010-0x1FFE are allowed to carry a PCR.
 */
type SrsTsPid int

const (
	_ SrsTsPid = iota
	// Program Association Table(see Table 2-25).
	SrsTsPidPAT = 0x00
	// Conditional Access Table (see Table 2-27).
	SrsTsPidCAT = 0x01
	// Transport Stream Description Table
	SrsTsPidTSDT = 0x02
	// Reserved
	SrsTsPidReservedStart = 0x03
	SrsTsPidReservedEnd   = 0x0f
	// May be assigned as network_PID, Program_map_PID, elementary_PID, or for other purposes
	SrsTsPidAppStart = 0x10
	SrsTsPidAppEnd   = 0x1ffe
	// null packets (see Table 2-3)
	SrsTsPidNULL = 0x01FFF
)

/**
* the transport_scrambling_control of ts packet,
* Table 2-4 - Scrambling control values, hls-mpeg-ts-iso13818-1.pdf, page 38
 */
type SrsTsScrambled int

const (
	_ SrsTsScrambled = iota
	// Not scrambled
	SrsTsScrambledDisabled = 0x00
	// User-defined
	SrsTsScrambledUserDefined1 = 0x01
	// User-defined
	SrsTsScrambledUserDefined2 = 0x02
	// User-defined
	SrsTsScrambledUserDefined3 = 0x03
)

const SRS_TS_SYNC_BYTE = 0x47

const TS_PMT_NUMBER = 1
const TS_PMT_PID = 0x1001
const TS_VIDEO_AVC_PID = 0x100
const TS_AUDIO_AAC_PID = 0x101
const TS_AUDIO_MP3_PID = 0x102

type SrsTsPayload interface {
	Encode(stream *utils.SrsStream)
	Decode(stream *utils.SrsStream) error
	Size() uint32
}

type SrsTsMessage struct {
	channel *SrsTsChannel
	packet 	*SrsTsPacket
}

func (this *SrsTsMessage) IsAudio() bool {
	return false
}