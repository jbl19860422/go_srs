package protocol

import (
	"bufio"
	"time"
	"net"
	log "github.com/sirupsen/logrus"
	// "fmt"
)

type SrsHandshakeBytes struct {
	C0C1 []byte
	S0S1S2 []byte
	C2 []byte
}

func (this *SrsHandshakeBytes) ReadC0C1(c *net.Conn) int {
	var ret int = 0
	if len(this.C0C1) > 0 {
		return ret
	}

	this.C0C1 = make([]byte, 1537, 1537)
	(*c).SetReadDeadline(time.Now().Add(1000*time.Millisecond))
	reader := bufio.NewReader(*c)
	for {
		n, err := reader.Read(this.C0C1)
		if err != nil {
			return -1
		} else {
			log.Info("read bytes len=", n)
		}
		var _ = n
	}
	return 0
}

