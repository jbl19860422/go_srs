package app


/**
* the profile for avc/h.264.
* @see Annex A Profiles and levels, H.264-AVC-ISO_IEC_14496-10.pdf, page 205.
*/
type SrsAvcProfile int
const (
	_ SrsAvcProfile = iota
	SrsAvcProfileReserved = 0
    
    // @see ffmpeg, libavcodec/avcodec.h:2713
    SrsAvcProfileBaseline = 66
    // FF_PROFILE_H264_CONSTRAINED  (1<<9)  // 8+1; constraint_set1_flag
    // FF_PROFILE_H264_CONSTRAINED_BASELINE (66|FF_PROFILE_H264_CONSTRAINED)
    SrsAvcProfileConstrainedBaseline = 578
    SrsAvcProfileMain = 77
    SrsAvcProfileExtended = 88
    SrsAvcProfileHigh = 100
    SrsAvcProfileHigh10 = 110
    SrsAvcProfileHigh10Intra = 2158
    SrsAvcProfileHigh422 = 122
    SrsAvcProfileHigh422Intra = 2170
    SrsAvcProfileHigh444 = 144
    SrsAvcProfileHigh444Predictive = 244
    SrsAvcProfileHigh444Intra = 2192
)

/**
* the level for avc/h.264.
* @see Annex A Profiles and levels, H.264-AVC-ISO_IEC_14496-10.pdf, page 207.
*/
type SrsAvcLevel int
const (
	_ SrsAvcLevel = iota
	SrsAvcLevelReserved = 0
    
    SrsAvcLevel_1 = 10
    SrsAvcLevel_11 = 11
    SrsAvcLevel_12 = 12
    SrsAvcLevel_13 = 13
    SrsAvcLevel_2 = 20
    SrsAvcLevel_21 = 21
    SrsAvcLevel_22 = 22
    SrsAvcLevel_3 = 30
    SrsAvcLevel_31 = 31
    SrsAvcLevel_32 = 32
    SrsAvcLevel_4 = 40
    SrsAvcLevel_41 = 41
    SrsAvcLevel_5 = 50
    SrsAvcLevel_51 = 51
)

/**
* the avc payload format, must be ibmf or annexb format.
* we guess by annexb first, then ibmf for the first time,
* and we always use the guessed format for the next time.
*/
type SrsAvcPayloadFormat int
const (
	_ SrsAvcPayloadFormat = iota
	SrsAvcPayloadFormatGuess = 0
    SrsAvcPayloadFormatAnnexb = 1
    SrsAvcPayloadFormatIbmf = 2
)


/**
* the aac object type, for RTMP sequence header
* for AudioSpecificConfig, @see aac-mp4a-format-ISO_IEC_14496-3+2001.pdf, page 33
* for audioObjectType, @see aac-mp4a-format-ISO_IEC_14496-3+2001.pdf, page 23
*/
type SrsAacObjectType int
const (
	_ SrsAacObjectType = iota
    SrsAacObjectTypeReserved = 0
    
    // Table 1.1 - Audio Object Type definition
    // @see @see aac-mp4a-format-ISO_IEC_14496-3+2001.pdf, page 23
    SrsAacObjectTypeAacMain = 1
    SrsAacObjectTypeAacLC = 2
    SrsAacObjectTypeAacSSR = 3
    
    // AAC HE = LC+SBR
    SrsAacObjectTypeAacHE = 5
    // AAC HEv2 = LC+SBR+PS
    SrsAacObjectTypeAacHEV2 = 29
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
	duration 		int
	width			int
	height			int
	frameRate		int
	videoCodecId	int
	videoDataRate	int
	audioDataRate	int
	audioCodecId	int
	// profile_idc, H.264-AVC-ISO_IEC_14496-10.pdf, page 45.
	avcProfile		SrsAvcProfile
	// level_idc, H.264-AVC-ISO_IEC_14496-10.pdf, page 45.
    avcLevel		SrsAvcLevel 
	NALUintLength	int8
	sequenceParameterSetLength 	int16
	sequenceParameterSetNALUnit []byte
	pictureParameterSetLength	int16
	pictureParameterSetNALUnit	[]byte

	payloadFormat	SrsAvcPayloadFormat
	aacObject		SrsAacObjectType
	aacSampleRate	int8
	aacChannels		int8

	avcExtraSize	int
	avcExtraData	[]byte

	aacExtraSize	int
	aacExtraData	[]byte

	avcParseSps		bool
}

