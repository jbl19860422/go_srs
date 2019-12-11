/*
The MIT License (MIT)

Copyright (c) 2019 GOSRS(gosrs)

Permission is hereby granted, free of charge, to any person obtaining a copy of
this software and associated documentation files (the "Software"), to deal in
the Software without restriction, including without limitation the rights to
use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
the Software, and to permit persons to whom the Software is furnished to do so,
subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

package app

import (
	"encoding/binary"
	"go_srs/srs/utils"
)

/**
* 2.4.4.4 Table_id assignments, hls-mpeg-ts-iso13818-1.pdf, page 62
* The table_id field identifies the contents of a Transport Stream PSI section as shown in Table 2-26.
 */
type SrsTsPsiTableId int

const (
	_ SrsTsPsiTableId = iota
	// program_association_section
	SrsTsPsiTableIdPas = 0x00
	// conditional_access_section (CA_section)
	SrsTsPsiTableIdCas = 0x01
	// TS_program_map_section
	SrsTsPsiTableIdPms = 0x02
	// TS_description_section
	SrsTsPsiTableIdDs = 0x03
	// ISO_IEC_14496_scene_description_section
	SrsTsPsiTableIdSds = 0x04
	// ISO_IEC_14496_object_descriptor_section
	SrsTsPsiTableIdOds = 0x05
	// ITU-T Rec. H.222.0 | ISO/IEC 13818-1 reserved
	SrsTsPsiIdTableIso138181Start = 0x06
	SrsTsPsiIdTableIso138181End   = 0x37
	// Defined in ISO/IEC 13818-6
	SrsTsPsiIdTableIso138186Start = 0x38
	SrsTsPsiIdTableIso138186End   = 0x3F
	// User private
	SrsTsPsiTableIdUserStart = 0x40
	SrsTsPsiTableIdUserEnd   = 0xFE
	// forbidden
	SrsTsPsiTableIdForbidden = 0xFF
)

type SrsTsPayloadPSI struct {
	packet 		*SrsTsPacket
	/**
	 * This is an 8-bit field whose value shall be the number of bytes, immediately following the pointer_field
	 * until the first byte of the first section that is present in the payload of the Transport Stream packet (so a value of 0x00 in
	 * the pointer_field indicates that the section starts immediately after the pointer_field). When at least one section begins in
	 * a given Transport Stream packet, then the payload_unit_start_indicator (refer to 2.4.3.2) shall be set to 1 and the first
	 * byte of the payload of that Transport Stream packet shall contain the pointer. When no section begins in a given
	 * Transport Stream packet, then the payload_unit_start_indicator shall be set to 0 and no pointer shall be sent in the
	 * payload of that packet.
	 */
	pointerField int8
	// 1B
	/**
	 * This is an 8-bit field, which shall be set to 0x00 as shown in Table 2-26.
	 */
	tableId SrsTsPsiTableId //PAT表固定为0x00,PMT表为0x02
	/*
	 The section_syntax_indicator is a 1-bit field which shall be set to '1'.
	*/
	sectionSyntaxIndicator int8 //固定为二进制1
	/**
	 * const value, must be '0'
	 */
	const0Value int8 //1bit
	/**
	 * reverved value, must be '1'
	 */
	const1Value0 int8 //2bits
	/*
	 This is a 12-bit field, the first two bits of which shall be '00'. The remaining 10 bits specify the number
	 of bytes of the section, starting immediately following the section_length field, and including the CRC. The value in this
	 field shall not exceed 1021 (0x3FD).
	*/
	sectionLength int16 //12bits 	后面数据的长度
}

func NewSrsTsPayloadPSI(p *SrsTsPacket) *SrsTsPayloadPSI {
	return &SrsTsPayloadPSI{
		pointerField: 0,
		packet:p,
	}
}

func (this *SrsTsPayloadPSI) Encode(stream *utils.SrsStream) {//4B
	if this.packet.tsHeader.payloadUnitStartIndicator == 1 {
		stream.WriteByte(byte(this.pointerField))
	}
	
	stream.WriteByte(byte(this.tableId))

	var slv int16 = 0
	slv |= this.sectionLength & 0x0fff
	this.const1Value0 = 0x3
	slv |= (int16(this.const1Value0) << 12) & 0x3000
	this.const0Value = 0
	slv |= (int16(this.const0Value) << 14) & 0x4000
	slv |= int16(this.sectionSyntaxIndicator & 0x01) << 15
	stream.WriteInt16(slv, binary.BigEndian)
}

func (this *SrsTsPayloadPSI) Decode(stream *utils.SrsStream) error {
	return nil
}
