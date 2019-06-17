package packet

import "go_srs/srs/protocol/amf0"
import "go_srs/srs/global"
import "go_srs/srs/utils"

type SrsSampleAccessPacket struct {
	CommandName amf0.SrsAmf0String
	/**
	 * whether allow access the sample of video.
	 * @see: https://github.com/ossrs/srs/issues/49
	 * @see: http://help.adobe.com/en_US/FlashPlatform/reference/actionscript/3/flash/net/NetStream.html#videoSampleAccess
	 */
	VideoSampleAccess amf0.SrsAmf0Boolean
	/**
	 * whether allow access the sample of audio.
	 * @see: https://github.com/ossrs/srs/issues/49
	 * @see: http://help.adobe.com/en_US/FlashPlatform/reference/actionscript/3/flash/net/NetStream.html#audioSampleAccess
	 */
	AudioSampleAccess amf0.SrsAmf0Boolean
}

func NewSrsSampleAccessPacket() *SrsSampleAccessPacket {
	return &SrsSampleAccessPacket{
		CommandName: amf0.SrsAmf0String{Value: amf0.SrsAmf0Utf8{Value: amf0.RTMP_AMF0_DATA_SAMPLE_ACCESS}},
		VideoSampleAccess:amf0.SrsAmf0Boolean{Value:false},
		AudioSampleAccess:amf0.SrsAmf0Boolean{Value:false},
	}
}

func (this *SrsSampleAccessPacket) GetMessageType() int8 {
	return global.RTMP_MSG_AMF0DataMessage
}

func (this *SrsSampleAccessPacket) GetPreferCid() int32 {
	return global.RTMP_CID_OverStream
}

func (this *SrsSampleAccessPacket) Decode(stream *utils.SrsStream) error {
	return nil
}

func (this *SrsSampleAccessPacket) Encode(stream *utils.SrsStream) error {
	_ = this.CommandName.Encode(stream)
	_ = this.VideoSampleAccess.Encode(stream)
	_ = this.AudioSampleAccess.Encode(stream)
	return nil
}


