package rtmp

import (
	_ "context"
	"log"
	"net"
	_ "net/url"
	_ "strings"
	_ "time"
	"go_srs/srs/protocol/skt"
	"go_srs/srs/protocol/packet"
	"go_srs/srs/protocol/amf0"
	"go_srs/srs/global"
	// "fmt"
)

type SrsRtmpServer struct {
	io   			*skt.SrsIOReadWriter
	Protocol    	*SrsProtocol
	HandShaker  	*SrsSimpleHandShake
	IOErrListener 	skt.SrsIOErrListener
}

func NewSrsRtmpServer(conn net.Conn, listener skt.SrsIOErrListener) *SrsRtmpServer {
	io_ := skt.NewSrsIOReadWriter(conn)
	return &SrsRtmpServer{
		io: io_,
		Protocol: NewSrsProtocol(io_), 
		HandShaker: NewSrsSimpleHandShake(io_),
		IOErrListener:listener,
	}
}

func (this *SrsRtmpServer) Close() {
	this.io.Close()
}

func (this *SrsRtmpServer) GetClientIP() string {
	return this.io.GetClientIP()
}

func (this *SrsRtmpServer) HandShake() error {
	err := this.HandShaker.HandShakeWithClient()
	return err
}

func (this *SrsRtmpServer) ConnectApp() (packet.SrsPacket, error) {
	connPacket := packet.NewSrsConnectAppPacket()
	if err := this.Protocol.ExpectMessage(connPacket); err != nil {
		return nil, err
	}
	return connPacket, nil
}

func (this *SrsRtmpServer) RecvMessage() (*SrsRtmpMessage, error) {
	return this.Protocol.RecvMessage()
}

func (this *SrsRtmpServer) DecodeMessage(msg *SrsRtmpMessage) (packet.SrsPacket, error) {
	return this.Protocol.DecodeMessage(msg)
}

func (this *SrsRtmpServer) OnRecvError(err error) {
	if this.IOErrListener != nil {
		this.IOErrListener.OnRecvError(err)
	}
}

func (this *SrsRtmpServer) IdentifyClient(streamId int) (SrsRtmpConnType, string, float64, error) {
	var typ SrsRtmpConnType
	var streamname string
	var duration float64
	for {
		msg, err := this.Protocol.RecvMessage()
		if err != nil {
			continue
		}
		header := msg.GetHeader()
		if header.IsAckledgement() || header.IsSetChunkSize() || header.IsWindowAckledgementSize() || header.IsUserControlMessage() {
			continue
		}

		if !header.IsAmf0Command() && !header.IsAmf3Command() {
			continue
		}

		pkt, err := this.Protocol.DecodeMessage(msg)
		switch pkt.(type) {
			//todo
			case *packet.SrsCreateStreamPacket: {
				typ, streamname, duration, err = this.identify_create_stream_client(pkt.(*packet.SrsCreateStreamPacket), streamId)
				return typ, streamname, duration, err
			}
			case *packet.SrsFMLEStartPacket: {
				typ, streamname, err = this.identify_fmle_publish_client(pkt.(*packet.SrsFMLEStartPacket))
				if err != nil {
					return typ, streamname, 0, nil
				}
				return typ, streamname, 0, nil
			}
			case *packet.SrsPlayPacket:{
				typ, streamname, duration, err = this.identify_play_client(pkt.(*packet.SrsPlayPacket))
				return typ, streamname, duration, err
			}
		}
		return typ, streamname, 0, nil
	}
	_ = typ
	return typ, streamname, 0, nil
}

func (this *SrsRtmpServer) identify_create_stream_client(req *packet.SrsCreateStreamPacket, streamId int) (SrsRtmpConnType, string, float64, error) {
	typ := SrsRtmpConnType(SrsRtmpConnFMLEPublish)
	resPkt := packet.NewSrsCreateStreamResPacket(req.TransactionId.GetValue().(float64), float64(streamId))
	err := this.Protocol.SendPacket(resPkt, 0)
	var streamname string
	var duration float64
	if err != nil {
		return SrsRtmpConnType(SrsRtmpConnFMLEPublish), "", 0, err
	}

	for {
		msg, err := this.Protocol.RecvMessage()
		if err != nil {
			// log.Print("identify_client err, msg=", err)
			continue
		}
		header := msg.GetHeader()
		if header.IsAckledgement() || header.IsSetChunkSize() || header.IsWindowAckledgementSize() || header.IsUserControlMessage() {
			continue
		}

		if !header.IsAmf0Command() && !header.IsAmf3Command() {
			continue
		}

		pkt, err := this.Protocol.DecodeMessage(msg)
		switch pkt.(type) {
			case *packet.SrsPlayPacket:{
				log.Print("SrsPlayPacket")
				typ, streamname, duration, err = this.identify_play_client(pkt.(*packet.SrsPlayPacket))
				return typ, streamname, duration, err
			}
			case *packet.SrsCreateStreamPacket: {
				// log.Print("SrsCreateStreamPacket")
				typ, streamname, duration, err = this.identify_create_stream_client(pkt.(*packet.SrsCreateStreamPacket), streamId)
				return typ, streamname, duration, err
			}
			case *packet.SrsFMLEStartPacket: {
				// log.Print("SrsFMLEStartPacket streamname=", pkt.(*packet.SrsFMLEStartPacket).StreamName)
				typ, streamname, err = this.identify_fmle_publish_client(pkt.(*packet.SrsFMLEStartPacket))
				if err != nil {
					log.Print("identify_fmle_publish_client reeturn")
					return typ, streamname, 0, nil
				}
				return typ, streamname, 0, nil
			}
			//todo
			//identify_haivision_publish_client
			//identify_flash_publish_client
		}
	}
	_ = typ
	return typ, streamname, 0, nil
}

