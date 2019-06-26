package app

/**
* the pid of ts packet,
* Table 2-3 - PID table, hls-mpeg-ts-iso13818-1.pdf, page 37
* NOTE - The transport packets with PID values 0x0000, 0x0001, and 0x0010-0x1FFE are allowed to carry a PCR.
*/
type SrsTsPid	int
const (
	_ SrsTsPid = iota
	// Program Association Table(see Table 2-25).
    SrsTsPidPAT             = 0x00
    // Conditional Access Table (see Table 2-27).
    SrsTsPidCAT             = 0x01
    // Transport Stream Description Table
    SrsTsPidTSDT            = 0x02
    // Reserved
    SrsTsPidReservedStart   = 0x03
    SrsTsPidReservedEnd     = 0x0f
    // May be assigned as network_PID, Program_map_PID, elementary_PID, or for other purposes
    SrsTsPidAppStart        = 0x10
    SrsTsPidAppEnd          = 0x1ffe
    // null packets (see Table 2-3)
    SrsTsPidNULL    = 0x01FFF
)

/**
* the transport_scrambling_control of ts packet,
* Table 2-4 - Scrambling control values, hls-mpeg-ts-iso13818-1.pdf, page 38
*/
type SrsTsScrambled int
const (
	_ SrsTsScrambled = iota
	// Not scrambled
    SrsTsScrambledDisabled      = 0x00
    // User-defined
    SrsTsScrambledUserDefined1  = 0x01
    // User-defined
    SrsTsScrambledUserDefined2  = 0x02
    // User-defined
    SrsTsScrambledUserDefined3  = 0x03
)

/*
* the adaptation_field_control of ts packet
* Table 2-5 - Adaptation field control values hls-mpeg-ts-iso13818-1.pdf, page 38
*/
type SrsTsAdapationControl int
const (
	_ SrsTsAdapationControl 			= 	iota
	SrsTsAdapationControlReserved 		=	0
	SrsTsAdapationControlPayloadOnly	=	1
	SrsTsAdapationControlFieldOnly		=	2
	SrsTsAdapationControlBoth			=	3  
)

const SRS_TS_SYNC_BYTE 	= 0x47

const TS_PMT_NUMBER  	= 1
const TS_PMT_PID 	 	= 0x1001
const TS_VIDEO_AVC_PID 	= 0x100
const TS_AUDIO_AAC_PID  = 0x101
const TS_AUDIO_MP3_PID  = 0x102

type SrsTsPayload interface {
	Encode(stream *utils.SrsStream) error
	Decode(stream *utils.SrsStream) error
}

type SrsTsHeader struct {
	/*
	The sync_byte is a fixed 8-bit field whose value is '0100 0111' (0x47). Sync_byte emulation in the choice of
	values for other regularly occurring fields, such as PID, should be avoided.
	*/
	syncByte					int8	//8bit 同步字节，固定为0x47
	transportErrorIndicator		int8	//1bit 传输错误指示符，表明在ts头的adapt域后由一个无用字节，通常都为0，这个字节算在adapt域长度内
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
	payloadUnitStartIndicator	int8 //1bit
	/*
	The transport_priority is a 1-bit indicator. When set to '1' it indicates that the associated packet is
	of greater priority than other packets having the same PID which do not have the bit set to '1'. The transport mechanism
	can use this to prioritize its data within an elementary stream. Depending on the application the transport_priority field
	may be coded regardless of the PID or within one PID only. This field may be changed by channel specific encoders or
	decoders
	*/
	transportPriority			int8 //通常为0
	/*
	The PID is a 13-bit field, indicating the type of the data stored in the packet payload. PID value 0x0000 is
	reserved for the Program Association Table (see Table 2-25). PID value 0x0001 is reserved for the Conditional Access
	Table (see Table 2-27). PID values 0x0002 – 0x000F are reserved. PID value 0x1FFF is reserved for null packets (see
	Table 2-3).
	*/
	PID							SrsTsPid //13bit
	/*
	This 2-bit field indicates the scrambling mode of the Transport Stream packet payload.
	The Transport Stream packet header, and the adaptation field when present, shall not be scrambled. In the case of a null
	packet the value of the transport_scrambling_control field shall be set to '00' (see Table 2-4).
	*/
	transportScrambingControl	SrsTsScrambled
	/*
	This 2-bit field indicates whether this Transport Stream packet header is followed by an
	adaptation field and/or payload (see Table 2-5).
	*/
	adatpationFieldControl		SrsTsAdapationControl
	
	/*
	The continuity_counter is a 4-bit field incrementing with each Transport Stream packet with the
	same PID. The continuity_counter wraps around to 0 after its maximum value. The continuity_counter shall not be
	incremented when the adaptation_field_control of the packet equals '00' or '10'.

	The continuity counter may be discontinuous when the discontinuity_indicator is set to '1' (refer to 2.4.3.4). In the case of
	a null packet the value of the continuity_counter is undefined.
	*/
	continuityCounter			int8
}

func NewSrsTsHeader() *SrsTsHeader {
	return &SrsTsHeader{
		syncByte:0x47,
		transportErrorIndicator:0,
		payloadUnitStartIndicator:1,
		transportPriority:0,
		PID:SrsTsPidPAT,
		transportScrambingControl:SrsTsScrambledDisabled,
		adatpationFieldControl:SrsTsAdapationControlPayloadOnly,
		continuityCounter:0,
	}
}

