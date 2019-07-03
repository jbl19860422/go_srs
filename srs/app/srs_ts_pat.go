package app
//discuss @iso13818-1.pdf, page 61
import (
	"encoding/binary"
	"go_srs/srs/utils"
	// "hash/crc32"
)

type SrsTsPayloadPATProgram struct {
	// 4B
	/**
	 * Program_number is a 16-bit field. It specifies the program to which the program_map_PID is
	 * applicable. When set to 0x0000, then the following PID reference shall be the network PID. For all other cases the value
	 * of this field is user defined. This field shall not take any single value more than once within one version of the Program
	 * Association Table.
	 */
	number int16 // 16bits
	/**
	 * reverved value, must be '1'
	 */
	const1Value int8 //3bits
	/**
	 * program_map_PID/network_PID 13bits
	 * network_PID - The network_PID is a 13-bit field, which is used only in conjunction with the value of the
	 * program_number set to 0x0000, specifies the PID of the Transport Stream packets which shall contain the Network
	 * Information Table. The value of the network_PID field is defined by the user, but shall only take values as specified in
	 * Table 2-3. The presence of the network_PID is optional.
	 */
	pid int16 //13bits
}

func NewSrsTsPayloadPATProgram(program_number int16, p int16) *SrsTsPayloadPATProgram {
	return &SrsTsPayloadPATProgram{
		number:      program_number,
		const1Value: 0x7,
		pid:         p,
	}
}

func (this *SrsTsPayloadPATProgram) Encode(stream *utils.SrsStream) {
	var tmpv int32 = int32(this.pid) & 0x1FFF
    tmpv |= int32((uint32(this.number) << 16) & 0xFFFF0000)
	tmpv |= (int32(this.const1Value) << 13) & 0xE000
	stream.WriteInt32(tmpv, binary.BigEndian)
}

func (this *SrsTsPayloadPATProgram) Size() uint32 {
	return 4
}

type SrsTsPayloadPAT struct {
	psiHeader *SrsTsPayloadPSI
	/*
		This is a 16-bit field which serves as a label to identify this Transport Stream from any other
		multiplex within a network. Its value is defined by the user
	*/
	transportStreamId int16 //固定为0x0001
	const1Value0      int8  //2bits
	/*
		This 5-bit field is the version number of the whole Program Association Table. The version number
		shall be incremented by 1 modulo 32 whenever the definition of the Program Association Table changes. When the
		current_next_indicator is set to '1', then the version_number shall be that of the currently applicable Program Association
		Table. When the current_next_indicator is set to '0', then the version_number shall be that of the next applicable Program
		Association Table.
	*/
	versionNumber int8 //5bits	版本号，固定为二进制00000，如果PAT有变化则版本号加1
	/*
		A 1-bit indicator, which when set to '1' indicates that the Program Association Table sent is
		currently applicable. When the bit is set to '0', it indicates that the table sent is not yet applicable and shall be the next
		table to become valid.
	*/
	currentNextIndicator int8 //固定为二进制1，表示这个PAT表可以用，如果为0则要等待下一个PAT表
	/*
		This 8-bit field gives the number of this section. The section_number of the first section in the
		Program Association Table shall be 0x00. It shall be incremented by 1 with each additional section in the Program
		Association Table.
	*/
	sectionNumber int8 //固定为0x00
	/*
		This 8-bit field specifies the number of the last section (that is, the section with the highest
		section_number) of the complete Program Association Table.
	*/
	lastSectionNumber int8 //固定为0x00
	// multiple 4B program data.
	programs []*SrsTsPayloadPATProgram
	// 4B
	/**
	 * This is a 32-bit field that contains the CRC value that gives a zero output of the registers in the decoder
	 * defined in Annex A after processing the entire section.
	 * @remark crc32(bytes without pointer field, before crc32 field)
	 */
	crc32 int32
}

func NewSrsTsPayloadPAT(p *SrsTsPacket) *SrsTsPayloadPAT {
	return &SrsTsPayloadPAT{
		psiHeader: NewSrsTsPayloadPSI(p),
	}
}

func CreatePAT(context *SrsTsContext, pmt_number int16, pmt_pid int16) *SrsTsPacket {
	pkt := NewSrsTsPacket()

	pkt.tsHeader.syncByte = SRS_TS_SYNC_BYTE
	pkt.tsHeader.transportErrorIndicator = 0
	pkt.tsHeader.payloadUnitStartIndicator = 1
	pkt.tsHeader.transportPriority = 0
	pkt.tsHeader.PID = SrsTsPidPAT
	pkt.tsHeader.transportScrambingControl = SrsTsScrambledDisabled
	pkt.tsHeader.adaptationFieldControl = SrsTsAdapationControlPayloadOnly
	pkt.tsHeader.continuityCounter = 0

	pat := NewSrsTsPayloadPAT(pkt)
	pat.psiHeader.pointerField = 0
	pat.psiHeader.tableId = SrsTsPsiTableIdPas
	pat.psiHeader.sectionSyntaxIndicator = 1
	pat.psiHeader.const0Value = 0
	pat.psiHeader.const1Value0 = 0x03 //2bits
	pat.psiHeader.sectionLength = 0   //calc in size

	pat.transportStreamId = 1
	pat.const1Value0 = 0x3 //2bits
	pat.versionNumber = 0
	pat.currentNextIndicator = 1
	pat.sectionNumber = 0
	pat.lastSectionNumber = 0
	program := NewSrsTsPayloadPATProgram(pmt_number, pmt_pid)
	pat.programs = append(pat.programs, program)
	//calc section length
	pat.psiHeader.sectionLength = int16(pat.Size())
	//填充payload
	s := utils.NewSrsStream([]byte{})
	pat.Encode(s)
	pkt.payload = s.Data()
	return pkt
}

func (this *SrsTsPayloadPAT) Encode(stream *utils.SrsStream) {
	s := utils.NewSrsStream([]byte{})//3
	this.psiHeader.Encode(s) //5
	s.WriteInt16(this.transportStreamId, binary.BigEndian)
	this.const1Value0 = 0x03
	var b byte = 0
	b |= byte(this.currentNextIndicator & 0x01)
	b |= byte((this.versionNumber << 1) & 0x3e)
	b |= byte(this.const1Value0 << 6) & 0xC0
	s.WriteByte(b)

	s.WriteByte(byte(this.sectionNumber))
	s.WriteByte(byte(this.lastSectionNumber))//5
	for i := 0; i < len(this.programs); i++ {
		this.programs[i].Encode(s)//4
	}

	CRC32 := utils.MpegtsCRC32(s.Data()[1:])
	s.WriteInt32(int32(CRC32), binary.BigEndian)//4
	stream.WriteBytes(s.Data())
	if len(stream.Data()) + 4 < 188 {
		i := 188 - len(stream.Data()) - 4
		for j := 0; j < i; j++ {
			stream.WriteByte(0xff)
		}
	}
}

func (this *SrsTsPayloadPAT) Decode(stream *utils.SrsStream) error {
	return nil
}

func (this *SrsTsPayloadPAT) Size() uint32 {
	var m uint32 = 0
	for i := 0; i < len(this.programs); i++ {
		m += this.programs[i].Size()
	}

	return 5 + m + 4
}
