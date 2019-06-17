package app

import (
	// log "github.com/sirupsen/logrus"
	// "go_srs/srs/protocol"
	// "fmt"

	_ "log"
	"net"
	"strconv"
)

type SrsServer struct {
	conns []*SrsRtmpConn
}

func (this *SrsServer) StartProcess(port int) error {
	ln, err := net.Listen("tcp", ":"+strconv.Itoa(port))
	if err != nil {
		return err
	}

	for {
		conn, _ := ln.Accept()
		go this.HandleConnection(conn)
	}
	return nil
}

func (this *SrsServer) HandleConnection(conn net.Conn) {
	rtmpConn := NewSrsRtmpConn(conn, this)
	this.conns = append(this.conns, rtmpConn)
	rtmpConn.Start()
}

func (this *SrsServer) OnPublish(s *SrsSource, r *SrsRequest) error {
	return nil
}
	
func (this *SrsServer) OnUnpublish(s *SrsSource, r *SrsRequest) error {
	return nil
}
