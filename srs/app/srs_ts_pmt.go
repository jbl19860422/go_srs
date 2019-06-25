package app

import "go_srs/srs/utils"

type SrsTsPayloadPMTESInfo struct {
	/*
		This is an 8-bit field specifying the type of program element carried within the packets with the PID
		whose value is specified by the elementary_PID. The values of stream_type are specified in Table 2-29.
	*/
	streamType SrsTsStream //流类型，标志是Video还是Audio还是其他数据，h.264编码对应0x1b，aac编码对应0x0f，mp3编码对应0x03

	const1Value0 int8 //3bits
	/*
		This is a 13-bit field specifying the PID of the Transport Stream packets which carry the associated
		program element
	*/
	elemenaryPID int16 //13bits
	const1Value1 int8  //4bits
	/*
		This is a 12-bit field, the first two bits of which shall be '00'. The remaining 10 bits specify the number
		of bytes of the descriptors of the associated program element immediately following the ES_info_length field.
	*/
	ESInfoLength int16 //12bits
	ESInfo       []byte
}

type SrsTsPayloadPMT struct {
	psiHeader *SrsTsPayloadPSI
	/*
		program_number is a 16-bit field. It specifies the program to which the program_map_PID is
		applicable. One program definition shall be carried within only one TS_program_map_section. This implies that a
		program definition is never longer than 1016 (0x3F8). See Informative Annex C for ways to deal with the cases when
		that length is not sufficient. The program_number may be used as a designation for a broadcast channel, for example. By
		describing the different program elements belonging to a program, data from different sources (e.g. sequential events)
		can be concatenated together to form a continuous set of streams using a program_number. For examples of applications
		refer to Annex C.
	*/
	programNumber int16 //频道号码，表示当前的PMT关联到的频道，取值0x0001

	// 1B
	/**
	 * reverved value, must be '1'
	 */
	const1Value0 int8 //2bits
	/*
		This 5-bit field is the version number of the TS_program_map_section. The version number shall be
		incremented by 1 modulo 32 when a change in the information carried within the section occurs. Version number refers
		to the definition of a single program, and therefore to a single section. When the current_next_indicator is set to '1', then
		the version_number shall be that of the currently applicable TS_program_map_section. When the current_next_indicator
		is set to '0', then the version_number shall be that of the next applicable TS_program_map_section.
	*/
	versionNumber int8 //5bits 版本号，固定为00000，如果PAT有变化则版本号加1
	/*
		A 1-bit field, which when set to '1' indicates that the TS_program_map_section sent is
		currently applicable. When the bit is set to '0', it indicates that the TS_program_map_section sent is not yet applicable
		and shall be the next TS_program_map_section to become valid.
	*/
	currentNextIndicator int8 //1bit	固定为1就好，没那么复杂
	/*
		The value of this 8-bit field shall be 0x00
	*/
	sectionNumber int8
	/*
		The value of this 8-bit field shall be 0x00.
	*/
	lastSectionNumber int8

	/*
		This is a 13-bit field indicating the PID of the Transport Stream packets which shall contain the PCR fields
		valid for the program specified by program_number. If no PCR is associated with a program definition for private
		streams, then this field shall take the value of 0x1FFF. Refer to the semantic definition of PCR in 2.4.3.5 and Table 2-3
		for restrictions on the choice of PCR_PID value
	*/
	PCR_PID int16
	// 2B
	const1Value2 int8 //4bits
	/*
		This is a 12-bit field, the first two bits of which shall be '00'. The remaining 10 bits specify the
		number of bytes of the descriptors immediately following the program_info_length field.
	*/
	programInfoLength int16
	programDescriptor []byte //the len is programInfoLength

	infoes []*SrsTsPayloadPMTESInfo
}

func NewSrsTsPayloadPMT() *SrsTsPayloadPMT {
	return &SrsTsPayloadPMT{
		psiHeader: NewSrsTsPayloadPSI(),
	}
}

func (this *SrsTsPayloadPMT) Encode(stream *utils.SrsStream) {

}

func (this *SrsTsPayloadPMT) Decode(stream *utils.SrsStream) error {
	return nil
}
