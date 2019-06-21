package app

import (
	// log "github.com/sirupsen/logrus"
	// "go_srs/srs/protocol"
	"fmt"
	"sync"
	_ "log"
	"net"
	"strconv"
	"go_srs/srs/utils"
	"runtime"
	"time"
)

type SrsServer struct {
	conns 		[]*SrsRtmpConn
	connsMtx	sync.Mutex
}

func NewSrsServer() *SrsServer {
	return &SrsServer{
		conns:make([]*SrsRtmpConn, 0),
	}
}

func (this *SrsServer) OnRecvError(err error, c *SrsRtmpConn) {
	this.RemoveConn(c)
}

func (this *SrsServer) RemoveConn(c *SrsRtmpConn) {
	this.connsMtx.Lock()
	defer this.connsMtx.Unlock()
	for i := 0; i < len(this.conns); i++ {
		if this.conns[i] == c {
			fmt.Println("remove conn len=", len(this.conns))
			this.conns = append(this.conns[:i], this.conns[i+1:]...)
			fmt.Println("conns.len=", len(this.conns))
			break
		}
	}
}

func (this *SrsServer) AddConn(c *SrsRtmpConn) {
	this.connsMtx.Lock()
	this.conns = append(this.conns, c)
	fmt.Println("xxxxxxxxxxxxxxxxxxxconns.len=", len(this.conns), "xxxxxxxxxxxxxxxxxxxxxxxxxxxx")
	this.connsMtx.Unlock()
}

func (this *SrsServer) StartProcess(port int) error {
	ln, err := net.Listen("tcp", ":"+strconv.Itoa(port))
	if err != nil {
		return err
	}

	go func() {
		for {
			time.Sleep(time.Second*2)
			runtime.GC()
			utils.TraceMemStats()
		}
	}()

	for {
		conn, _ := ln.Accept()
		go this.HandleConnection(conn)
	}
	return nil
}

func (this *SrsServer) HandleConnection(conn net.Conn) {
	rtmpConn := NewSrsRtmpConn(conn, this)
	this.AddConn(rtmpConn)
	err := rtmpConn.Start()
	_ = err
	this.RemoveConn(rtmpConn)

	fmt.Println("HandleConnection done")
}

func (this *SrsServer) OnPublish(s *SrsSource, r *SrsRequest) error {
	return nil
}
	
func (this *SrsServer) OnUnpublish(s *SrsSource, r *SrsRequest) error {
	return nil
}
