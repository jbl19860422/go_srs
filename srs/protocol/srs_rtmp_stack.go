package protocol

import (
	_ "bufio"
	"time"
	"net"
	"bytes"
	// log "github.com/sirupsen/logrus"
	"math/rand"
	"encoding/binary"
	"log"
	// "fmt"
)

type SrsHandshakeBytes struct {
	C0C1 []byte
	S0S1S2 []byte
	C2 []byte
	conn *net.Conn
}

func NewSrsHandshakeBytes(c *net.Conn) *SrsHandshakeBytes {
	return &SrsHandshakeBytes{
		conn:c,
	}
}

func (this *SrsHandshakeBytes) ReadC0C1() int {
	var ret int = 0
	if len(this.C0C1) > 0 {
		return ret
	}

	this.C0C1 = make([]byte, 1537, 1537)
	left := 1537
	for {
		n, err := (*this.conn).Read(this.C0C1[1537-left:1537])
		if err != nil {
			return -1
		}
		
		left = left - n
		if left <= 0 {
			return 0
		}
	}
}

func (this *SrsHandshakeBytes) CreateS0S1S2() int {
	if len(this.S0S1S2) > 0 {
		return -1
	}
	rand.Seed(time.Now().UnixNano())
	this.S0S1S2 = make([]byte, 3073)
	//s0 = version
	this.S0S1S2[0] = 0x3
	//s1 for bytes(timestamp)
	binary.Write(bytes.NewBuffer(this.S0S1S2[1:5]), binary.LittleEndian, time.Now().Unix())
	//s1 rand bytes
	if n, err := rand.Read(this.S0S1S2[9:1537]); err != nil || n != 1528 {
		return -2
	}
	//s2=c1
	copy(this.S0S1S2[1537:], this.C0C1[1:])
	return 0
}

func (this *SrsHandshakeBytes) ReadC2() int {
	if len(this.C2) > 0 {
		return -1
	}

	this.C2 = make([]byte, 1536)
	left := 1536
	for {
		n, err := (*this.conn).Read(this.C2[1536-left:1536])
		if err != nil {
			return -1
		}
		log.Print("read n=", n)
		left = left - n
		if left <= 0 {
			return 0
		}
	}
}

func (this *SrsHandshakeBytes) CheckC2() bool {
	return bytes.Equal(this.C2, this.S0S1S2[1:1537])
}