func (this *SrsRtmpServer) identify_play_client(pkt *packet.SrsPlayPacket) (SrsRtmpConnType, string, float64, error) {
	return SrsRtmpConnPlay, pkt.StreamName.GetValue().(string), pkt.Duration.GetValue().(float64),nil
}

func (this *SrsRtmpServer) SendMsg(msg *SrsRtmpMessage, streamId int) error {
	msgs := make([]*SrsRtmpMessage, 1)
	msgs[0] = msg
	return this.Protocol.SendMessages(msgs, streamId)
}

func (this *SrsRtmpServer) identify_fmle_publish_client(req *packet.SrsFMLEStartPacket) (SrsRtmpConnType, string, error) {
	typ := SrsRtmpConnType(SrsRtmpConnFMLEPublish)
	pkt := packet.NewSrsFMLEStartResPacket(req.TransactionId.Value)
	err := this.Protocol.SendPacket(pkt, 0)
	if err != nil {
		return typ, req.StreamName.Value.Value, err
	}
	return typ, req.StreamName.Value.Value, nil
}

func (this *SrsRtmpServer) StartPlay(streamId int) error {
	 // StreamBegin
	pkt := packet.NewSrsUserControlPacket()
	pkt.EventType = global.SrcPCUCStreamBegin
	pkt.EventData = int32(streamId)
	err := this.Protocol.SendPacket(pkt, 0)
	if err != nil {
		return err
	}

	// onStatus(NetStream.Play.Reset)
	callPkt := packet.NewSrsOnStatusCallPacket()
	callPkt.Data.Set(global.StatusLevel, global.StatusLevelStatus)
	callPkt.Data.Set(global.StatusCode, global.StatusCodeStreamReset)
	callPkt.Data.Set(global.StatusDescription, "Playing and resetting stream.")
	callPkt.Data.Set(global.StatusDetails, "stream")
	callPkt.Data.Set(global.StatusClientId, global.RTMP_SIG_CLIENT_ID)
	
	err = this.Protocol.SendPacket(callPkt, int32(streamId))
	if err != nil {
		return err
	}

	// onStatus(NetStream.Play.Start)
	startPkt := packet.NewSrsOnStatusCallPacket()
	startPkt.Data.Set(global.StatusLevel, global.StatusLevelStatus)
	startPkt.Data.Set(global.StatusCode, global.StatusCodeStreamStart)
	startPkt.Data.Set(global.StatusDescription, "Started playing stream.")
	startPkt.Data.Set(global.StatusDetails, "stream")
	startPkt.Data.Set(global.StatusClientId, global.RTMP_SIG_CLIENT_ID)
	// |RtmpSampleAccess(false, false)
	err = this.Protocol.SendPacket(startPkt, int32(streamId))
	if err != nil {
		return err
	}
	// allow audio/video sample.
    // @see: https://github.com/ossrs/srs/issues/49
	samplePkt := packet.NewSrsSampleAccessPacket()
	samplePkt.VideoSampleAccess.Value = true
	samplePkt.AudioSampleAccess.Value = true
	err = this.Protocol.SendPacket(samplePkt, int32(streamId))
	if err != nil {
		return err
	}

	statusPkt := packet.NewSrsOnStatusDataPacket()
	statusPkt.Data.Set(global.StatusCode, global.StatusCodeDataStart)
	err = this.Protocol.SendPacket(statusPkt, int32(streamId))
	if err != nil {
		return err
	}

	return nil
}

func (this *SrsRtmpServer) SetWindowAckSize(act_size int32) error {
	pkt := packet.NewSrsSetWindowAckSizePacket()
	pkt.AckowledgementWindowSize = act_size
	err := this.Protocol.SendPacket(pkt, 0)
	if err != nil {
		return err
	}
	return nil
}

func (this *SrsRtmpServer) SetChunkSize(chunkSize uint32) error {
	pkt := packet.NewSrsSetChunkSizePacket()
	pkt.ChunkSize = int32(chunkSize)
	err := this.Protocol.SendPacket(pkt, 0)
	return err
}

