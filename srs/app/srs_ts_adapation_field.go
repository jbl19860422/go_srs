package app

import "go_srs/srs/utils"

type SrsTsAdapationField struct {
	packet 				*SrsTsPacket
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
	adaptationFieldLength uint8
	/*
		This is a 1-bit field which when set to '1' indicates that the discontinuity state is true for the
		current Transport Stream packet. When the discontinuity_indicator is set to '0' or is not present, the discontinuity state is
		false. The discontinuity indicator is used to indicate two types of discontinuities, system time-base discontinuities and
		continuity_counter discontinuities.
		@discontinuity_indicator iso13818-1.pdf page 39
	*/
	discontinuityIndicator int8
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
	randomAccessIndicator int8
	/*
		The elementary_stream_priority_indicator is a 1-bit field. It indicates, among
		packets with the same PID, the priority of the elementary stream data carried within the payload of this Transport Stream
		packet. A '1' indicates that the payload has a higher priority than the payloads of other Transport Stream packets. In the
		case of video, this field may be set to '1' only if the payload contains one or more bytes from an intra-coded slice. A
		value of '0' indicates that the payload has the same priority as all other packets which do not have this bit set to '1'.
	*/
	elementaryStreamPriorityIndicator int8 //一般为0
	/*
		The PCR_flag is a 1-bit flag. A value of '1' indicates that the adaptation_field contains a PCR field coded in
		two parts. A value of '0' indicates that the adaptation field does not contain any PCR field.
	*/
	PCRFlag int8
	/*
		The OPCR_flag is a 1-bit flag. A value of '1' indicates that the adaptation_field contains an OPCR field
		coded in two parts. A value of '0' indicates that the adaptation field does not contain any OPCR field.
	*/
	OPCRFlag int8
	/*
		The splicing_point_flag is a 1-bit flag. When set to '1', it indicates that a splice_countdown field
		shall be present in the associated adaptation field, specifying the occurrence of a splicing point. A value of '0' indicates
		that a splice_countdown field is not present in the adaptation field.
		@iso13818-1.pdf, page 41
	*/
	splicingPointFlag int8
	/*
		The transport_private_data_flag is a 1-bit flag. A value of '1' indicates that the
		adaptation field contains one or more private_data bytes. A value of '0' indicates the adaptation field does not contain any
		private_data bytes.
	*/
	transportPrivateDataFlag int8
	/*
		The adaptation_field_extension_flag is a 1-bit field which when set to '1' indicates
		the presence of an adaptation field extension. A value of '0' indicates that an adaptation field extension is not present in
		the adaptation field.
	*/
	adaptationFieldExtensionFlag int8
	/*
		The program_clock_reference (PCR) is a
		42-bit field coded in two parts. The first part, program_clock_reference_base, is a 33-bit field whose value is given by
		PCR_base(i), as given in equation 2-2. The second part, program_clock_reference_extension, is a 9-bit field whose value
		is given by PCR_ext(i), as given in equation 2-3. The PCR indicates the intended time of arrival of the byte containing
		the last bit of the program_clock_reference_base at the input of the system target decoder
	*/
	programClockReferenceBase int64

	const1Value0                   int8 //6bits
	programClockReferenceExtension int8
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
	originalProgramClockReferenceBase      int64 //33bits
	const1Value1                           int8  //6bits set to 1
	originalProgramClockReferenceExtension int16 //9bits
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
	spliceDown int8 //8bits
	/*
		The transport_private_data_length is an 8-bit field specifying the number of
		private_data bytes immediately following the transport private_data_length field. The number of private_data bytes shall
		not be such that private data extends beyond the adaptation field.
	*/
	transportPrivateDataLength uint8 //8bits
	/*
		The private_data_byte is an 8-bit field that shall not be specified by ITU-T | ISO/IEC.
	*/
	privateData []byte
	/*
		The adaptation_field_extension_length is an 8-bit field. It indicates the number of
		bytes of the extended adaptation field data immediately following this field, including reserved bytes if present
	*/
	adaptationFieldExtensionLength uint8
	/*
		This is a 1-bit field which when set to '1' indicates the presence of the ltw_offset field.
	*/
	ltwFlag int8
	/*
		This is a 1-bit field which when set to '1' indicates the presence of the piecewise_rate field
	*/
	piecewiseRateFlag int8
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
	seamlessSpliceFlag int8
	/*
		This is a 1-bit field which when set to '1' indicates that the value of the
		ltw_offset shall be valid. A value of '0' indicates that the value in the ltw_offset field is undefined.
	*/
	ltwValidFlag int8
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
	ltwOffset int16 //15 bits
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
	piecewiseRate int32 //22 bits
	/*
		This is a 4-bit field. From the first occurrence of this field onwards, it shall have the same value in all the
		subsequent Transport Stream packets of the same PID in which it is present, until the packet in which the
		splice_countdown reaches zero (including this packet). If the elementary stream carried in that PID is an audio stream,
		this field shall have the value '0000'. If the elementary stream carried in that PID is a video stream, this field indicates the
		conditions that shall be respected by this elementary stream for splicing purposes. These conditions are defined as a
		function of profile, level and splice_type in Table 2-7 through Table 2-16
		@iso13818-1.pdf, page 43
	*/
	spliceType int8 //4bits
	/*
		(decoding time stamp next access unit) – This is a 33-bit field, coded in three parts. In the case of
		continuous and periodic decoding through this splicing point it indicates the decoding time of the first access unit
		following the splicing point. This decoding time is expressed in the time base which is valid in the Transport Stream
		packet in which the splice_countdown reaches zero. From the first occurrence of this field onwards, it shall have the
		same value in all the subsequent Transport Stream packets of the same PID in which it is present, until the packet in
		which the splice_countdown reaches zero (including this packet).
	*/
	DTSNextAU0 int8  //3bits
	markerBit0 int8  //1bit
	DTSNextAU1 int16 //15bits
	markerBit1 int8  //1bit
	DTSNextAU2 int16 //15bits
	markerBit2 int8  //1bit
	/*
		This is a fixed 8-bit value equal to '1111 1111' that can be inserted by the encoder. It is discarded by the decoder.
	*/
	staffingByte []byte
}

func NewSrsTsAdaptationField(p *SrsTsPacket) *SrsTsAdapationField {
	return &SrsTsAdapationField{
		packet:p,
	}
}

func (this *SrsTsAdapationField) Decode(stream *utils.SrsStream) error {
	return nil
}

func (this *SrsTsAdapationField) Encode(stream *utils.SrsStream) {

}

func (this *SrsTsAdapationField) Size() uint32 {
	return 0
}
