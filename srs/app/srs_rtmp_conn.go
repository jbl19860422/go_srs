package app

import (
	"go_srs/srs/protocol/rtmp"
	"go_srs/srs/protocol/skt"
	"go_srs/srs/protocol/packet"
	"go_srs/srs/codec/flv"
	"go_srs/srs/utils"
	"net"
	"strings"
	"net/url"
	// "log"
	"time"
	"fmt"
	// "context"
	"errors"
)

type SrsRtmpConn struct {
	io   					*skt.SrsIOReadWriter
	rtmp 					*rtmp.SrsRtmpServer
	req						*SrsRequest
	res 					*SrsResponse
	server				*SrsServer
	source				*SrsSource
	clientType 		rtmp.SrsRtmpConnType
	publishThread *SrsAppPublishRecvThread			
}

func NewSrsRtmpConn(conn net.Conn, s *SrsServer) *SrsRtmpConn {
	socketIO := skt.NewSrsIOReadWriter(conn)
	// ctx, cancelFun := context.WithCancel(context.Background())
	rtmpConn := &SrsRtmpConn{
		io: socketIO,
		req:NewSrsRequest(),
		res:NewSrsResponse(1),
		server:s,
	}
	rtmpConn.rtmp = rtmp.NewSrsRtmpServer(socketIO, rtmpConn)
	return rtmpConn
}

func (this *SrsRtmpConn) Start() error {
	return this.do_cycle()
}

func (this *SrsRtmpConn) Stop() {
	// this.stopFun()
	this.io.Close()
	fmt.Println("typ=", this.req.typ)
	if this.req.typ == rtmp.SrsRtmpConnFMLEPublish || this.req.typ == rtmp.SrsRtmpConnFlashPublish || this.req.typ == rtmp.SrsRtmpConnHaivisionPublish {
		fmt.Println("RemoveConsumers")
		this.source.RemoveConsumers()
		RemoveSrsSource(this.source)
	}

	fmt.Println("remove source conn")
}

func (this *SrsRtmpConn) do_cycle() error {
	if err := this.rtmp.HandShake(); err != nil {
		return err
	}

	pkt, err := this.rtmp.ConnectApp()
	if err != nil {
		return err
	}
	
	err = pkt.(*packet.SrsConnectAppPacket).CommandObj.Get("tcUrl", &this.req.tcUrl)
	if err != nil {
		return err
	}

	err = pkt.(*packet.SrsConnectAppPacket).CommandObj.Get("tcUrl", &this.req.pageUrl)
	if err != nil {
		return err
	}

	err = pkt.(*packet.SrsConnectAppPacket).CommandObj.Get("tcUrl", &this.req.swfUrl)
	if err != nil {
		return err
	}

	u, err := url.Parse(this.req.tcUrl)
	this.req.schema = u.Scheme
	this.req.host = u.Host
	p := strings.Split(u.Path, "/")
	if len(p) >= 2 {
		this.req.app = p[1]
	}

	if len(p) >= 3 {
		this.req.stream = p[2]
	}

	m, _ := url.ParseQuery(u.RawQuery)
	// fmt.Println(this.req.tcUrl)
	// log.Print(m)
	vhost, ok := m["vhost"]
	if ok {
		this.req.vhost = vhost[0]
	}

	this.service_cycle()
	return nil
}

func (this *SrsRtmpConn) service_cycle() error {
	err := this.rtmp.SetWindowAckSize((int32)(1000000))
	if err != nil {
		// log.Print("set_window_ack_size failed")
		return err
	}

	err = this.rtmp.SetPeerBandwidth(1000*1000, 2)
	if err != nil {
		// log.Print("set_peer_bandwidth failed")
		return err
	}

	this.req.ip = this.io.GetClientIP()

	err = this.rtmp.SetChunkSize(4096)
	if err != nil {
		// log.Print("set_chunk_size failed")
		return err
	}

	err = this.rtmp.ResponseConnectApp()
	if err != nil {
		// log.Print("response_connect_app error")
		return err
	}

	return this.stream_service_cycle()
}

func (this *SrsRtmpConn) stream_service_cycle() error {
	var dur float64
	this.req.typ, this.req.stream, dur, _ = this.rtmp.IdentifyClient(this.res.StreamId)
	_ = dur
	var err error
	this.req.schema, this.req.host, this.req.vhost, this.req.app, _, this.req.port, this.req.param, err = utils.SrsDiscoveryTCUrl(this.req.tcUrl)
	if err != nil {
		return errors.New("srs_discovery_tc_url failed")
	}
	// fmt.Println("Srs_discovery_tc_url succeed, stream_name=", this.req.stream)
	this.source, err = FetchOrCreate(this.req, this.server)
	if err != nil {
		fmt.Println("FetchOrCreate failed")
		return err
	}

	this.clientType = this.req.typ
	// fmt.Println("*************clientType=", this.clientType, "*************")

	switch(this.req.typ) {
	case rtmp.SrsRtmpConnPlay:{
		if err := this.rtmp.StartPlay(this.res.StreamId); err != nil {
			return err
		}

		//todo http_hooks_on_play

		return this.playing(this.source)
	}
	case rtmp.SrsRtmpConnFMLEPublish:{
		// log.Print("******************start SrsRtmpConnFMLEPublish*******************")
		this.rtmp.Start_fmle_publish(0)
		return this.publishing(this.source)
	}
	}
	return nil
}