func (this *SrsRtmpServer) SetPeerBandwidth(bandwidth int, typ int8) error {
	pkt := packet.NewSrsSetPeerBandwidthPacket()
	pkt.Bandwidth = int32(bandwidth)
	pkt.Type = typ
	err := this.Protocol.SendPacket(pkt, 0)
	return err
}

func (this *SrsRtmpServer) ResponseConnectApp(objectEncoding float64) error {
	pkt := packet.NewSrsConnectAppResPacket()
	_ = pkt
	pkt.Props.Set("fmsVer", "FMS/"+global.RTMP_SIG_FMS_VER)
	pkt.Props.Set("capabilities", float64(127))
	pkt.Props.Set("mode", float64(1))
	pkt.Info.Set("level", global.StatusLevelStatus)
	pkt.Info.Set("code", global.StatusCodeConnectSuccess)
	pkt.Info.Set("description", "Connection succeeded")
	pkt.Info.Set("objectEncoding", objectEncoding)

	data := amf0.NewSrsAmf0EcmaArray()
	data.Set("version", global.RTMP_SIG_FMS_VER)
    data.Set("srs_sig", global.RTMP_SIG_SRS_KEY)
    data.Set("srs_server", global.RTMP_SIG_SRS_SERVER)
    data.Set("srs_license", global.RTMP_SIG_SRS_LICENSE)
    data.Set("srs_role", global.RTMP_SIG_SRS_ROLE)
    data.Set("srs_url", global.RTMP_SIG_SRS_URL)
    data.Set("srs_version", global.RTMP_SIG_SRS_VERSION)
    data.Set("srs_site", global.RTMP_SIG_SRS_WEB)
    data.Set("srs_email", global.RTMP_SIG_SRS_EMAIL)
    data.Set("srs_copyright", global.RTMP_SIG_SRS_COPYRIGHT)
    data.Set("srs_primary", global.RTMP_SIG_SRS_PRIMARY)
	data.Set("srs_authors", global.RTMP_SIG_SRS_AUTHROS)
	data.Set("srs_pid", float64(12345));
    data.Set("srs_id", float64(12345));
	pkt.Info.Set("data", data)
	
	err := this.Protocol.SendPacket(pkt, 0)
	return err
}

func (this *SrsRtmpServer) OnBwDone() error {
	pkt := packet.NewSrsOnBwDonePacket()
	err := this.Protocol.SendPacket(pkt, 0)
	return err
}

func (this *SrsRtmpServer) StartFmlePublish(streamId int) error {
	// FCPublish
	var fc_publish_tid float64 = 0
	{
		startPkt := packet.NewSrsFMLEStartPacket("")
		if err := this.Protocol.ExpectMessage(startPkt); err != nil {
			return err
		}
		fc_publish_tid = startPkt.TransactionId.GetValue().(float64)
		startResPkt := packet.NewSrsFMLEStartResPacket(fc_publish_tid)
		err := this.Protocol.SendPacket(startResPkt, 0)
		if err != nil {
			// log.Print("send start fmle start res packet failed")
			return err
		}
	}

	var create_stream_tid float64 = 0
	{
		createPkt := packet.NewSrsCreateStreamPacket()
		if err := this.Protocol.ExpectMessage(createPkt); err != nil {
			return err
		}
		create_stream_tid = createPkt.TransactionId.Value
		createResPkt := packet.NewSrsCreateStreamResPacket(create_stream_tid, float64(streamId))
		err := this.Protocol.SendPacket(createResPkt, 0)
		if err != nil {
			return err
		}
	}

	// publish
	{
		publishPacket := packet.NewSrsPublishPacket()
		if err := this.Protocol.ExpectMessage(publishPacket); err != nil {
			return err
		}
	}

	// publish response onFCPublish(NetStream.Publish.Start)
	{
		statusPacket := packet.NewSrsOnStatusCallPacket()
		statusPacket.CommandName.Value.Value = global.RTMP_AMF0_COMMAND_ON_FC_PUBLISH
		statusPacket.Data.Set(global.StatusCode, global.StatusCodePublishStart)
		statusPacket.Data.Set(global.StatusDescription, "Started publishing stream.")
		err := this.Protocol.SendPacket(statusPacket, 0)
		if err != nil {
			log.Print("response onFCPublish failed")
			return err
		}
	}

	{
		statusPacket := packet.NewSrsOnStatusCallPacket()
		statusPacket.Data.Set(global.StatusLevel, global.StatusLevelStatus)
		statusPacket.Data.Set(global.StatusCode, global.StatusCodePublishStart)
		statusPacket.Data.Set(global.StatusDescription, "Started publishing stream.")
		statusPacket.Data.Set(global.StatusClientId, global.RTMP_SIG_CLIENT_ID)
		err := this.Protocol.SendPacket(statusPacket, 0)
		if err != nil {
			log.Print("response onFCPublish failed")
			return err
		}
	}
	log.Print("Start_fmle_publish succeed")

	return nil
}
