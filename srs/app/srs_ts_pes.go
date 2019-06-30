/*
The MIT License (MIT)

Copyright (c) 2013-2015 GOSRS(gosrs)

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

//see iso-13818.pdf, page 49
type SrsTsPayloadPES struct {
	packet *SrsTsPacket
	/*
		The packet_start_code_prefix is a 24-bit code. Together with the stream_id that follows, it
		constitutes a packet start code that identifies the beginning of a packet. The packet_start_code_prefix is the bit string
		'0000 0000 0000 0000 0000 0001' (0x000001 in hexadecimal).
	*/
	packetStartCodePrefix int32 //24bit 0x000001
	/*
		In Program Streams, the stream_id specifies the type and number of the elementary stream as defined by the
		stream_id Table 2-18. In Transport Streams, the stream_id may be set to any valid value which correctly describes the
		elementary stream type as defined in Table 2-18. In Transport Streams, the elementary stream type is specified in the
		Program Specific Information as specified in 2.4.4.
	*/
	streamId int8 //音频取值（0xc0-0xdf），通常为0xc0 视频取值（0xe0-0xef），通常为0xe0
	/*
		The PES_packet_length is a 16-bit field indicating the total number of bytes in the
		program_stream_directory immediately following this field (refer to Table 2-18).
	*/
	PESPacketLength uint16

	const2Bits int8 //2bits '10'
	/*
		The 2-bit PES_scrambling_control field indicates the scrambling mode of the PES packet
		payload. When scrambling is performed at the PES level, the PES packet header, including the optional fields when
		present, shall not be scrambled (see Table 2-19)
	*/
	PESScramblingControl int8 //2bit
	/**
	 * This is a 1-bit field indicating the priority of the payload in this PES packet. A '1' indicates a higher
	 * priority of the payload of the PES packet payload than a PES packet payload with this field set to '0'. A multiplexor can
	 * use the PES_priority bit to prioritize its data within an elementary stream. This field shall not be changed by the transport
	 * mechanism.
	 */
	PESPriority int8 //1bit 一般为0
	/**
	 * This is a 1-bit flag. When set to a value of '1' it indicates that the PES packet header is
	 * immediately followed by the video start code or audio syncword indicated in the data_stream_alignment_descriptor
	 * in 2.6.10 if this descriptor is present. If set to a value of '1' and the descriptor is not present, alignment as indicated in
	 * alignment_type '01' in Table 2-47 and Table 2-48 is required. When set to a value of '0' it is not defined whether any such
	 * alignment occurs or not.
	 */
	dataAlignmentIndicator int8 //1bit 一般为0
	/**
	 * This is a 1-bit field. When set to '1' it indicates that the material of the associated PES packet payload is
	 * protected by copyright. When set to '0' it is not defined whether the material is protected by copyright. A copyright
	 * descriptor described in 2.6.24 is associated with the elementary stream which contains this PES packet and the copyright
	 * flag is set to '1' if the descriptor applies to the material contained in this PES packet
	 */
	copyright int8 //1bit 一般为0
	/**
	 * This is a 1-bit field. When set to '1' the contents of the associated PES packet payload is an original.
	 * When set to '0' it indicates that the contents of the associated PES packet payload is a copy.
	 */
	originalOrCopy int8 //1bit 一般为1

	// 1B
	/**
	 * This is a 2-bit field. When the PTS_DTS_flags field is set to '10', the PTS fields shall be present in
	 * the PES packet header. When the PTS_DTS_flags field is set to '11', both the PTS fields and DTS fields shall be present
	 * in the PES packet header. When the PTS_DTS_flags field is set to '00' no PTS or DTS fields shall be present in the PES
	 * packet header. The value '01' is forbidden.
	 */
	PTSDTSFlags int8 //2bits
	/**
	 * A 1-bit flag, which when set to '1' indicates that ESCR base and extension fields are present in the PES
	 * packet header. When set to '0' it indicates that no ESCR fields are present.
	 */
	ESCRFlag int8 //1bit	一般为 0吧
	/**
	 * A 1-bit flag, which when set to '1' indicates that the ES_rate field is present in the PES packet header.
	 * When set to '0' it indicates that no ES_rate field is present.
	 */
	ESRateFlag int8 //1bit	一般为0
	/**
	 * A 1-bit flag, which when set to '1' it indicates the presence of an 8-bit trick mode field. When
	 * set to '0' it indicates that this field is not present.
	 */
	DSMTrickModeFlag int8 //1bit
	/**
	 * A 1-bit flag, which when set to '1' indicates the presence of the additional_copy_info field.
	 * When set to '0' it indicates that this field is not present.
	 */
	additionalCopyInfoFlag int8 //1bit
	/**
	 * A 1-bit flag, which when set to '1' indicates that a CRC field is present in the PES packet. When set to
	 * '0' it indicates that this field is not present.
	 */
	PESCRCFlag int8 //1bit
	/**
	 * A 1-bit flag, which when set to '1' indicates that an extension field exists in this PES packet
	 * header. When set to '0' it indicates that this field is not present.
	 */
	PESExtensionFlag int8 //1bit

	// 1B
	/**
	 * An 8-bit field specifying the total number of bytes occupied by the optional fields and any
	 * stuffing bytes contained in this PES packet header. The presence of optional fields is indicated in the byte that precedes
	 * the PES_header_data_length field.
	 */
	PESHeaderDataLength uint8 //8bits

	// 5B
	/**
	 * Presentation times shall be related to decoding times as follows: The PTS is a 33-bit
	 * number coded in three separate fields. It indicates the time of presentation, tp n (k), in the system target decoder of a
	 * presentation unit k of elementary stream n. The value of PTS is specified in units of the period of the system clock
	 * frequency divided by 300 (yielding 90 kHz). The presentation time is derived from the PTS according to equation 2-11
	 * below. Refer to 2.7.4 for constraints on the frequency of coding presentation timestamps.
	 */
	// ===========1B
	// 4bits const
	// 3bits PTS [32..30]
	// 1bit const '1'
	// ===========2B
	// 15bits PTS [29..15]
	// 1bit const '1'
	// ===========2B
	// 15bits PTS [14..0]
	// 1bit const '1'
	pts int64 // 33bits

	// 5B
	/**
	 * The DTS is a 33-bit number coded in three separate fields. It indicates the decoding time,
	 * td n (j), in the system target decoder of an access unit j of elementary stream n. The value of DTS is specified in units of
	 * the period of the system clock frequency divided by 300 (yielding 90 kHz).
	 */
	// ===========1B
	// 4bits const
	// 3bits DTS [32..30]
	// 1bit const '1'
	// ===========2B
	// 15bits DTS [29..15]
	// 1bit const '1'
	// ===========2B
	// 15bits DTS [14..0]
	// 1bit const '1'
	dts int64 // 33bits

	// 6B
	/**
	 * The elementary stream clock reference is a 42-bit field coded in two parts. The first
	 * part, ESCR_base, is a 33-bit field whose value is given by ESCR_base(i), as given in equation 2-14. The second part,
	 * ESCR_ext, is a 9-bit field whose value is given by ESCR_ext(i), as given in equation 2-15. The ESCR field indicates the
	 * intended time of arrival of the byte containing the last bit of the ESCR_base at the input of the PES-STD for PES streams
	 * (refer to 2.5.2.4).
	 */
	// 2bits reserved
	// 3bits ESCR_base[32..30]
	// 1bit const '1'
	// 15bits ESCR_base[29..15]
	// 1bit const '1'
	// 15bits ESCR_base[14..0]
	// 1bit const '1'
	// 9bits ESCR_extension
	// 1bit const '1'
	ESCRBase      int64 //33bits
	ESCRExtension int16 //9bits

	// 3B
	/**
	 * The ES_rate field is a 22-bit unsigned integer specifying the rate at which the
	 * system target decoder receives bytes of the PES packet in the case of a PES stream. The ES_rate is valid in the PES
	 * packet in which it is included and in subsequent PES packets of the same PES stream until a new ES_rate field is
	 * encountered. The value of the ES_rate is measured in units of 50 bytes/second. The value 0 is forbidden. The value of the
	 * ES_rate is used to define the time of arrival of bytes at the input of a P-STD for PES streams defined in 2.5.2.4. The
	 * value encoded in the ES_rate field may vary from PES_packet to PES_packet.
	 */
	// 1bit const '1'
	// 22bits ES_rate
	// 1bit const '1'
	ESRate int32 //22bits

	// 1B
	/**
	 * A 3-bit field that indicates which trick mode is applied to the associated video stream. In cases of
	 * other types of elementary streams, the meanings of this field and those defined by the following five bits are undefined.
	 * For the definition of trick_mode status, refer to the trick mode section of 2.4.2.3.
	 */
	trickModeControl int8 //3bits
	trickModeValue   int8 //5bits

	// 1B
	// 1bit const '1'
	/**
	 * This 7-bit field contains private data relating to copyright information.
	 */
	additionalCopyInfo int8 //7bits

	// 2B
	/**
	 * The previous_PES_packet_CRC is a 16-bit field that contains the CRC value that yields
	 * a zero output of the 16 registers in the decoder similar to the one defined in Annex A,
	 */
	previousPESPacketCRC int8 //16bits

	// 1B
	/**
	 * A 1-bit flag which when set to '1' indicates that the PES packet header contains private data.
	 * When set to a value of '0' it indicates that private data is not present in the PES header.
	 */
	PESPrivateDataFlag int8 //1bit
	/**
	 * A 1-bit flag which when set to '1' indicates that an ISO/IEC 11172-1 pack header or a
	 * Program Stream pack header is stored in this PES packet header. If this field is in a PES packet that is contained in a
	 * Program Stream, then this field shall be set to '0'. In a Transport Stream, when set to the value '0' it indicates that no pack
	 * header is present in the PES header.
	 */
	packHeaderFieldFlag int8 //1bit
	/**
	 * A 1-bit flag which when set to '1' indicates that the
	 * program_packet_sequence_counter, MPEG1_MPEG2_identifier, and original_stuff_length fields are present in this
	 * PES packet. When set to a value of '0' it indicates that these fields are not present in the PES header.
	 */
	programPacketSequenceCounterFlag int8 //1bit
	/**
	 * A 1-bit flag which when set to '1' indicates that the P-STD_buffer_scale and P-STD_buffer_size
	 * are present in the PES packet header. When set to a value of '0' it indicates that these fields are not present in the
	 * PES header.
	 */
	PSTDBufferFlag int8 //1bit
	/**
	 * reverved value, must be '1'
	 */
	const1Value0 int8 //3bits
	/**
	 * A 1-bit field which when set to '1' indicates the presence of the PES_extension_field_length
	 * field and associated fields. When set to a value of '0' this indicates that the PES_extension_field_length field and any
	 * associated fields are not present.
	 */
	PESExtensionFlag2 int8 //1bit

	// 16B
	/**
	 * This is a 16-byte field which contains private data. This data, combined with the fields before and
	 * after, shall not emulate the packet_start_code_prefix (0x000001).
	 */
	PESPrivateData []byte //128bits

	// (1+x)B
	/**
	 * This is an 8-bit field which indicates the length, in bytes, of the pack_header_field().
	 */
	packFieldLength uint8  //8bits
	packField       []byte //[pack_field_length] bytes

	// 2B
	// 1bit const '1'
	/**
	 * The program_packet_sequence_counter field is a 7-bit field. It is an optional
	 * counter that increments with each successive PES packet from a Program Stream or from an ISO/IEC 11172-1 Stream or
	 * the PES packets associated with a single program definition in a Transport Stream, providing functionality similar to a
	 * continuity counter (refer to 2.4.3.2). This allows an application to retrieve the original PES packet sequence of a Program
	 * Stream or the original packet sequence of the original ISO/IEC 11172-1 stream. The counter will wrap around to 0 after
	 * its maximum value. Repetition of PES packets shall not occur. Consequently, no two consecutive PES packets in the
	 * program multiplex shall have identical program_packet_sequence_counter values.
	 */
	programPacketSequenceCounter int8 //7bits
	// 1bit const '1'
	/**
	 * A 1-bit flag which when set to '1' indicates that this PES packet carries information from
	 * an ISO/IEC 11172-1 stream. When set to '0' it indicates that this PES packet carries information from a Program Stream.
	 */
	MPEG1MPEG2Identifier int8 //1bit
	/**
	 * This 6-bit field specifies the number of stuffing bytes used in the original ITU-T
	 * Rec. H.222.0 | ISO/IEC 13818-1 PES packet header or in the original ISO/IEC 11172-1 packet header.
	 */
	originalStuffLength int8 //6bits

	// 2B
	// 2bits const '01'
	/**
	 * The P-STD_buffer_scale is a 1-bit field, the meaning of which is only defined if this PES packet
	 * is contained in a Program Stream. It indicates the scaling factor used to interpret the subsequent P-STD_buffer_size field.
	 * If the preceding stream_id indicates an audio stream, P-STD_buffer_scale shall have the value '0'. If the preceding
	 * stream_id indicates a video stream, P-STD_buffer_scale shall have the value '1'. For all other stream types, the value
	 * may be either '1' or '0'.
	 */
	PSTDBufferScale int8 //1bit
	/**
	 * The P-STD_buffer_size is a 13-bit unsigned integer, the meaning of which is only defined if this
	 * PES packet is contained in a Program Stream. It defines the size of the input buffer, BS n , in the P-STD. If
	 * P-STD_buffer_scale has the value '0', then the P-STD_buffer_size measures the buffer size in units of 128 bytes. If
	 * P-STD_buffer_scale has the value '1', then the P-STD_buffer_size measures the buffer size in units of 1024 bytes.
	 */
	PSTDBufferSize int16 //13bits

	// (1+x)B
	// 1bit const '1'
	/**
	 * This is a 7-bit field which specifies the length, in bytes, of the data following this field in
	 * the PES extension field up to and including any reserved bytes.
	 */
	PESExtensionFieldLength uint8 //7bits
	PESExtensionField       []byte

	// NB
	/**
	 * This is a fixed 8-bit value equal to '1111 1111' that can be inserted by the encoder, for example to meet
	 * the requirements of the channel. It is discarded by the decoder. No more than 32 stuffing bytes shall be present in one
	 * PES packet header.
	 */
	stuffingBytes []byte

	// NB
	/**
	 * PES_packet_data_bytes shall be contiguous bytes of data from the elementary stream
	 * indicated by the packet's stream_id or PID. When the elementary stream data conforms to ITU-T
	 * Rec. H.262 | ISO/IEC 13818-2 or ISO/IEC 13818-3, the PES_packet_data_bytes shall be byte aligned to the bytes of this
	 * Recommendation | International Standard. The byte-order of the elementary stream shall be preserved. The number of
	 * PES_packet_data_bytes, N, is specified by the PES_packet_length field. N shall be equal to the value indicated in the
	 * PES_packet_length minus the number of bytes between the last byte of the PES_packet_length field and the first
	 * PES_packet_data_byte.
	 *
	 * In the case of a private_stream_1, private_stream_2, ECM_stream, or EMM_stream, the contents of the
	 * PES_packet_data_byte field are user definable and will not be specified by ITU-T | ISO/IEC in the future.
	 */

	dataBytes []byte

	// NB
	/**
	 * This is a fixed 8-bit value equal to '1111 1111'. It is discarded by the decoder.
	 */
	paddings []byte
}

