package app

import (
	"errors"
	"fmt"
	"encoding/binary"
	"go_srs/srs/codec"
	"go_srs/srs/utils"
)

/**
* the h264/avc and aac codec, for media stream.
*
* to demux the FLV/RTMP video/audio packet to sample,
* add each NALUs of h.264 as a sample unit to sample,
* while the entire aac raw data as a sample unit.
*
* for sequence header,
* demux it and save it in the avc_extra_data and aac_extra_data,
*
* for the codec info, such as audio sample rate,
* decode from FLV/RTMP header, then use codec info in sequence
* header to override it.
 */
type SrsAvcAacCodec struct {
	stream *utils.SrsStream
	/*
	* metadata
	 */
	duration      int
	width         int
	height        int
	frameRate     int
	videoCodecId  int
	videoDataRate int
	audioDataRate int
	audioCodecId  int
	// profile_idc, H.264-AVC-ISO_IEC_14496-10.pdf, page 45.
	avcProfile codec.SrsAvcProfile
	// level_idc, H.264-AVC-ISO_IEC_14496-10.pdf, page 45.
	avcLevel                    codec.SrsAvcLevel
	NALUintLength               int8
	sequenceParameterSetLength  int16
	sequenceParameterSetNALUnit []byte
	pictureParameterSetLength   int16
	pictureParameterSetNALUnit  []byte

	payloadFormat      codec.SrsAvcPayloadFormat
	aacObject          codec.SrsAacObjectType
	aacSampleRateIndex int8
	aacChannels        int8

	avcExtraSize int
	avcExtraData []byte

	aacExtraSize int
	aacExtraData []byte

	avcParseSps bool
}

func NewSrsAvcAacCodec() *SrsAvcAacCodec {
	return &SrsAvcAacCodec{
		avcParseSps:   true,
		width:         0,
		height:        0,
		duration:      0,
		NALUintLength: 0,
		frameRate:     0,
		videoDataRate: 0,
		videoCodecId:  0,
		audioDataRate: 0,
		audioCodecId:  0,

		avcProfile:         codec.SrsAvcProfileReserved,
		avcLevel:           codec.SrsAvcLevelReserved,
		aacObject:          codec.SrsAacObjectTypeReserved,
		aacSampleRateIndex: codec.SRS_AAC_SAMPLE_RATE_UNSET,
		aacChannels:        0,
		avcExtraSize:       0,
		aacExtraSize:       0,

		sequenceParameterSetLength: 0,
		pictureParameterSetLength:  0,
		payloadFormat:              codec.SrsAvcPayloadFormatGuess,
		stream:                     utils.NewSrsStream([]byte{}),
	}
}

func (this *SrsAvcAacCodec) is_avc_codec_ok() bool {
	return this.avcExtraSize > 0 && this.avcExtraData != nil
}

func (this *SrsAvcAacCodec) is_aac_codec_ok() bool {
	return this.aacExtraSize > 0 && this.aacExtraData != nil
}