func (this *SrsTsHeader) Encode(stream *utils.SrsStream) {

}

func (this *SrsTsHeader) Decode(stream *utils.SrsStream) error {
	return nil
}

type SrsTsAdapationField struct {
	/*
	The adaptation_field_length is an 8-bit field specifying the number of bytes in the
	adaptation_field immediately following the adaptation_field_length. The value 0 is for inserting a single stuffing byte in
	a Transport Stream packet. When the adaptation_field_control value is '11', the value of the adaptation_field_length shall
	be in the range 0 to 182. When the adaptation_field_control value is '10', the value of the adaptation_field_length shall
	be 183. For Transport Stream packets carrying PES packets, stuffing is needed when there is insufficient PES packet data
	to completely fill the Transport Stream packet payload bytes. Stuffing is accomplished by defining an adaptation field
	longer than the sum of the lengths of the data elements in it, so that the payload bytes remaining after the adaptation field
	exactly accommodates the available PES packet data. The extra space in the adaptation field is filled with stuffing bytes.
	This is the only method of stuffing allowed for Transport Stream packets carrying PES packets. For Transport Stream
	packets carrying PSI, an alternative stuffing method is described in 2.4.4.
	*/
	adaptationFieldLength				uint8
	/*
	This is a 1-bit field which when set to '1' indicates that the discontinuity state is true for the
	current Transport Stream packet. When the discontinuity_indicator is set to '0' or is not present, the discontinuity state is
	false. The discontinuity indicator is used to indicate two types of discontinuities, system time-base discontinuities and
	continuity_counter discontinuities.
	@discontinuity_indicator iso13818-1.pdf page 39
	*/
	discontinuityIndicator				int8
	/*
	The random_access_indicator is a 1-bit field that indicates that the current Transport
	Stream packet, and possibly subsequent Transport Stream packets with the same PID, contain some information to aid
	random access at this point. Specifically, when the bit is set to '1', the next PES packet to start in the payload of Transport
	Stream packets with the current PID shall contain the first byte of a video sequence header if the PES stream type (refer
	to Table 2-29) is 1 or 2, or shall contain the first byte of an audio frame if the PES stream type is 3 or 4. In addition, in
	the case of video, a presentation timestamp shall be present in the PES packet containing the first picture following the
	sequence header. In the case of audio, the presentation timestamp shall be present in the PES packet containing the first
	byte of the audio frame. In the PCR_PID the random_access_indicator may only be set to '1' in Transport Stream packet
	containing the PCR fields.
	当设置为1，如果是视频下一个PES包的第一字节必须是视频sequence header（视频包的第一个包？）；如果是音频，也必须是音频的第一个字节
	（PS:是指payloadUnitStartIndicator需要为1吗？）
	*/
	randomAccessIndicator				int8
	/*
	The elementary_stream_priority_indicator is a 1-bit field. It indicates, among
	packets with the same PID, the priority of the elementary stream data carried within the payload of this Transport Stream
	packet. A '1' indicates that the payload has a higher priority than the payloads of other Transport Stream packets. In the
	case of video, this field may be set to '1' only if the payload contains one or more bytes from an intra-coded slice. A
	value of '0' indicates that the payload has the same priority as all other packets which do not have this bit set to '1'.
	*/
	elementaryStreamPriorityIndicator	int8  //一般为0
	/*
	The PCR_flag is a 1-bit flag. A value of '1' indicates that the adaptation_field contains a PCR field coded in
	two parts. A value of '0' indicates that the adaptation field does not contain any PCR field.
	*/
	PCRFlag								int8
	/*
	The OPCR_flag is a 1-bit flag. A value of '1' indicates that the adaptation_field contains an OPCR field
	coded in two parts. A value of '0' indicates that the adaptation field does not contain any OPCR field.
	*/
	OPCRFlag							int8
	/*
	The splicing_point_flag is a 1-bit flag. When set to '1', it indicates that a splice_countdown field
	shall be present in the associated adaptation field, specifying the occurrence of a splicing point. A value of '0' indicates
	that a splice_countdown field is not present in the adaptation field.
	@iso13818-1.pdf, page 41
	*/
	splicingPointFlag					int8
	/*
	The transport_private_data_flag is a 1-bit flag. A value of '1' indicates that the
	adaptation field contains one or more private_data bytes. A value of '0' indicates the adaptation field does not contain any
	private_data bytes.
	*/
	transportPrivateDataFlag			int8	
	/*
	The adaptation_field_extension_flag is a 1-bit field which when set to '1' indicates
	the presence of an adaptation field extension. A value of '0' indicates that an adaptation field extension is not present in
	the adaptation field.
	*/
	adaptationFieldExtensionFlag		int8
	/*
	The program_clock_reference (PCR) is a
	42-bit field coded in two parts. The first part, program_clock_reference_base, is a 33-bit field whose value is given by
	PCR_base(i), as given in equation 2-2. The second part, program_clock_reference_extension, is a 9-bit field whose value
	is given by PCR_ext(i), as given in equation 2-3. The PCR indicates the intended time of arrival of the byte containing
	the last bit of the program_clock_reference_base at the input of the system target decoder
	*/
	programClockReferenceBase			int64

	const1Value0						int8	//6bits
	programClockReferenceExtension		int8
	/*
	The optional original
	program reference (OPCR) is a 42-bit field coded in two parts. These two parts, the base and the extension, are coded
	identically to the two corresponding parts of the PCR field. The presence of the OPCR is indicated by the OPCR_flag.
	The OPCR field shall be coded only in Transport Stream packets in which the PCR field is present. OPCRs are permitted
	in both single program and multiple program Transport Streams.
	OPCR(i) = OPCR_base(i) × 300 + OPCR_ext(i)
	OPCR_base(i) = ((system_clock_ frequency × t(i)) DIV 300) % 233
	OPCR_ext(i) = ((system_clock_ frequency × t(i)) DIV 1) % 300
	*/
	originalProgramClockReferenceBase		int64	//33bits
	const1Value1							int8 	//6bits set to 1
	originalProgramClockReferenceExtension	int16	//9bits
	/*
	The splice_countdown is an 8-bit field, representing a value which may be positive or negative. A
	positive value specifies the remaining number of Transport Stream packets, of the same PID, following the associated
	Transport Stream packet until a splicing point is reached. Duplicate Transport Stream packets and Transport Stream
	packets which only contain adaptation fields are excluded. The splicing point is located immediately after the last byte of
	the Transport Stream packet in which the associated splice_countdown field reaches zero. In the Transport Stream packet
	where the splice_countdown reaches zero, the last data byte of the Transport Stream packet payload shall be the last byte
	of a coded audio frame or a coded picture. In the case of video, the corresponding access unit may or may not be
	terminated by a sequence_end_code. Transport Stream packets with the same PID, which follow, may contain data from
	a different elementary stream of the same type.

	The payload of the next Transport Stream packet of the same PID (duplicate packets and packets without payload being
	excluded) shall commence with the first byte of a PES packet. In the case of audio, the PES packet payload shall
	commence with an access point. In the case of video, the PES packet payload shall commence with an access point, or
	with a sequence_end_code, followed by an access point. Thus, the previous coded audio frame or coded picture aligns
	with the packet boundary, or is padded to make this so. Subsequent to the splicing point, the countdown field may also
	be present. When the splice_countdown is a negative number whose value is minus n (−n), it indicates that the associated
	Transport Stream packet is the n-th packet following the splicing point (duplicate packets and packets without payload
	being excluded).

	For the purposes of this subclause, an access point is defined as follows:
	• Video – The first byte of a video_sequence_header.
	• Audio – The first byte of an audio frame.

	@iso13818-1.pdf, page 42
	*/
	spliceDown								int8	//8bits
	/*
	The transport_private_data_length is an 8-bit field specifying the number of
	private_data bytes immediately following the transport private_data_length field. The number of private_data bytes shall
	not be such that private data extends beyond the adaptation field.
	*/
	transportPrivateDataLength				uint8	//8bits
	/*
	The private_data_byte is an 8-bit field that shall not be specified by ITU-T | ISO/IEC.
	*/
	privateData								[]byte
	/*
	The adaptation_field_extension_length is an 8-bit field. It indicates the number of
	bytes of the extended adaptation field data immediately following this field, including reserved bytes if present
	*/
	adaptationFieldExtensionLength			uint8
	/*
	This is a 1-bit field which when set to '1' indicates the presence of the ltw_offset field.
	*/
	ltwFlag									int8
	/*
	This is a 1-bit field which when set to '1' indicates the presence of the piecewise_rate field
	*/
	piecewiseRateFlag						int8
	/*
	This is a 1-bit flag which when set to '1' indicates that the splice_type and DTS_next_AU fields
	are present. A value of '0' indicates that neither splice_type nor DTS_next_AU fields are present. This field shall not be
	set to '1' in Transport Stream packets in which the splicing_point_flag is not set to '1'. Once it is set to '1' in a Transport
	Stream packet in which the splice_countdown is positive, it shall be set to '1' in all the subsequent Transport Stream
	packets of the same PID that have the splicing_point_flag set to '1', until the packet in which the splice_countdown
	reaches zero (including this packet). When this flag is set, if the elementary stream carried in this PID is an audio stream,
	the splice_type field shall be set to '0000'. If the elementary stream carried in this PID is a video stream, it shall fulfil the
	constraints indicated by the splice_type value.
	*/
	seamlessSpliceFlag						int8
	/*
	This is a 1-bit field which when set to '1' indicates that the value of the
	ltw_offset shall be valid. A value of '0' indicates that the value in the ltw_offset field is undefined.
	*/
	ltwValidFlag							int8
	/*
	(legal time window offset) – This is a 15-bit field, the value of which is defined only if the ltw_valid flag has
	a value of '1'. When defined, the legal time window offset is in units of (300/fs) seconds, where fs is the system clock
	frequency of the program that this PID belongs to, and fulfils:
	offset = t1(i) – t(i)
	ltw_offset = offset//1
	where i is the index of the first byte of this Transport Stream packet, offset is the value encoded in this field, t(i) is the
	arrival time of byte i in the T-STD, and t1(i) is the upper bound in time of a time interval called the Legal Time Window
	which is associated with this Transport Stream packet.
	The Legal Time Window has the property that if this Transport Stream is delivered to a T-STD starting at time t1(i),
	i.e. at the end of its Legal Time Window, and all other Transport Stream packets of the same program are delivered at the
	end of their Legal Time Windows, then
	• For video – The MBn buffer for this PID in the T-STD shall contain less than 184 bytes of elementary
	stream data at the time the first byte of the payload of this Transport Stream packet enters it, and no buffer
	violations in the T-STD shall occur.
	• For audio – The Bn buffer for this PID in the T-STD shall contain less than BSdec + 1 bytes of elementary
	stream data at the time the first byte of this Transport Stream packet enters it, and no buffer violations in
	the T-STD shall occur.
	Depending on factors including the size of the buffer MBn and the rate of data transfer between MBn and EBn, it is
	possible to determine another time t0(i), such that if this packet is delivered anywhere in the interval [t0(i), t1(i)], no
	T-STD buffer violations will occur. This time interval is called the Legal Time Window. The value of t0 is not defined in
	this Recommendation | International Standard.
	The information in this field is intended for devices such as remultiplexers which may need this information in order to
	reconstruct the state of the buffers MBn.
	*/
	ltwOffset								int16 //15 bits
	/*
	The meaning of this 22-bit field is only defined when both the ltw_flag and the ltw_valid_flag are set
	to '1'. When defined, it is a positive integer specifying a hypothetical bitrate R which is used to define the end times of
	the Legal Time Windows of Transport Stream packets of the same PID that follow this packet but do not include the
	legal_time_window_offset field.
	Assume that the first byte of this Transport Stream packet and the N following Transport Stream packets of the same PID
	have indices Ai, Ai + 1, ..., Ai + N, respectively, and that the N latter packets do not have a value encoded in the field
	legal_time_window_offset. Then the values t1(Ai + j) shall be determined by:
	t1
	(Ai + j) = t1(Ai) + j × 188 × 8-bits/byte /R
	where j goes from 1 to N.
	All packets between this packet and the next packet of the same PID to include a legal_time_window_offset field shall
	be treated as if they had the value:
	offset = t1(Ai) – t(Ai)
	corresponding to the value t1(.) as computed by the formula above encoded in the legal_time_window_offset field. t(j) is
	the arrival time of byte j in the T-STD.
	The meaning of this field is not defined when it is present in a Transport Stream packet with no
	legal_time_window_offset field.
	*/
	piecewiseRate							int32 //22 bits
	/*
	This is a 4-bit field. From the first occurrence of this field onwards, it shall have the same value in all the
	subsequent Transport Stream packets of the same PID in which it is present, until the packet in which the
	splice_countdown reaches zero (including this packet). If the elementary stream carried in that PID is an audio stream,
	this field shall have the value '0000'. If the elementary stream carried in that PID is a video stream, this field indicates the
	conditions that shall be respected by this elementary stream for splicing purposes. These conditions are defined as a
	function of profile, level and splice_type in Table 2-7 through Table 2-16
	@iso13818-1.pdf, page 43
	*/
	spliceType								int8 //4bits	
	/*
	(decoding time stamp next access unit) – This is a 33-bit field, coded in three parts. In the case of
	continuous and periodic decoding through this splicing point it indicates the decoding time of the first access unit
	following the splicing point. This decoding time is expressed in the time base which is valid in the Transport Stream
	packet in which the splice_countdown reaches zero. From the first occurrence of this field onwards, it shall have the
	same value in all the subsequent Transport Stream packets of the same PID in which it is present, until the packet in
	which the splice_countdown reaches zero (including this packet).
	*/
	DTSNextAU0								int8 //3bits
	markerBit0								int8 //1bit
	DTSNextAU1								int16 //15bits
	markerBit1								int8 //1bit
	DTSNextAU2								int16 //15bits
	markerBit2								int8 //1bit	
	/*
	This is a fixed 8-bit value equal to '1111 1111' that can be inserted by the encoder. It is discarded by the decoder.
	*/
	staffingByte							[]byte
}

