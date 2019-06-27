package app

import (
	// "fmt"
	"encoding/binary"
	"go_srs/srs/utils"
)

/*
* the adaptation_field_control of ts packet
* Table 2-5 - Adaptation field control values hls-mpeg-ts-iso13818-1.pdf, page 38
 */
type SrsTsAdapationControl int

const (
	_                                SrsTsAdapationControl = iota
	SrsTsAdapationControlReserved                          = 0
	SrsTsAdapationControlPayloadOnly                       = 1
	SrsTsAdapationControlFieldOnly                         = 2
	SrsTsAdapationControlBoth                              = 3
)

type SrsTsHeader struct {
	/*
		The sync_byte is a fixed 8-bit field whose value is '0100 0111' (0x47). Sync_byte emulation in the choice of
		values for other regularly occurring fields, such as PID, should be avoided.
	*/
	syncByte                int8 //8bit 同步字节，固定为0x47
	transportErrorIndicator int8 //1bit 传输错误指示符，表明在ts头的adapt域后由一个无用字节，通常都为0，这个字节算在adapt域长度内
	/*
		The payload_unit_start_indicator is a 1-bit flag which has normative meaning for
		Transport Stream packets that carry PES packets (refer to 2.4.3.6) or PSI data (refer to 2.4.4).
		When the payload of the Transport Stream packet contains PES packet data, the payload_unit_start_indicator has the
		following significance: a '1' indicates that the payload of this Transport Stream packet will commence with the first byte
		of a PES packet and a '0' indicates no PES packet shall start in this Transport Stream packet. If the
		payload_unit_start_indicator is set to '1', then one and only one PES packet starts in this Transport Stream packet. This
		also applies to private streams of stream_type 6 (refer to Table 2-29).
		When the payload of the Transport Stream packet contains PSI data, the payload_unit_start_indicator has the following
		significance: if the Transport Stream packet carries the first byte of a PSI section, the payload_unit_start_indicator value
		shall be '1', indicating that the first byte of the payload of this Transport Stream packet carries the pointer_field. If the
		Transport Stream packet does not carry the first byte of a PSI section, the payload_unit_start_indicator value shall be '0',
		indicating that there is no pointer_field in the payload. Refer to 2.4.4.1 and 2.4.4.2. This also applies to private streams of
		stream_type 5 (refer to Table 2-29)
	*/
	payloadUnitStartIndicator int8 //1bit
	/*
		The transport_priority is a 1-bit indicator. When set to '1' it indicates that the associated packet is
		of greater priority than other packets having the same PID which do not have the bit set to '1'. The transport mechanism
		can use this to prioritize its data within an elementary stream. Depending on the application the transport_priority field
		may be coded regardless of the PID or within one PID only. This field may be changed by channel specific encoders or
		decoders
	*/
	transportPriority int8 //通常为0
	/*
		The PID is a 13-bit field, indicating the type of the data stored in the packet payload. PID value 0x0000 is
		reserved for the Program Association Table (see Table 2-25). PID value 0x0001 is reserved for the Conditional Access
		Table (see Table 2-27). PID values 0x0002 – 0x000F are reserved. PID value 0x1FFF is reserved for null packets (see
		Table 2-3).
	*/
	PID SrsTsPid //13bit
	/*
		This 2-bit field indicates the scrambling mode of the Transport Stream packet payload.
		The Transport Stream packet header, and the adaptation field when present, shall not be scrambled. In the case of a null
		packet the value of the transport_scrambling_control field shall be set to '00' (see Table 2-4).
	*/
	transportScrambingControl SrsTsScrambled
	/*
		This 2-bit field indicates whether this Transport Stream packet header is followed by an
		adaptation field and/or payload (see Table 2-5).
	*/
	adaptationFieldControl SrsTsAdapationControl

	/*
		The continuity_counter is a 4-bit field incrementing with each Transport Stream packet with the
		same PID. The continuity_counter wraps around to 0 after its maximum value. The continuity_counter shall not be
		incremented when the adaptation_field_control of the packet equals '00' or '10'.

		The continuity counter may be discontinuous when the discontinuity_indicator is set to '1' (refer to 2.4.3.4). In the case of
		a null packet the value of the continuity_counter is undefined.
	*/
	continuityCounter int8
}

func NewSrsTsHeader() *SrsTsHeader {
	return &SrsTsHeader{
		syncByte:                  0x47,
		transportErrorIndicator:   0,
		payloadUnitStartIndicator: 1,
		transportPriority:         0,
		PID:                       SrsTsPidPAT,
		transportScrambingControl: SrsTsScrambledDisabled,
		adaptationFieldControl:    SrsTsAdapationControlPayloadOnly,
		continuityCounter:         0,
	}
}

func (this *SrsTsHeader) Encode(stream *utils.SrsStream) {//4B
	stream.WriteByte(byte(this.syncByte))

	var pidv int16 = 0
	pidv = int16(this.PID) & 0x1FFF
	pidv |= int16(this.transportPriority<<13) & 0x2000
	pidv |= int16(int16(this.payloadUnitStartIndicator)<<14) & 0x4000
	pidv |= int16(uint16(this.transportErrorIndicator<<15) & 0x8000)
	stream.WriteInt16(pidv, binary.BigEndian)

	var b byte = 0
	b |= byte(this.continuityCounter&0x0f)
	b |= byte(this.transportScrambingControl << 6) & 0xC0
	b |= byte((this.adaptationFieldControl << 4) & 0x30)
	stream.WriteByte(b)
}

func (this *SrsTsHeader) Decode(stream *utils.SrsStream) error {
	return nil
}

func (this *SrsTsHeader) Size() uint32 {
	return 4
}