func (this *SrsRtmpConn) playing(source *SrsSource) error {
	consumer := source.CreateConsumer(this, true, true, true)
	return this.do_playing(source, consumer)
}

func (this *SrsRtmpConn) RemoveSelf() {
	this.server.RemoveConn(this)
}

func (this *SrsRtmpConn) OnRecvError(err error) {
	//判断如果是publish，则删除源
	this.Stop()
	this.server.OnRecvError(err, this)
}

func (this *SrsRtmpConn) do_playing(source *SrsSource, consumer *SrsConsumer) error {
	//todo refer check
	//todo srsprint
	realtime := false

	for {
		// fmt.Println("*************do_playing start***************")
		//todo expired
		for !consumer.queueRecvThread.Empty() {//process signal message
			msg := consumer.queueRecvThread.GetMsg()
			if msg != nil {
				err := this.process_play_control_msg(consumer, msg)
				if err != nil {
					return err
				}
			}
		}
		//todo process trd error
		//todo process realtime stream
		if realtime {

		} else {
			msg, err := consumer.Wait(1, 100)
			if err != nil {
				return err
			}

			if msg != nil {
				// fmt.Println("send to consumer")
				if msg.GetHeader().IsVideo() {
					//fmt.Println("xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxsendmsg video");
					if flvcodec.VideoIsKeyframe(msg.GetPayload()) {
						// fmt.Println("send key frame")
					}
				} else {
					//fmt.Println("xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxsendmsg audio");
				}
				
				err := this.rtmp.SendMsg(msg, this.res.StreamId)
				_ = err
			}
		}

		//time.Sleep(time.Millisecond*1)
	}

	return nil
}

func (this *SrsRtmpConn) process_play_control_msg(consumer *SrsConsumer, msg *rtmp.SrsRtmpMessage) error {
	if !msg.GetHeader().IsAmf0Command() && !msg.GetHeader().IsAmf3Command() {
		//ignore 
		return nil
	}
	
	pkt, err := this.rtmp.DecodeMessage(msg)
	if err != nil {
		return err
	}
	//todo add callpacket 
	//todo process pause message
	switch pkt.(type) {
	case *packet.SrsCloseStreamPacket:{
		//todo fix close stream action
		return errors.New("get close stream packet")
	}
	case *packet.SrsPausePacket:{
		return nil
	}
	}
	return nil
}

func (this *SrsRtmpConn) publishing(s *SrsSource) error {
	//TODO
	//refer.check
	//http_hooks_on_publish
	//judge edge host
	if err := this.acquirePublish(s, false); err != nil {
		return err
	}

	err := this.doPublishing(s)
	return err
}

func (this *SrsRtmpConn) acquirePublish(source *SrsSource, isEdge bool) error {
	//TODO edge process
	return nil
}

func (this *SrsRtmpConn) doPublishing(source *SrsSource) error {
	// fmt.Println("******************doPublishing*******************")
	this.publishThread = NewSrsAppPublishRecvThread(this.rtmp, this.req, this, source, false, false)
	this.publishThread.Start()
	for {
		time.Sleep(time.Second)
	}
	return nil
}

func (this *SrsRtmpConn) HandlePublishMessage(source *SrsSource, msg *rtmp.SrsRtmpMessage, isFmle bool, isEdge bool) error {
	if msg.GetHeader().IsAmf0Command() || msg.GetHeader().IsAmf3Command() {
		pkt, err := this.rtmp.DecodeMessage(msg)
		if err != nil {
			return err
		}
		_ = pkt
		//todo isfmle process
	}

	return this.ProcessPublishMessage(source, msg, isEdge)
}

func (this *SrsRtmpConn) ProcessPublishMessage(source *SrsSource, msg *rtmp.SrsRtmpMessage, isEdge bool) error {
	//todo fix edge process
	if msg.GetHeader().IsAudio() {
		//process audio
		// fmt.Println("onaudio*******************")
		if err := source.OnAudio(msg); err != nil {

		}
	}

	if msg.GetHeader().IsVideo() {
		// fmt.Println("onvideo******************")
		if err := source.OnVideo(msg); err != nil {
			
		}
		//process video
	}
	//todo fix aggregate message
	//todo fix amf0 or amf3 data

	// process onMetaData
    if (msg.GetHeader().IsAmf0Data() || msg.GetHeader().IsAmf3Data()) {
		pkt, err := this.rtmp.DecodeMessage(msg)
		if err != nil {
			return err
		}

		switch pkt.(type) {
			case *packet.SrsOnMetaDataPacket: {
				// fmt.Println("xxxxxxxxxxxxxxxxxxxxxxxxxxxmetadata")
				err := source.on_meta_data(msg, pkt.(*packet.SrsOnMetaDataPacket))
				if err != nil {
					return err
				}
			}
		}
    }
	return nil
}

func (this *SrsRtmpConn) Playing(source *SrsSource) {
	//todo
}