func (this *SrsTsAdapationField) Decode(stream *utils.SrsStream) error {
	return nil
}

func (this *SrsTsAdapationField) Encode(stream *utils.SrsStream) {

}

/**
* 2.4.4.4 Table_id assignments, hls-mpeg-ts-iso13818-1.pdf, page 62
* The table_id field identifies the contents of a Transport Stream PSI section as shown in Table 2-26.
*/
type SrsTsPsiTableId int
const (
	_ SrsTsPsiId = 	iota
	// program_association_section
    SrsTsPsiTableIdPas               = 0x00
    // conditional_access_section (CA_section)
    SrsTsPsiTableIdCas               = 0x01
    // TS_program_map_section
    SrsTsPsiTableIdPms               = 0x02
    // TS_description_section
    SrsTsPsiTableIdDs                = 0x03
    // ISO_IEC_14496_scene_description_section
    SrsTsPsiTableIdSds               = 0x04
    // ISO_IEC_14496_object_descriptor_section
    SrsTsPsiTableIdOds               = 0x05
    // ITU-T Rec. H.222.0 | ISO/IEC 13818-1 reserved
    SrsTsPsiIdTableIso138181Start    = 0x06
    SrsTsPsiIdTableIso138181End      = 0x37
    // Defined in ISO/IEC 13818-6
    SrsTsPsiIdTableIso138186Start    = 0x38
    SrsTsPsiIdTableIso138186End      = 0x3F
    // User private
    SrsTsPsiTableIdUserStart         = 0x40
    SrsTsPsiTableIdUserEnd           = 0xFE
    // forbidden
    SrsTsPsiTableIdForbidden         = 0xFF
)

