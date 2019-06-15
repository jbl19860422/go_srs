package handshake

type HandShaker interface {
	HandShakeWithClient() error
	HandShakeWithServer() error
	conn *net.Conn
}

type SrsSimpleHandShake struct {
	HSBytes *SrsHandshakeBytes
}

func NewSrsSimpleHandShake(c *net.Conn) *SrsSimpleHandShake {
	return &SrsSimpleHandShake{
		HSBytes:NewSrsHandshakeBytes(c)
	}
}

func (this *SrsSimpleHandShake) HandShakeWithClient() error {
	ret := this.HSBytes.ReadC0C1()
	if 0 != ret {
		log.Printf("HandShake ReadC0C1 failed")
		return -1
	}

	if this.HandShaker.C0C1[0] != 0x03 {
		log.Printf("only support rtmp plain text.")
		return -2
	}

	if 0 != this.HandShaker.CreateS0S1S2() {
		return -2
	}

	n, err := this.Conn.Conn.Write(this.HandShaker.S0S1S2)
	if err != nil {
		log.Printf("write s0s1s2 failed")
	} else {
		log.Printf("write s0s1s2 succeed, count=", len(this.HandShaker.S0S1S2))
	}

	if 0 != this.HandShaker.ReadC2() {
		log.Printf("HandShake ReadC2 failed")
		return -3
	}

	if !this.HandShaker.CheckC2() {
		log.Printf("HandShake CheckC2 failed")
	}

	log.Printf("HandShake Succeed")
	_ = n

	return 0
}

func (this *SrsSimpleHandShake) HandShakeWithServer() error {

}