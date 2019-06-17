package app

import (
	"go_srs/srs/protocol/rtmp"
	"go_srs/srs/protocol/skt"
	"go_srs/srs/protocol/packet"
	"go_srs/srs/utils"
	"net"
	"strings"
	"net/url"
	"log"
	"time"
	"fmt"
	"errors"
)

type SrsRtmpConn struct {
	io   			*skt.SrsIOReadWriter
	rtmp 			*rtmp.SrsRtmpServer
	req				*SrsRequest
	server			*SrsServer
	clientType 		rtmp.SrsRtmpConnType
	publishThread 	*SrsAppPublishRecvThread
}

func NewSrsRtmpConn(conn net.Conn, s *SrsServer) *SrsRtmpConn {
	socketIO := skt.NewSrsIOReadWriter(conn)
	
	return &SrsRtmpConn{
		io: socketIO,
		rtmp:rtmp.NewSrsRtmpServer(socketIO),
		req:NewSrsRequest(),
		server:s,
	}
}

func (this *SrsRtmpConn) Start() {
	err := this.do_cycle()
	_ = err
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
	fmt.Println(this.req.tcUrl)
	log.Print(m)
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
		log.Print("set_window_ack_size failed")
		return err
	}

	err = this.rtmp.SetPeerBandwidth(1000*1000, 2)
	if err != nil {
		log.Print("set_peer_bandwidth failed")
		return err
	}

	this.req.ip = this.io.GetClientIP()

	err = this.rtmp.SetChunkSize(4096)
	if err != nil {
		log.Print("set_chunk_size failed")
		return err
	}

	err = this.rtmp.ResponseConnectApp()
	if err != nil {
		log.Print("response_connect_app error")
		return err
	}

	for {
		this.stream_service_cycle()
		for {
			time.Sleep(time.Second * 1)
		}
	}
	return nil
}

func (this *SrsRtmpConn) stream_service_cycle() error {
	var typ rtmp.SrsRtmpConnType
	this.req.typ, this.req.stream, _ = this.rtmp.IdentifyClient()
	//log.Print("***************identify_client done ,type=", typ);
	var err error
	this.req.schema, this.req.host, this.req.vhost, this.req.app, _, this.req.port, this.req.param, err = utils.SrsDiscoveryTCUrl(this.req.tcUrl)
	if err != nil {
		return errors.New("srs_discovery_tc_url failed")
	}
	fmt.Println("Srs_discovery_tc_url succeed, stream_name=", this.req.stream)
	source, err2 := FetchOrCreate(this.req, this.server)
	if err2 != nil {
		return err2
	}

	this.clientType = this.req.typ


	switch(this.req.typ) {
	case rtmp.SrsRtmpConnFMLEPublish:{
		log.Print("******************start SrsRtmpConnFMLEPublish*******************")
		this.rtmp.Start_fmle_publish(0)
		return this.publishing(source)
	}
	}
	_ = typ
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
		fmt.Println("onaudio")
	}

	if msg.GetHeader().IsVideo() {
		fmt.Println("onvideo")
		//process video
	}
	//todo fix aggregate message
	//todo fix amf0 or amf3 data
	return nil
}
