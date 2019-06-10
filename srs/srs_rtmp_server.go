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

	this.request.ip = this.Conn.Conn.RemoteAddr().String()

	time.Sleep(10 * time.Millisecond)
	err = this.set_chunk_size(60000)
	if err != nil {
		log.Print("set_chunk_size failed")
		return -3
	}

	err = this.response_connect_app()
	if err != nil {
		log.Print("response_connect_app error")
	}

	for {
		this.stream_service_cycle()
		for {
			time.Sleep(time.Second * 1)
		}
	}

	
	return 0
}

func (this *SrsRtmpServer) stream_service_cycle() {
	var typ protocol.SrsRtmpConnType
	this.request.typ, this.request.stream, _ = this.identify_client()
	log.Print("***************identify_client done ,type=", typ);
	var err error
	this.request.schema, this.request.host, this.request.vhost, this.request.app, _, this.request.port, this.request.param, err = protocol.Srs_discovery_tc_url(this.request.tcUrl)
	if err != nil {
		log.Print("Srs_discovery_tc_url failed")
		return
	} else {
		log.Print("Srs_discovery_tc_url succeed, stream_name=", this.request.stream)
	}

	switch(this.request.typ) {
	case protocol.SrsRtmpConnFMLEPublish:{
		log.Print("******************start SrsRtmpConnFMLEPublish*******************")
		this.Protocol.Start_fmle_publish(0)
	}
	}
	_ = typ
}

func (this *SrsRtmpServer) identify_client() (protocol.SrsRtmpConnType, string, error) {
	var typ protocol.SrsRtmpConnType
	var streamname string
	for {
		msg, err := this.Protocol.RecvMessage()
		if err != nil {
			log.Print("identify_client err, msg=", err)
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
		// case SrsCreateStreamPacket: {
		// 	log.Print("SrsCreateStreamPacket")
		// }
		case (*protocol.SrsFMLEStartPacket):
			{
				log.Print("SrsFMLEStartPacket streamname=", pkt.(*protocol.SrsFMLEStartPacket).Stream_name)
				typ, streamname, err = this.identify_fmle_publish_client(pkt.(*protocol.SrsFMLEStartPacket))
				if err != nil {
					log.Print("identify_fmle_publish_client reeturn")
					return typ, streamname, nil
				}
				return typ, streamname, nil
			}
			// case SrsPlayPacket:{
			// 	log.Print("SrsPlayPacket")
			// }
		}
		return typ, streamname, nil
	}
	_ = typ
	return typ, streamname, nil
}

func (this *SrsRtmpServer) identify_fmle_publish_client(req *protocol.SrsFMLEStartPacket) (protocol.SrsRtmpConnType, string, error) {
	typ := protocol.SrsRtmpConnType(protocol.SrsRtmpConnFMLEPublish)
	log.Print("")
	pkt := protocol.NewSrsFMLEStartResPacket(req.Transaction_id)
	err := this.Protocol.SendPacket(pkt, 0)
	if err != nil {
		return typ, req.Stream_name, err
	}
	return typ, req.Stream_name, nil
}

