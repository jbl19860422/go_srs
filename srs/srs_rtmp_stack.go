package srs

import (
	"bufio"
	// "time"
	log "github.com/sirupsen/logrus"
	// "fmt"
)

type SrsHandshakeBytes struct {
	C0C1 []byte
	S0S1S2 []byte
	C2 []byte
}

func (this *SrsHandshakeBytes) ReadC0C1(c *SrsRtmpConn) int {
	// var ret int = 0
	// if len(this.C0C1) > 0 {
	// 	return ret
	// }

	// this.C0C1 = make([]byte, 1537, 1537)
	// var data []byte = make([]byte, 1024, 1024)
	// _, err := c.Conn.Read(data)
	// if err != nil {
	// 	fmt.Println("read data failed, err=" + err.Error())
	// 	return 0
	// }
	// str := string(data)
	// fmt.Println("server recv:" + str)
	// return 0
	log.Info("start read c0c1", c.Conn)
	reader := bufio.NewReader(c.Conn)
	for {
		n, err := reader.Read(this.C0C1)
		if err != nil {
			log.Error("read c0c1 failed, err=", err.Error())
		} else {
			log.Info("read bytes len=", n)
			// if n != 0 {
			// 	log.Info("read bytes len=", n)
			// } else {

			// }
			// time.Sleep(10*time.Millisecond)
			// log.Info("read bytes len=", n)
			// break;
		}

		var _ = n
	}
	return 0
}