type SrsTsPayloadPSI struct {
	/**
    * This is an 8-bit field whose value shall be the number of bytes, immediately following the pointer_field
    * until the first byte of the first section that is present in the payload of the Transport Stream packet (so a value of 0x00 in
    * the pointer_field indicates that the section starts immediately after the pointer_field). When at least one section begins in
    * a given Transport Stream packet, then the payload_unit_start_indicator (refer to 2.4.3.2) shall be set to 1 and the first
    * byte of the payload of that Transport Stream packet shall contain the pointer. When no section begins in a given
    * Transport Stream packet, then the payload_unit_start_indicator shall be set to 0 and no pointer shall be sent in the
    * payload of that packet.
    */
	pointerField 				int8
	// 1B
    /**
    * This is an 8-bit field, which shall be set to 0x00 as shown in Table 2-26.
    */
	tableId						SrsTsPsiTableId	//PAT表固定为0x00,PMT表为0x02
	/*
	The section_syntax_indicator is a 1-bit field which shall be set to '1'.
	*/
	sectionSyntaxIndicator		int8	//固定为二进制1
	/**
    * const value, must be '0'
    */
	const0Value					int8	//1bit
	/**
    * reverved value, must be '1'
    */
	const1Value0				int8 	//2bits
	/*
	This is a 12-bit field, the first two bits of which shall be '00'. The remaining 10 bits specify the number
	of bytes of the section, starting immediately following the section_length field, and including the CRC. The value in this
	field shall not exceed 1021 (0x3FD).
	*/
	sectionLength				int16	//12bits 	后面数据的长度
}

