package srs

import (
	"net"
)

type SrsRtmpConn struct {
	Svr *SrsServer
	Conn net.Conn
}

func (c *SrsRtmpConn) Start() {
	
}

func (c *SrsRtmpConn) doCycle() {

}