func (this *SrsAvcAacCodec) audio_aac_demux(data []byte, sample *SrsCodecSample) error {
	sample.SetIsVideo(false)

	stream := utils.NewSrsStream(data)

	soundFormat, err := stream.ReadByte()
	if err != nil {
		return err
	}
	soundType := soundFormat & 0x01
	soundSize := (soundFormat >> 1) & 0x01
	soundRate := (soundFormat >> 2) & 0x03
	soundFormat = (soundFormat >> 4) & 0x0f

	this.audioCodecId = int(soundFormat)
	sample.ACodec = codec.SrsCodecAudio(this.audioCodecId)
	sample.SoundType = codec.SrsCodecAudioSoundType(soundType)
	sample.SoundRate = codec.SrsCodecAudioSampleRate(soundRate)
	sample.SoundSize = codec.SrsCodecAudioSampleSize(soundSize)

	if this.audioCodecId == codec.SrsCodecAudioMP3 {
		return errors.New("error hls try mp3")
	}

	if this.audioCodecId != codec.SrsCodecAudioAAC {
		return errors.New("aac only support mp3/aac codec")
	}

	aacPacketType, err := stream.ReadByte()
	if err != nil {
		return err
	}

	sample.AacPacketType = codec.SrsCodecAudioType(aacPacketType)
	if aacPacketType == codec.SrsCodecAudioTypeSequenceHeader {
		this.aacExtraData = stream.ReadLeftBytes()
		this.aacExtraSize = len(this.aacExtraData)
		if err := this.audio_aac_sequence_header_demux(this.aacExtraData); err != nil {
			return err
		}
	} else if aacPacketType == codec.SrsCodecAudioTypeRawData {
		if !this.is_aac_codec_ok() {
			return fmt.Errorf("aac ignore type=%d for no sequence header. ret=%d", aacPacketType)
		}
		// Raw AAC frame data in UI8 []
		// 6.3 Raw Data, aac-iso-13818-7.pdf, page 28
		if err := sample.AddSampleUnit(stream.ReadLeftBytes()); err != nil {
			return errors.New("aac add sample failed.")
		}
	}

	// reset the sample rate by sequence header
	if this.aacSampleRateIndex != codec.SRS_AAC_SAMPLE_RATE_UNSET {
		aacSampleRates := []int{
			96000, 88200, 64000, 48000,
			44100, 32000, 24000, 22050,
			16000, 12000, 11025, 8000,
			7350, 0, 0, 0,
		}

		switch aacSampleRates[this.aacSampleRateIndex] {
		case 11025:
			sample.SoundRate = codec.SrsCodecAudioSampleRate11025
			break
		case 22050:
			sample.SoundRate = codec.SrsCodecAudioSampleRate22050
			break
		case 44100:
			sample.SoundRate = codec.SrsCodecAudioSampleRate44100
			break
		default:
			break
		}
	}

	return nil
}

func (this *SrsAvcAacCodec) audio_aac_sequence_header_demux(data []byte) error {
	stream := utils.NewSrsStream(data)
	// only need to decode the first 2bytes:
	//      audioObjectType, aac_profile, 5bits.
	//      samplingFrequencyIndex, aac_sample_rate, 4bits.
	//      channelConfiguration, aac_channels, 4bits
	profileObjectType, err := stream.ReadByte()
	if err != nil {
		return err
	}
	samplingFrequencyIndex, err := stream.ReadByte()
	if err != nil {
		return err
	}

	this.aacChannels = (int8(samplingFrequencyIndex) >> 3) & 0x0f
	samplingFrequencyIndex = ((profileObjectType << 1) & 0x0e) | ((samplingFrequencyIndex >> 7) & 0x01)
	profileObjectType = (profileObjectType >> 3) & 0x1f

	this.aacSampleRateIndex = int8(samplingFrequencyIndex)

	// convert the object type in sequence header to aac profile of ADTS.
	this.aacObject = codec.SrsAacObjectType(profileObjectType)
	if this.aacObject == codec.SrsAacObjectTypeReserved {
		return errors.New("audio codec decode aac sequence header failed, adts object invalid")
	}

	// TODO: FIXME: to support aac he/he-v2, see: ngx_rtmp_codec_parse_aac_header
	// @see: https://github.com/winlinvip/nginx-rtmp-module/commit/3a5f9eea78fc8d11e8be922aea9ac349b9dcbfc2
	//
	// donot force to LC, @see: https://github.com/ossrs/srs/issues/81
	// the source will print the sequence header info.
	//if (aac_profile > 3) {
	// Mark all extended profiles as LC
	// to make Android as happy as possible.
	// @see: ngx_rtmp_hls_parse_aac_header
	//aac_profile = 1;
	//}

	return nil
}

func (this *SrsAvcAacCodec) video_avc_demux(data []byte, sample *SrsCodecSample) error {
	sample.SetIsVideo(true)

	stream := utils.NewSrsStream(data)
	// @see: E.4.3 Video Tags, video_file_format_spec_v10_1.pdf, page 78
	frameType, err := stream.ReadByte()
	if err != nil {
		return err
	}

	codecId := frameType & 0x0f
	frameType = (frameType >> 4) & 0x0f

	sample.FrameType = codec.SrsCodecVideoAVCFrame(frameType)

	if sample.FrameType == codec.SrsCodecVideoAVCFrameVideoInfoFrame {
		return errors.New("avc ignore the info frame")
	}

	if codecId != codec.SrsCodecVideoAVC {
		return errors.New("avc only support video h.264/avc codec")
	}

	this.videoCodecId = int(codecId)

	avcPacketType, err := stream.ReadByte()
	if err != nil {
		return err
	}

	ctsTmp, err := stream.ReadBytes(3)
	if err != nil {
		return err
	}

	var cts int32 = 0
	cts |= int32(ctsTmp[2])
	cts |= int32(ctsTmp[1] << 8)
	cts |= int32(ctsTmp[0] << 16)
	//pts = dts + cts
	sample.Cts = cts
	sample.AvcPacketType = codec.SrsCodecVideoAVCType(avcPacketType)

	if avcPacketType == codec.SrsCodecVideoAVCTypeSequenceHeader {

	}
	return nil
}