func CreatePAT(context *SrsTsContext, pmt_number int16, pmt_pid int16) *SrsTsPacket {
	pkt := NewSrsTsPacket()

	pkt.tsHeader.syncByte = SRS_TS_SYNC_BYTE
	pkt.tsHeader.transportErrorIndicator = 0
	pkt.tsHeader.payloadUnitStartIndicator = 1
	pkt.tsHeader.transportPriority = 0
	pkt.tsHeader.PID = SrsTsPidPAT
	pkt.tsHeader.transportScrambingControl = SrsTsScrambledDisabled
	pkt.tsHeader.adatpationFieldControl = SrsTsAdapationControlPayloadOnly
	pkt.tsHeader.continuityCounter = 0

	pkt.payload = NewSrsTsPayloadPAT()
	var pat *SrsTsPayloadPAT = pkt.payload.(*SrsTsPayloadPAT)
	pat.psiHeader.pointerField	= 0
	pat.psiHeader.tableId = SrsTsPsiIdPas
	pat.psiHeader.sectionSyntaxIndicator = 1
	pat.psiHeader.const0Value = 0
	pat.psiHeader.const1Value0 = 0xff
	pat.psiHeader.sectionLength = 0	//calc in size

	pat.transportStreamId 	= 1
	pat.const1Value0		= 0xff
	pat.versionNumber 		= 0
	pat.currentNextIndicator= 1
	pat.sectionNumber		= 0
	pat.lastSectionNumber	= 0
	program := NewSrsTsPayloadPATProgram(pmt_number, pmt_pid)
	pat.programs = append(pat.programs, program)
}

type   struct {
	/*
	Program_number is a 16-bit field. It specifies the program to which the program_map_PID is
	applicable. When set to 0x0000, then the following PID reference shall be the network PID. For all other cases the value
	of this field is user defined. This field shall not take any single value more than once within one version of the Program
	Association Table.
	*/
	programNumber			int16	//节目号为0x0000时表示这是NIT，节目号为0x0001时,表示这是PMT
	const1Value0			int8	//3bit '1'
	/*
	Program_number is a 16-bit field. It specifies the program to which the program_map_PID is
	applicable. When set to 0x0000, then the following PID reference shall be the network PID. For all other cases the value
	of this field is user defined. This field shall not take any single value more than once within one version of the Program
	Association Table.
	*/
	pid						int16	//13bit
}

type SrsTsPayloadPATProgram struct {
	// 4B
    /**
    * Program_number is a 16-bit field. It specifies the program to which the program_map_PID is
    * applicable. When set to 0x0000, then the following PID reference shall be the network PID. For all other cases the value
    * of this field is user defined. This field shall not take any single value more than once within one version of the Program
    * Association Table.
    */
    number 			int16 // 16bits
    /**
    * reverved value, must be '1'
    */
    const1Value		int8 //3bits
    /**
    * program_map_PID/network_PID 13bits
    * network_PID - The network_PID is a 13-bit field, which is used only in conjunction with the value of the
    * program_number set to 0x0000, specifies the PID of the Transport Stream packets which shall contain the Network
    * Information Table. The value of the network_PID field is defined by the user, but shall only take values as specified in
    * Table 2-3. The presence of the network_PID is optional.
    */
    pid				int16 //13bits
}

func NewSrsTsPayloadPATProgram(program_number int16, p int16) *SrsTsPayloadPATProgram {
	return &SrsTsPayloadPATProgram{
		number:program_number,
		const1Value:0x7,
		pid:p,
	}
}

