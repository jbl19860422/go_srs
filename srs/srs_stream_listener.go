package srs
import (
	// "fmt"
	"net"
	"strconv"
	"errors"
	// log "github.com/sirupsen/logrus"
)

type SrsStreamListener struct {
	Svr *SrsServer
}

func (s *SrsStreamListener) ListenAndAccept() error {
	ln, err := net.Listen("tcp", ":" + strconv.Itoa(RTMP_PORT))
	if err != nil {
		return errors.New("list rtmp port 1935 failed")
	}

	for {
		conn, _ := ln.Accept()
		go s.HandleConnection(conn)
	}
}

func (s *SrsStreamListener) HandleConnection(conn net.Conn) {
	c := &SrsRtmpConn{Svr:s.Svr, Conn:conn}
	s.Svr.AcceptConnection(c)
}