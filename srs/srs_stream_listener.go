package srs
import (
	// "fmt"
	"net"
	"strconv"
	"errors"
	"log"
	// log "github.com/sirupsen/logrus"
)

type SrsStreamListener struct {
	Svr *SrsServer
}

func (s *SrsStreamListener) ListenAndAccept(port int) error {
	ln, err := net.Listen("tcp", ":" + strconv.Itoa(port))
	if err != nil {
		return errors.New("list rtmp port 1935 failed")
	}

	for {
		log.Print("^^^^^^^^^^^^^start accept^^^^^^^^^^^^^^^^^^^")
		conn, _ := ln.Accept()
		log.Print("^^^^^^^^^^^^^^^^^^^^^^^^^get a new connection^^^^^^^^^^^^^^^^^^^^^")
		go HandleConnection(s, conn)
	}
}

func HandleConnection(s *SrsStreamListener, conn net.Conn) {
	c := &SrsRtmpConn{Svr:s.Svr, Conn:conn}
	s.Svr.AcceptConnection(c)
}