func NewSrsTsPayloadPES() *SrsTsPayloadPES {
	return &SrsTsPayloadPES{
		const2Bits: 0x02,
	}
}

func (this *SrsTsPayloadPES) Decode(stream *utils.SrsStream) error {
	return nil
}

func (this *SrsTsPayloadPES) Size() uint32 {
	return 0
}

func CreatePes(context *SrsTsContext, pid int16, sid SrsTsPESStreamId, continuityCounter *uint8, discontinuity int8, pcr int64, dts int64, pts int64, data []byte) []*SrsTsPacket {
	pkts := make([]*SrsTsPacket, 0)
	pes := NewSrsTsPayloadPES()
	pes.dataBytes = data
	if len(data) > 0xffff {
		pes.PESPacketLength = 0
	} else {
		pes.PESPacketLength = uint16(len(data))
	}

	pes.packetStartCodePrefix = 0x01
	pes.streamId = int8(sid)
	pes.PESScramblingControl = 0
	pes.PESPriority = 0
	pes.dataAlignmentIndicator = 0
	pes.copyright = 0
	pes.originalOrCopy = 0
	if dts == pts {
		pes.PTSDTSFlags = 0x02
	} else {
		pes.PTSDTSFlags = 0x03
	}
	pes.ESCRFlag = 0
	pes.ESRateFlag = 0
	pes.DSMTrickModeFlag = 0
	pes.additionalCopyInfoFlag = 0
	pes.PESCRCFlag = 0
	pes.PESExtensionFlag = 0
	pes.PESHeaderDataLength = 0 // calc in size.
	pes.pts = pts
	pes.dts = dts

	s := utils.NewSrsStream([]byte{})
	pes.Encode(s)
	payload := s.Data()
	leftCount := len(payload)
	var currPos int = 0
	for leftCount > 0 {
		pkt := NewSrsTsPacket()

		pkt.tsHeader.syncByte = SRS_TS_SYNC_BYTE
		pkt.tsHeader.transportErrorIndicator = 0
		pkt.tsHeader.payloadUnitStartIndicator = 0
		if leftCount == len(payload) {
			pkt.tsHeader.payloadUnitStartIndicator = 1
		}
		pkt.tsHeader.transportPriority = 0
		pkt.tsHeader.PID = SrsTsPid(pid)
		pkt.tsHeader.transportScrambingControl = SrsTsScrambledDisabled
		pkt.tsHeader.adaptationFieldControl = SrsTsAdapationControlPayloadOnly
		pkt.tsHeader.continuityCounter = int8(*continuityCounter)
		*continuityCounter++
		var spaceLeft int = 0
		if leftCount == len(payload) {
			if pcr >= 0 {
				af := NewSrsTsAdaptationField(pkt)
				pkt.adaptationField = af
				pkt.tsHeader.adaptationFieldControl = SrsTsAdaptationFieldTypeBoth

				af.adaptationFieldLength = 0 // calc in size.
				af.discontinuityIndicator = discontinuity
				af.randomAccessIndicator = 0
				af.elementaryStreamPriorityIndicator = 0
				af.PCRFlag = 1
				af.OPCRFlag = 0
				af.splicingPointFlag = 0
				af.transportPrivateDataFlag = 0
				af.adaptationFieldExtensionFlag = 0
				af.programClockReferenceBase = pcr
				af.programClockReferenceExtension = 0

				spaceLeft = int(188 - 4 - af.Size())
			} else {
				spaceLeft = 188 - 4
			}
		} else {
			spaceLeft = 188 - 4
		}

		if leftCount < spaceLeft {
			af := NewSrsTsAdaptationField(pkt)
			pkt.adaptationField = af
			pkt.tsHeader.adaptationFieldControl = SrsTsAdaptationFieldTypeBoth

			af.adaptationFieldLength = 0 // calc in size.
			af.discontinuityIndicator = 0
			af.randomAccessIndicator = 0
			af.elementaryStreamPriorityIndicator = 0
			af.PCRFlag = 0
			af.OPCRFlag = 0
			af.splicingPointFlag = 0
			af.transportPrivateDataFlag = 0
			af.adaptationFieldExtensionFlag = 0
			af.programClockReferenceBase = 0
			af.programClockReferenceExtension = 0
		}

		var adaptationFieldLength int = 0
		if pkt.adaptationField != nil {
			adaptationFieldLength = pkt.adaptationField.Padding(leftCount)
		}

		consumed := 188 - 4 - adaptationFieldLength
		pkt.payload = payload[currPos:(currPos + consumed)]
		currPos += consumed
		leftCount -= consumed
		pkts = append(pkts, pkt)
	}
	return pkts
}