func (this *SrsAvcAacCodec) avc_demux_sps_pps(stream *utils.SrsStream) error {
	sps_pps := stream.ReadLeftBytes()
	if len(sps_pps) <= 0 {
		return errors.New("the length of sps and pps must big than 0")
	}

	//int8_t configurationVersion = stream->read_1bytes();
	_, err := stream.ReadByte()
	if err != nil {
		return err
	}
	//int8_t AVCProfileIndication = stream->read_1bytes();
	tmp, err1 := stream.ReadByte()
	if err1 != nil {
		return err1
	}
	this.avcProfile = codec.SrsAvcProfile(tmp)

	_, err = stream.ReadByte()
	if err != nil {
		return err
	}
	//int8_t AVCLevelIndication = stream->read_1bytes();
	tmp, err2 := stream.ReadByte()
	if err2 != nil {
		return err2
	}
	this.avcLevel = codec.SrsAvcLevel(tmp)
	// parse the NALU size.
	lengthSizeMinusOne, err3 := stream.ReadInt8()
	if err3 != nil {
		return err3
	}
	lengthSizeMinusOne &= 0x03
	this.NALUintLength = lengthSizeMinusOne
	// 5.3.4.2.1 Syntax, H.264-AVC-ISO_IEC_14496-15.pdf, page 16
	// 5.2.4.1 AVC decoder configuration record
	// 5.2.4.1.2 Semantics
	// The value of this field shall be one of 0, 1, or 3 corresponding to a
	// length encoded with 1, 2, or 4 bytes, respectively.
	if this.NALUintLength == 2 {
		return errors.New("sps lengthSizeMinusOne should never be 2")
	}

	// 1 sps, 7.3.2.1 Sequence parameter set RBSP syntax
	// H.264-AVC-ISO_IEC_14496-10.pdf, page 45.
	numOfSequenceParameterSets, err4 := stream.ReadInt8()
	if err4 != nil {
		return err4
	}
	numOfSequenceParameterSets &= 0x1f
	if numOfSequenceParameterSets != 1 {
		return errors.New("avc decode sequence header sps failed")
	}
	this.sequenceParameterSetLength, err = stream.ReadInt16(binary.BigEndian)
	if err != nil {
		return err
	}

	if this.sequenceParameterSetLength > 0 {
		this.sequenceParameterSetNALUnit, err = stream.ReadBytes(uint32(this.sequenceParameterSetLength))
		if err != nil {
			return err
		}
	}

	numOfPictureParameterSets, err7 := stream.ReadInt8()
	if err7 != nil {
		return err7
	}
	numOfPictureParameterSets &= 0x1f
	if numOfPictureParameterSets != 1 {
		return errors.New("avc decode sequence header pps failed")
	}
	this.pictureParameterSetLength, err = stream.ReadInt16(binary.BigEndian)
	if err != nil {
		return err
	}

	if this.pictureParameterSetLength > 0 {
		this.pictureParameterSetNALUnit, err = stream.ReadBytes(uint32(this.pictureParameterSetLength))
		if err != nil {
			return err
		}
	}
	return nil
}

