package app

import (
	// log "github.com/sirupsen/logrus"
	// "go_srs/srs/protocol"
	"fmt"
	"sync"
	_ "log"
	"net"
	"strconv"
)

type SrsServer struct {
	conns 		[]*SrsRtmpConn
	connsMtx	sync.Mutex
}

func (this *SrsServer) OnRecvError(err error, c *SrsRtmpConn) {
	this.RemoveConn(c)
}

func (this *SrsServer) RemoveConn(c *SrsRtmpConn) {
	this.connsMtx.Lock()
	defer this.connsMtx.Unlock()
	for i := 0; i < len(this.conns); i++ {
		if this.conns[i] == c {
			fmt.Println("remove conn")
			this.conns = append(this.conns[:i], this.conns[i+1:]...)
			break
		}
	}
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
	err := rtmpConn.Start()
	_ = err
	this.RemoveConn(rtmpConn)
}

func (this *SrsServer) OnPublish(s *SrsSource, r *SrsRequest) error {
	return nil
}
	
func (this *SrsServer) OnUnpublish(s *SrsSource, r *SrsRequest) error {
	return nil
}