func (this *SrsTsPayloadPES) Encode(stream *utils.SrsStream) {
	this.PESHeaderDataLength = 0
	if this.PTSDTSFlags == 0x03 {
		this.PESHeaderDataLength = 10
	} else if this.PTSDTSFlags == 0x02 {
		this.PESHeaderDataLength = 5
	}

	//start code
	stream.WriteByte(0x00)
	stream.WriteByte(0x00)
	stream.WriteByte(0x01)
	//stream id(1B)
	stream.WriteByte(byte(this.streamId))
	// 2B
	// the PES_packet_length is the actual bytes size, the pplv write to ts
	// is the actual bytes plus the header size.
	var pplv int32 = 0
	if this.PESPacketLength > 0 {
		pplv = int32(this.PESPacketLength + 3 + uint16(this.PESHeaderDataLength))
		if pplv > 0xFFFF {
			pplv = 0
		}
	}
	stream.WriteInt16(int16(pplv), binary.BigEndian)

	var oocv int8 = this.originalOrCopy & 0x01
	oocv |= int8(int32(this.const2Bits<<6) & 0xC0)
	oocv |= (this.PESScramblingControl << 4) & 0x30
	oocv |= (this.PESPriority << 3) & 0x08
	oocv |= (this.dataAlignmentIndicator << 2) & 0x04
	oocv |= (this.copyright << 1) & 0x02

	stream.WriteByte(byte(oocv))

	// 1B
	var pefv int8 = this.PESExtensionFlag & 0x01
	pefv |= int8(int32(this.PTSDTSFlags<<6) & 0xC0)
	pefv |= (this.ESCRFlag << 5) & 0x20
	pefv |= (this.ESRateFlag << 4) & 0x10
	pefv |= (this.DSMTrickModeFlag << 3) & 0x08
	pefv |= (this.additionalCopyInfoFlag << 2) & 0x04
	pefv |= (this.PESCRCFlag << 1) & 0x02
	stream.WriteByte(byte(pefv))
	// 1B
	stream.WriteByte(this.PESHeaderDataLength)

	if this.PTSDTSFlags == 0x02 {
		this.encode_33bits_dts_pts(stream, 0x02, this.pts)
	}

	// fmt.Println("PTSDTSFlags=", this.PTSDTSFlags)
	// os.Exit(0)
	if this.PTSDTSFlags == 0x03 {
		this.encode_33bits_dts_pts(stream, 0x03, this.pts)
		this.encode_33bits_dts_pts(stream, 0x03, this.dts)
		//todo message
	}

	stream.WriteBytes(this.dataBytes)
}

func (this *SrsTsPayloadPES) encode_33bits_dts_pts(stream *utils.SrsStream, fb uint8, v int64) {
	var val int32 = 0
	val = int32(int64(fb)<<4 | (((v >> 30) & 0x07) << 1) | 1)
	stream.WriteByte(byte(val))

	val = int32((((v >> 15) & 0x7fff) << 1) | 1)
	stream.WriteByte(byte(val >> 8))
	stream.WriteByte(byte(val))

	val = int32((((v) & 0x7fff) << 1) | 1)
	stream.WriteByte(byte(val >> 8))
	stream.WriteByte(byte(val))
}
