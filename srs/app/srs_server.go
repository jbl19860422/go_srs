/*
The MIT License (MIT)

Copyright (c) 2013-2015 GOSRS(gosrs)

Permission is hereby granted, free of charge, to any person obtaining a copy of
this software and associated documentation files (the "Software"), to deal in
the Software without restriction, including without limitation the rights to
use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
the Software, and to permit persons to whom the Software is furnished to do so,
subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/
package app

import (
	"net/http"
	"sync"
	_ "log"
	"net"
	"strconv"
	"go_srs/srs/utils"
	"runtime"
	"time"
	"go_srs/srs/app/config"
	"fmt"
	log "github.com/sirupsen/logrus"
)

type SrsServer struct {
	conns 		[]*SrsRtmpConn
	flvServer 	*SrsHttpStreamServer
	connsMtx	sync.Mutex
}

func NewSrsServer() *SrsServer {
	return &SrsServer{
		conns:make([]*SrsRtmpConn, 0),
		flvServer:NewSrsHttpStreamServer(),
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
			this.conns = append(this.conns[:i], this.conns[i+1:]...)
			break
		}
	}
}

func (this *SrsServer) AddConn(c *SrsRtmpConn) error {
	this.connsMtx.Lock()
	if uint32(len(this.conns) + 1) > config.GetInstance().MaxConnections {
		return fmt.Errorf("exceed the max connections, drop client:clients=%d, max=%d", len(this.conns), config.GetInstance().MaxConnections);
	}
	this.conns = append(this.conns, c)
	this.connsMtx.Unlock()
	return nil
}

func (this *SrsServer) StartProcess(port uint32) error {
	log.Info("starting server...")

	ln, err := net.Listen("tcp", ":"+strconv.Itoa(int(port)))
	if err != nil {
		return err
	}

	go func() {
		http.Handle("/", this.flvServer)
		http.Handle("/hls/", http.StripPrefix("/hls/", http.FileServer(http.Dir("./html"))))
		http.ListenAndServe(":8080", nil)
	}()

	go func() {
		for {
			time.Sleep(time.Second*2)
			runtime.GC()
			utils.TraceMemStats()
		}
	}()

	log.Info("starting server succeed")
	for {
		conn, _ := ln.Accept()
		go this.HandleConnection(conn)
	}

	return nil
}

func (this *SrsServer) HandleConnection(conn net.Conn) {
	rtmpConn := NewSrsRtmpConn(conn, this)
	err := this.AddConn(rtmpConn)
	if err != nil {
		conn.Close()
		return
	}
	err = rtmpConn.ServiceLoop()
	this.RemoveConn(rtmpConn)
}

func (this *SrsServer) OnPublish(s *SrsSource, r *SrsRequest) error {
	return nil
}
	
func (this *SrsServer) OnUnpublish(s *SrsSource, r *SrsRequest) error {
	return nil
}