func NewSrsAvcAacCodec() *SrsAvcAacCodec {
	return &SrsAvcAacCodec{
		avcParseSps:true,
		width:0,
		height:0,
		duration:0,
		NALUintLength:0,
		frameRate:0,
		videoDataRate:0,
		videoCodecId:0,
		audioDataRate:0,
		audioCodecId:0,

		avcProfile:SrsAvcProfileReserved,
		avcLevel:SrsAvcLevelReserved,
		aacObject:SrsAacObjectTypeReserved,
		aacSampleRateIndex:SRS_AAC_SAMPLE_RATE_UNSET,
		aacChannels:0,
		avcExtraSize:0,
		aacExtraSize:0,
		
		sequenceParameterSetLength:0,
		pictureParameterSetLength:0,
		payloadFormat:SrsAvcPayloadFormatGuess
		stream:utils.NewSrsStream([]byte{})
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
	soundFomat = (soundFomat >> 4) & 0x0f

	this.audioCodecId = soundFomat
	sample.ACodec = (SrsCodecAudio)this.audioCodecId
	sample.SoundType = SrsCodecAudioSoundType(soundType)
	sample.soundRate = SrsCodecAudioSampleRate(soundRate)
	sample.soundSize = SrsCodecAudioSampleSize(soundSize)

	if this.audioCodecId == SrsCodecAudioMP3 {
		return errors.New("error hls try mp3")
	}

	if this.audioCodecId != SrsCodecAudioAAC {
		return errors.New("aac only support mp3/aac codec")
	}

	aacPacketType, err := stream.ReadByte()
	if err != nil {
		return err
	}

	sample.AacPacketType = aacPacketType
	if aacPacketType == SrsCodecAudioTypeSequenceHeader {
		this.aacExtraData = stream.ReadLeftBytes()
		this.aacExtraSize = len(this.aacExtraData)
		if err := this.audio_aac_sequence_header_demux(this.aacExtraData); err != nil {
			return err
		}
	} else if aacPacketType == SrsCodecAudioTypeRawData {
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
    if this.aacSampleRate != SRS_AAC_SAMPLE_RATE_UNSET {
        aacSampleRates = []int{
            96000, 88200, 64000, 48000,
            44100, 32000, 24000, 22050,
            16000, 12000, 11025,  8000,
            7350,     0,     0,    0
		}
		
        switch aacSampleRates[this.aacSampleRateIndex] {
            case 11025:
                sample.SoundRate = SrsCodecAudioSampleRate11025;
                break;
            case 22050:
                sample.SoundRate = SrsCodecAudioSampleRate22050;
                break;
            case 44100:
                sample.SoundRate = SrsCodecAudioSampleRate44100;
                break;
            default:
                break;
        };
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

	this.aacChannels = (samplingFrequencyIndex >> 3) & 0x0f
	samplingFrequencyIndex = ((profileObjectType << 1) & 0x0e) | ((samplingFrequencyIndex >> 7) & 0x01)
	profileObjectType = (profileObjectType >> 3) & 0x1f

	this.aacSampleRateIndex = samplingFrequencyIndex

    // convert the object type in sequence header to aac profile of ADTS.
    this.aacObject = (SrsAacObjectType)profileObjectType
    if (this.aacObject == SrsAacObjectTypeReserved) {
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

    return ret;
}