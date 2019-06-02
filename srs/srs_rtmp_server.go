package srs

import (
	// "fmt"
	_ "context"
	"go_srs/srs/protocol"
	"log"
	"net/url"
	"strings"
	"time"
	// log "github.com/sirupsen/logrus"
)

type SrsRtmpServer struct {
	Conn       *SrsRtmpConn
	Protocol   *protocol.SrsProtocol
	HandShaker protocol.SrsHandshakeBytes
	request    SrsRequest
}

func NewSrsRtmpServer(conn *SrsRtmpConn) *SrsRtmpServer {
	return &SrsRtmpServer{Conn: conn, Protocol: protocol.NewSrsProtocol(&conn.Conn), HandShaker: protocol.SrsHandshakeBytes{}}
}

func (this *SrsRtmpServer) HandShake() int {
	ret := this.HandShaker.ReadC0C1(&(this.Conn.Conn))
	if 0 != ret {
		log.Printf("HandShake ReadC0C1 failed")
		return -1
	}

	if this.HandShaker.C0C1[0] != 0x03 {
		log.Printf("only support rtmp plain text.")
		return -2
	}

	if 0 != this.HandShaker.CreateS0S1S2() {
		return -2
	}

	n, err := this.Conn.Conn.Write(this.HandShaker.S0S1S2)
	if err != nil {
		log.Printf("write s0s1s2 failed")
	} else {
		log.Printf("write s0s1s2 succeed, count=", len(this.HandShaker.S0S1S2))
	}

	if 0 != this.HandShaker.ReadC2(&(this.Conn.Conn)) {
		log.Printf("HandShake ReadC2 failed")
		return -3
	}

	if !this.HandShaker.CheckC2() {
		log.Printf("HandShake CheckC2 failed")
	}

	log.Printf("HandShake Succeed")
	_ = n

	return 0
}

func (this *SrsRtmpServer) Start() int {
	log.Printf("start rtmp server")
	ret := this.HandShake()
	if ret != 0 {
		log.Printf("HandShake failed")
		return -1
	}

	err := this.connect_app()
	if err != nil {
		log.Print("connect app failed")
		return -2
	}

	ret = this.service_cycle()

	// http_hooks_on_close() //结束回调http
	//this.request.ip = "xxx"

	// _ = msg

	return 0
}

func (this *SrsRtmpServer) service_cycle() int {
	err := this.set_window_ack_size((int32)(1000000))
	if err != nil {
		log.Print("set_window_ack_size failed")
		return -1
	}

	err = this.set_peer_bandwidth(1000*1000, 2)
	if err != nil {
		log.Print("set_peer_bandwidth failed")
		return -2
	}

	for {
		time.Sleep(10*time.Second)
	}
	return 0
}
func (this *SrsRtmpServer) connect_app() error {
	// ctx, cancel := context.WithCancel(context.Background())
	// go this.Protocol.LoopMessage(ctx, &(this.Conn.Conn))
	connPacket := protocol.NewSrsConnectAppPacket()
	packet := this.Protocol.ExpectMessage(&(this.Conn.Conn), connPacket)
	var err error
	this.request.tcUrl, err = packet.(*protocol.SrsConnectAppPacket).CommandObj.GetStringProperty("tcUrl")
	if err != nil {
		return err
	}

	this.request.pageUrl, err = packet.(*protocol.SrsConnectAppPacket).CommandObj.GetStringProperty("tcUrl")
	if err != nil {
		return err
	}

	this.request.swfUrl, err = packet.(*protocol.SrsConnectAppPacket).CommandObj.GetStringProperty("tcUrl")
	if err != nil {
		return err
	}

	u, err := url.Parse(this.request.tcUrl)
	this.request.schema = u.Scheme
	this.request.host = u.Host
	p := strings.Split(u.Path, "/")
	this.request.app = p[1]
	this.request.stream = p[2]
	m, _ := url.ParseQuery(u.RawQuery)
	log.Print("****************************", this.request.schema)
	log.Print("****************************", this.request.host)
	log.Print("****************************", this.request.app)
	log.Print("****************************", u.RawQuery)
	log.Print(m)
	vhost, ok := m["vhost"]
	if ok {
		this.request.vhost = vhost[0]
	}
	log.Print("****************************", this.request.vhost)
	// srs_discovery_tc_url(req->tcUrl,
	//     req->schema, req->host, req->vhost, req->app, req->stream, req->port,
	//     req->param);
	// req->strip();
	// for {
	// 	time.Sleep(10 * time.Second)
	// }
	return err
}

func (this *SrsRtmpServer) set_window_ack_size(act_size int32) error {
	pkt := protocol.NewSrsSetWindowAckSizePacket()
	pkt.Ackowledgement_window_size = act_size
	err := this.Protocol.SendPacket(pkt, 0)
	if err != nil {
		log.Print("send packet err ", err)
		return err
	}
	log.Print("send act size succeed")
	return nil
}

func (this *SrsRtmpServer) set_peer_bandwidth(bandwidth int, typ int8) error {
    pkt := protocol.NewSrsSetPeerBandwidthPacket();
    pkt.Bandwidth = int32(bandwidth)
	pkt.Typ = typ
	err := this.Protocol.SendPacket(pkt, 0)
	return err
}