type SrsTsPayloadPAT struct {
	psiHeader				*SrsTsPayloadPSI
	/*
	This is a 16-bit field which serves as a label to identify this Transport Stream from any other
	multiplex within a network. Its value is defined by the user
	*/
	transportStreamId		int16	//固定为0x0001
	const1Value0			int8	//2bits
	/*
	This 5-bit field is the version number of the whole Program Association Table. The version number
	shall be incremented by 1 modulo 32 whenever the definition of the Program Association Table changes. When the
	current_next_indicator is set to '1', then the version_number shall be that of the currently applicable Program Association
	Table. When the current_next_indicator is set to '0', then the version_number shall be that of the next applicable Program
	Association Table.
	*/
	versionNumber			int8	//5bits	版本号，固定为二进制00000，如果PAT有变化则版本号加1
	/*
	A 1-bit indicator, which when set to '1' indicates that the Program Association Table sent is
	currently applicable. When the bit is set to '0', it indicates that the table sent is not yet applicable and shall be the next
	table to become valid.
	*/
	currentNextIndicator	int8	//固定为二进制1，表示这个PAT表可以用，如果为0则要等待下一个PAT表
	/*
	This 8-bit field gives the number of this section. The section_number of the first section in the
	Program Association Table shall be 0x00. It shall be incremented by 1 with each additional section in the Program
	Association Table.
	*/
	sectionNumber			int8	//固定为0x00
	/*
	This 8-bit field specifies the number of the last section (that is, the section with the highest
	section_number) of the complete Program Association Table.
	*/
	lastSectionNumber		int8	//固定为0x00
	// multiple 4B program data.
	programs				[]*SrsTsPayloadPATProgram
	// 4B
    /**
    * This is a 32-bit field that contains the CRC value that gives a zero output of the registers in the decoder
    * defined in Annex A after processing the entire section.
    * @remark crc32(bytes without pointer field, before crc32 field)
	*/
	crc32					int32
}

func NewSrsTsPayloadPAT() *SrsTsPayloadPAT {
	return &SrsTsPayloadPAT{}
}

func (this *SrsTsPayloadPAT) Encode(stream *utils.SrsStream) {

}

func (this *SrsTsPayloadPAT) Decode(stream *utils.SrsStream) error {
	return nil
}

type SrsTsPayloadPMTESInfo struct {
	/*
	This is an 8-bit field specifying the type of program element carried within the packets with the PID
	whose value is specified by the elementary_PID. The values of stream_type are specified in Table 2-29.
	*/
	streamType			SrsTsStream	//流类型，标志是Video还是Audio还是其他数据，h.264编码对应0x1b，aac编码对应0x0f，mp3编码对应0x03

	const1Value0		int8	//3bits
	/*
	This is a 13-bit field specifying the PID of the Transport Stream packets which carry the associated
	program element
	*/
	elemenaryPID		int16	//13bits
	const1Value1		int8	//4bits
	/*
	This is a 12-bit field, the first two bits of which shall be '00'. The remaining 10 bits specify the number
	of bytes of the descriptors of the associated program element immediately following the ES_info_length field.
	*/
	ESInfoLength		int16	//12bits
	ESInfo				[]byte
}

type SrsTsPayloadPMT struct {
	psiHeader				*SrsTsPayloadPSI
	/*
	program_number is a 16-bit field. It specifies the program to which the program_map_PID is
	applicable. One program definition shall be carried within only one TS_program_map_section. This implies that a
	program definition is never longer than 1016 (0x3F8). See Informative Annex C for ways to deal with the cases when
	that length is not sufficient. The program_number may be used as a designation for a broadcast channel, for example. By
	describing the different program elements belonging to a program, data from different sources (e.g. sequential events)
	can be concatenated together to form a continuous set of streams using a program_number. For examples of applications
	refer to Annex C.
	*/
	programNumber			int16	//频道号码，表示当前的PMT关联到的频道，取值0x0001

	// 1B
    /**
    * reverved value, must be '1'
    */
	const1Value0			int8 	//2bits
	/*
	This 5-bit field is the version number of the TS_program_map_section. The version number shall be
	incremented by 1 modulo 32 when a change in the information carried within the section occurs. Version number refers
	to the definition of a single program, and therefore to a single section. When the current_next_indicator is set to '1', then
	the version_number shall be that of the currently applicable TS_program_map_section. When the current_next_indicator
	is set to '0', then the version_number shall be that of the next applicable TS_program_map_section.
	*/
	versionNumber			int8	//5bits 版本号，固定为00000，如果PAT有变化则版本号加1
	/*
	A 1-bit field, which when set to '1' indicates that the TS_program_map_section sent is
	currently applicable. When the bit is set to '0', it indicates that the TS_program_map_section sent is not yet applicable
	and shall be the next TS_program_map_section to become valid.
	*/
	currentNextIndicator	int8	//1bit	固定为1就好，没那么复杂
	/*
	The value of this 8-bit field shall be 0x00
	*/
	sectionNumber			int8
	/*
	The value of this 8-bit field shall be 0x00.
	*/
	lastSectionNumber		int8

	/*
	This is a 13-bit field indicating the PID of the Transport Stream packets which shall contain the PCR fields
	valid for the program specified by program_number. If no PCR is associated with a program definition for private
	streams, then this field shall take the value of 0x1FFF. Refer to the semantic definition of PCR in 2.4.3.5 and Table 2-3
	for restrictions on the choice of PCR_PID value
	*/
	PCR_PID					int16
	// 2B
	const1Value2			int8	//4bits
	/*
	This is a 12-bit field, the first two bits of which shall be '00'. The remaining 10 bits specify the
	number of bytes of the descriptors immediately following the program_info_length field.
	*/
	programInfoLength		int16
	programDescriptor		[]byte	//the len is programInfoLength

	infoes					[]*SrsTsPayloadPMTESInfo
}