func (this *SrsRtmpServer) connect_app() error {
	// ctx, cancel := context.WithCancel(context.Background())
	connPacket := protocol.NewSrsConnectAppPacket()
	packet := this.Protocol.ExpectMessage(connPacket)
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
	if len(p) >= 2 {
		this.request.app = p[1]
	}

	if len(p) >= 3 {
		this.request.stream = p[2]
	}

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

func (this *SrsRtmpServer) set_chunk_size(chunk_size int32) error {
	pkt := protocol.NewSrsSetChunkSizePacket()
	pkt.Chunk_size = chunk_size
	err := this.Protocol.SendPacket(pkt, 0)
	return err
}

func (this *SrsRtmpServer) set_peer_bandwidth(bandwidth int, typ int8) error {
	pkt := protocol.NewSrsSetPeerBandwidthPacket()
	pkt.Bandwidth = int32(bandwidth)
	pkt.Typ = typ
	err := this.Protocol.SendPacket(pkt, 0)
	return err
}

func (this *SrsRtmpServer) response_connect_app() error {
	pkt := protocol.NewSrsConnectAppResPacket()
	_ = pkt
	pkt.Props.SetStringProperty("fmsVer", "FMS/3,5,3,888")
	pkt.Props.SetNumberProperty("capabilities", float64(127))
	pkt.Props.SetNumberProperty("mode", float64(1))
	pkt.Info.SetStringProperty("level", "status")
	pkt.Info.SetStringProperty("code", "NetConnection.Connect.Success")
	pkt.Info.SetStringProperty("description", "Connection succeeded")
	pkt.Info.SetNumberProperty("objectEncoding", float64(0))

	data := protocol.NewSrsAmf0EcmaArray()
	data.Set("version", protocol.RTMP_SIG_FMS_VER)
    data.Set("srs_sig", protocol.RTMP_SIG_SRS_KEY)
    data.Set("srs_server", protocol.RTMP_SIG_SRS_SERVER)
    data.Set("srs_license", protocol.RTMP_SIG_SRS_LICENSE)
    data.Set("srs_role", protocol.RTMP_SIG_SRS_ROLE)
    data.Set("srs_url", protocol.RTMP_SIG_SRS_URL)
    data.Set("srs_version", protocol.RTMP_SIG_SRS_VERSION)
    data.Set("srs_site", protocol.RTMP_SIG_SRS_WEB)
    data.Set("srs_email", protocol.RTMP_SIG_SRS_EMAIL)
    data.Set("srs_copyright", protocol.RTMP_SIG_SRS_COPYRIGHT)
    data.Set("srs_primary", protocol.RTMP_SIG_SRS_PRIMARY)
	data.Set("srs_authors", protocol.RTMP_SIG_SRS_AUTHROS)
	data.Set("srs_server_ip", "172.19.5.107")
	data.Set("srs_pid", float64(12345));
    data.Set("srs_id", float64(12345));
	pkt.Info.Set("data", data)
	

	err := this.Protocol.SendPacket(pkt, 0)
	return err
	// pkt->props->set("fmsVer", SrsAmf0Any::str("FMS/"RTMP_SIG_FMS_VER));
	// pkt->props->set("capabilities", SrsAmf0Any::number(127));
	// pkt->props->set("mode", SrsAmf0Any::number(1));

	// pkt->info->set(StatusLevel, SrsAmf0Any::str(StatusLevelStatus));
	// pkt->info->set(StatusCode, SrsAmf0Any::str(StatusCodeConnectSuccess));
	// pkt->info->set(StatusDescription, SrsAmf0Any::str("Connection succeeded"));
	// pkt->info->set("objectEncoding", SrsAmf0Any::number(req->objectEncoding));
	// SrsAmf0EcmaArray* data = SrsAmf0Any::ecma_array();
	// pkt->info->set("data", data);

	// data->set("version", SrsAmf0Any::str(RTMP_SIG_FMS_VER));
	// data->set("srs_sig", SrsAmf0Any::str(RTMP_SIG_SRS_KEY));
	// data->set("srs_server", SrsAmf0Any::str(RTMP_SIG_SRS_SERVER));
	// data->set("srs_license", SrsAmf0Any::str(RTMP_SIG_SRS_LICENSE));
	// data->set("srs_role", SrsAmf0Any::str(RTMP_SIG_SRS_ROLE));
	// data->set("srs_url", SrsAmf0Any::str(RTMP_SIG_SRS_URL));
	// data->set("srs_version", SrsAmf0Any::str(RTMP_SIG_SRS_VERSION));
	// data->set("srs_site", SrsAmf0Any::str(RTMP_SIG_SRS_WEB));
	// data->set("srs_email", SrsAmf0Any::str(RTMP_SIG_SRS_EMAIL));
	// data->set("srs_copyright", SrsAmf0Any::str(RTMP_SIG_SRS_COPYRIGHT));
	// data->set("srs_primary", SrsAmf0Any::str(RTMP_SIG_SRS_PRIMARY));
	// data->set("srs_authors", SrsAmf0Any::str(RTMP_SIG_SRS_AUTHROS));

	// if (server_ip) {
	//     data->set("srs_server_ip", SrsAmf0Any::str(server_ip));
	// }
	// // for edge to directly get the id of client.
	// data->set("srs_pid", SrsAmf0Any::number(getpid()));
	// data->set("srs_id", SrsAmf0Any::number(_srs_context->get_id()));

	// if ((ret = protocol->send_and_free_packet(pkt, 0)) != ERROR_SUCCESS) {
	//     srs_error("send connect app response message failed. ret=%d", ret);
	//     return ret;
	// }
	// srs_info("send connect app response message success.");

}
