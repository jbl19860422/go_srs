package conn

import (
	"net"
)

type SrsRtmpConn struct {
	Svr *SrsServer
	conn net.Conn
	IOReader *bufio.Reader
	IOWriter *bufio.Writer
}

func (c *SrsRtmpConn) Start() {
	
}

func (c *SrsRtmpConn) doCycle() {

}

func (this *SrsRtmpConn) SetRecvTimeout() {

}

func (this *SrsRtmpConn) Read(b []byte) (n int, err error) {
	n, err = this.IOReader.Read(b)
	if err == io.EOF {

	}
}

func (this *SrsRtmpConn) ReadWithTimeout() {

}