//see iso-13818.pdf, page 49
type SrsTsPayloadPES struct {
	/*
	The packet_start_code_prefix is a 24-bit code. Together with the stream_id that follows, it
	constitutes a packet start code that identifies the beginning of a packet. The packet_start_code_prefix is the bit string
	'0000 0000 0000 0000 0000 0001' (0x000001 in hexadecimal).
	*/
	packetStartCodePrefix	int32		//24bit 0x000001
	/*
	In Program Streams, the stream_id specifies the type and number of the elementary stream as defined by the
	stream_id Table 2-18. In Transport Streams, the stream_id may be set to any valid value which correctly describes the
	elementary stream type as defined in Table 2-18. In Transport Streams, the elementary stream type is specified in the
	Program Specific Information as specified in 2.4.4.
	*/
	streamId				int8		//音频取值（0xc0-0xdf），通常为0xc0 视频取值（0xe0-0xef），通常为0xe0
	/*
	The PES_packet_length is a 16-bit field indicating the total number of bytes in the
	program_stream_directory immediately following this field (refer to Table 2-18).
	*/
	PESPacketLength			uint16		

	const2Bits				int8	//2bits '10'
	/*
	The 2-bit PES_scrambling_control field indicates the scrambling mode of the PES packet
	payload. When scrambling is performed at the PES level, the PES packet header, including the optional fields when
	present, shall not be scrambled (see Table 2-19)
	*/
	PESScramblingControl	int8	//2bit
	/**
    * This is a 1-bit field indicating the priority of the payload in this PES packet. A '1' indicates a higher
    * priority of the payload of the PES packet payload than a PES packet payload with this field set to '0'. A multiplexor can
    * use the PES_priority bit to prioritize its data within an elementary stream. This field shall not be changed by the transport
    * mechanism.
    */
    PESPriority				int8 	//1bit 一般为0
    /**
    * This is a 1-bit flag. When set to a value of '1' it indicates that the PES packet header is
    * immediately followed by the video start code or audio syncword indicated in the data_stream_alignment_descriptor
    * in 2.6.10 if this descriptor is present. If set to a value of '1' and the descriptor is not present, alignment as indicated in
    * alignment_type '01' in Table 2-47 and Table 2-48 is required. When set to a value of '0' it is not defined whether any such
    * alignment occurs or not.
    */
    dataAlignmentIndicator	int8 	//1bit 一般为0
    /**
    * This is a 1-bit field. When set to '1' it indicates that the material of the associated PES packet payload is
    * protected by copyright. When set to '0' it is not defined whether the material is protected by copyright. A copyright
    * descriptor described in 2.6.24 is associated with the elementary stream which contains this PES packet and the copyright
    * flag is set to '1' if the descriptor applies to the material contained in this PES packet
    */
    copyright				int8 	//1bit 一般为0
    /**
    * This is a 1-bit field. When set to '1' the contents of the associated PES packet payload is an original.
    * When set to '0' it indicates that the contents of the associated PES packet payload is a copy.
    */
    originalOrCopy			int8 	//1bit 一般为1

    // 1B
    /**
    * This is a 2-bit field. When the PTS_DTS_flags field is set to '10', the PTS fields shall be present in
    * the PES packet header. When the PTS_DTS_flags field is set to '11', both the PTS fields and DTS fields shall be present
    * in the PES packet header. When the PTS_DTS_flags field is set to '00' no PTS or DTS fields shall be present in the PES
    * packet header. The value '01' is forbidden.
    */
	PTSDTSflags				int8 	//2bits
    /**
    * A 1-bit flag, which when set to '1' indicates that ESCR base and extension fields are present in the PES
    * packet header. When set to '0' it indicates that no ESCR fields are present.
    */
    ESCRFlag				int8 	//1bit	一般为 0吧
    /**
    * A 1-bit flag, which when set to '1' indicates that the ES_rate field is present in the PES packet header.
    * When set to '0' it indicates that no ES_rate field is present.
    */
	ESRateFlag				int8 	//1bit	一般为0
    /**
    * A 1-bit flag, which when set to '1' it indicates the presence of an 8-bit trick mode field. When
    * set to '0' it indicates that this field is not present.
    */
    DSMTrickModeFlag		int8	//1bit
    /**
    * A 1-bit flag, which when set to '1' indicates the presence of the additional_copy_info field.
    * When set to '0' it indicates that this field is not present.
    */
    additionalCopyInfoFlag	int8 //1bit
    /**
    * A 1-bit flag, which when set to '1' indicates that a CRC field is present in the PES packet. When set to
    * '0' it indicates that this field is not present.
    */
    PESCRCFlag				int8 //1bit
    /**
    * A 1-bit flag, which when set to '1' indicates that an extension field exists in this PES packet
    * header. When set to '0' it indicates that this field is not present.
    */
    PESExtensionFlag		int8 //1bit

    // 1B
    /**
    * An 8-bit field specifying the total number of bytes occupied by the optional fields and any
    * stuffing bytes contained in this PES packet header. The presence of optional fields is indicated in the byte that precedes
    * the PES_header_data_length field.
    */
    PESHeaderDataLength		uint8 //8bits

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
    pts						int64 // 33bits

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
    dts						int64 // 33bits

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
    ESCRBase			int64 //33bits
    ESCRExtension		int16 //9bits

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
    ESRate				int32 //22bits

    // 1B
    /**
    * A 3-bit field that indicates which trick mode is applied to the associated video stream. In cases of
    * other types of elementary streams, the meanings of this field and those defined by the following five bits are undefined.
    * For the definition of trick_mode status, refer to the trick mode section of 2.4.2.3.
    */
    trickModeControl	int8 //3bits
    trickModeValue		int8 //5bits

    // 1B
    // 1bit const '1'
    /**
    * This 7-bit field contains private data relating to copyright information.
    */
    additionalCopyInfo	int8 //7bits

    // 2B
    /**
    * The previous_PES_packet_CRC is a 16-bit field that contains the CRC value that yields
    * a zero output of the 16 registers in the decoder similar to the one defined in Annex A,
    */
    previousPESPacketCRC	int8 //16bits

    // 1B
    /**
    * A 1-bit flag which when set to '1' indicates that the PES packet header contains private data.
    * When set to a value of '0' it indicates that private data is not present in the PES header.
    */
    PESPrivateDataFlag		int8 //1bit
    /**
    * A 1-bit flag which when set to '1' indicates that an ISO/IEC 11172-1 pack header or a
    * Program Stream pack header is stored in this PES packet header. If this field is in a PES packet that is contained in a
    * Program Stream, then this field shall be set to '0'. In a Transport Stream, when set to the value '0' it indicates that no pack
    * header is present in the PES header.
    */
    packHeaderFieldFlag		int8 //1bit
    /**
    * A 1-bit flag which when set to '1' indicates that the
    * program_packet_sequence_counter, MPEG1_MPEG2_identifier, and original_stuff_length fields are present in this
    * PES packet. When set to a value of '0' it indicates that these fields are not present in the PES header.
    */
    programPacketSequenceCounterFlag	int8 //1bit
    /**
    * A 1-bit flag which when set to '1' indicates that the P-STD_buffer_scale and P-STD_buffer_size
    * are present in the PES packet header. When set to a value of '0' it indicates that these fields are not present in the
    * PES header.
    */
    PSTDBufferFlag			int8 //1bit
    /**
    * reverved value, must be '1'
    */
    const1Value0			int8 //3bits
    /**
    * A 1-bit field which when set to '1' indicates the presence of the PES_extension_field_length
    * field and associated fields. When set to a value of '0' this indicates that the PES_extension_field_length field and any
    * associated fields are not present.
    */
    PESExtensionFlag2		int8//1bit

    // 16B
    /**
    * This is a 16-byte field which contains private data. This data, combined with the fields before and
    * after, shall not emulate the packet_start_code_prefix (0x000001).
    */
    PESPrivateData			[]byte//128bits

    // (1+x)B
    /**
    * This is an 8-bit field which indicates the length, in bytes, of the pack_header_field().
    */
    packFieldLength			uint8 //8bits
    packField				[]byte //[pack_field_length] bytes

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
    programPacketSequenceCounter	int8 //7bits
    // 1bit const '1'
    /**
    * A 1-bit flag which when set to '1' indicates that this PES packet carries information from
    * an ISO/IEC 11172-1 stream. When set to '0' it indicates that this PES packet carries information from a Program Stream.
    */
    MPEG1MPEG2Identifier			int8 //1bit
    /**
    * This 6-bit field specifies the number of stuffing bytes used in the original ITU-T
    * Rec. H.222.0 | ISO/IEC 13818-1 PES packet header or in the original ISO/IEC 11172-1 packet header.
    */
    originalStuffLength				int8 //6bits

    // 2B
    // 2bits const '01'
    /**
    * The P-STD_buffer_scale is a 1-bit field, the meaning of which is only defined if this PES packet
    * is contained in a Program Stream. It indicates the scaling factor used to interpret the subsequent P-STD_buffer_size field.
    * If the preceding stream_id indicates an audio stream, P-STD_buffer_scale shall have the value '0'. If the preceding
    * stream_id indicates a video stream, P-STD_buffer_scale shall have the value '1'. For all other stream types, the value
    * may be either '1' or '0'.
    */
    PSTDBufferScale				int8 //1bit
    /**
    * The P-STD_buffer_size is a 13-bit unsigned integer, the meaning of which is only defined if this
    * PES packet is contained in a Program Stream. It defines the size of the input buffer, BS n , in the P-STD. If
    * P-STD_buffer_scale has the value '0', then the P-STD_buffer_size measures the buffer size in units of 128 bytes. If
    * P-STD_buffer_scale has the value '1', then the P-STD_buffer_size measures the buffer size in units of 1024 bytes.
    */
    PSTDBufferSize				int16 //13bits

    // (1+x)B
    // 1bit const '1'
    /**
    * This is a 7-bit field which specifies the length, in bytes, of the data following this field in
    * the PES extension field up to and including any reserved bytes.
    */
	PESExtensionFieldLength		uint8 //7bits
	PESExtensionField			[]byte

    // NB
    /**
    * This is a fixed 8-bit value equal to '1111 1111' that can be inserted by the encoder, for example to meet
    * the requirements of the channel. It is discarded by the decoder. No more than 32 stuffing bytes shall be present in one
    * PES packet header.
	*/
	stuffingBytes					[]byte

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
	
    dataBytes						[]byte

    // NB
    /**
    * This is a fixed 8-bit value equal to '1111 1111'. It is discarded by the decoder.
    */
    paddings						[]byte
}

type SrsTsPacket struct {
	tsHeader			*SrsTsHeader
	adaptationField		*SrsTsAdapationField
	payload				SrsTsPayload
}

func NewSrsTsPacket() *SrsTsPacket {
	return &SrsTsPacket{
		tsHeader:NewSrsTsHeader()
	}
}

func (this *SrsTsMessage) Decode(stream *utils.SrsStream) error {
	return nil
}

func (this *SrsTsMessage) Encode(stream *utils.SrsStream) {

}