func (this *SrsAvcAacCodec) avc_demux_sps() error {
	if this.sequenceParameterSetLength <= 0 {
		return nil
	}

	stream := utils.NewSrsStream(this.sequenceParameterSetNALUnit)
	// for NALU, 7.3.1 NAL unit syntax
	// H.264-AVC-ISO_IEC_14496-10-2012.pdf, page 61.
	nutv, err := stream.ReadInt8()
	if err != nil {
		return err
	}

	forbidden_zero_bit := (nutv >> 7) & 0x01
	if forbidden_zero_bit != 0 {
		return errors.New("forbidden_zero_bit shall be equal to 0")
	}
	// nal_ref_idc not equal to 0 specifies that the content of the NAL unit contains a sequence parameter set or a picture
	// parameter set or a slice of a reference picture or a slice data partition of a reference picture.
	nal_ref_idc := (nutv >> 5) & 0x03
	if nal_ref_idc == 0 {
		return errors.New("for sps, nal_ref_idc shall be not be equal to 0.")
	}
	// 7.4.1 NAL unit semantics
    // H.264-AVC-ISO_IEC_14496-10-2012.pdf, page 61.
	// nal_unit_type specifies the type of RBSP data structure contained in the NAL unit as specified in Table 7-1.
	nal_unit_type := codec.SrsAvcNaluType(nutv & 0x1f)
	if nal_unit_type != 7 {
		return errors.New("for sps, nal_unit_type shall be equal to 7")
	}
	// decode the rbsp from sps.
	// rbsp[ i ] a raw byte sequence payload is specified as an ordered sequence of bytes.
	rbsp := make([]byte, 0)
	for !stream.Empty() {
		b, err := stream.ReadByte()
		if err != nil {
			return err
		}

		rbsp = append(rbsp, b)
		nb_rbsp := len(rbsp)
		// XX 00 00 03 XX, the 03 byte should be drop.
		if nb_rbsp > 2 && rbsp[nb_rbsp - 2] == 0 && rbsp[nb_rbsp - 1] == 0 && rbsp[nb_rbsp] == 3 {
			if stream.Empty() {
				break
			}
			c, err := stream.ReadByte()
			if err != nil {
				return err
			}
			rbsp[len(rbsp) - 1] = c
		}
	}
	return this.avc_demux_sps_rbsp(rbsp)
}

func (this *SrsAvcAacCodec) avc_demux_sps_rbsp(rbsp []byte) error {
	// we donot parse the detail of sps.
	// @see https://github.com/ossrs/srs/issues/474
	if !this.avcParseSps {
		return nil
	}

	stream := utils.NewSrsStream(rbsp)
	// for SPS, 7.3.2.1.1 Sequence parameter set data syntax
	// H.264-AVC-ISO_IEC_14496-10-2012.pdf, page 62.
	prifile_idc, err := stream.ReadInt8()
	if err != nil {
		return err
	}

	if prifile_idc == 0 {
		return errors.New("sps the profile_idc invalid")
	}

	flags, err := stream.ReadInt8()
	if err != nil {
		return err
	}

	if (flags & 0x03) != 0 {
		return errors.New("sps the flags invalid.")
	}

	level_idc, err := stream.ReadInt8()
	if err != nil {
		return err
	}

	if level_idc == 0 {
		return errors.New("sps the level_idc invalid.")
	}

	bs := utils.NewSrsBitStream(stream.ReadLeftBytes())
	_ = bs
	var seq_parameter_set_id int32 = -1
	_ = seq_parameter_set_id
	return nil
}

func srs_avc_nalu_read_uev(bs *utils.SrsBitStream, v *int32) error {
	if bs.Empty() {
		return errors.New("no enougth data")
	}

	// ue(v) in 9.1 Parsing process for Exp-Golomb codes
    // H.264-AVC-ISO_IEC_14496-10-2012.pdf, page 227.
    // Syntax elements coded as ue(v), me(v), or se(v) are Exp-Golomb-coded.
    //      leadingZeroBits = -1;
    //      for( b = 0; !b; leadingZeroBits++ )
    //          b = read_bits( 1 )
    // The variable codeNum is then assigned as follows:
	//      codeNum = (2<<leadingZeroBits) - 1 + read_bits( leadingZeroBits )
	
	var leadingZeroBits int = -1
	var b int8 = 0
	var err error
	
	for b = 0; b == 0 && !bs.Empty(); leadingZeroBits++ {
		b, err = bs.ReadBit()
		if err != nil {
			return err
		}
	}

	if leadingZeroBits >= 31 {
		return errors.New("")
	}

	*v = (1 << uint(leadingZeroBits)) - 1
	for i := 0; i < leadingZeroBits; i++ {
		b, err = bs.ReadBit()
		if err != nil {
			return err
		}
		*v += int32(b) << uint(leadingZeroBits - 1 - i)
	}
	return nil
}
