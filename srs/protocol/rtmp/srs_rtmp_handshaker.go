package rtmp

import (
	"errors"
	"go_srs/srs/protocol/skt"
	"log"
	// "fmt"
)

type HandShaker interface {
	HandShakeWithClient() error
	HandShakeWithServer() error
}

type SrsSimpleHandShake struct {
	HSBytes *SrsHandshakeBytes
	io      *skt.SrsIOReadWriter
}

func NewSrsSimpleHandShake(io_ *skt.SrsIOReadWriter) *SrsSimpleHandShake {
	return &SrsSimpleHandShake{
		HSBytes: NewSrsHandshakeBytes(io_),
		io:      io_,
	}
}

func (this *SrsSimpleHandShake) HandShakeWithClient() error {
	err := this.HSBytes.ReadC0C1()
	if err != nil {
		return err
	}

	if this.HSBytes.C0C1[0] != 0x03 {
		log.Printf("only support rtmp plain text.")
		return errors.New("only support rtmp plain text.")
	}

	if err = this.HSBytes.CreateS0S1S2(); err != nil {
		return err
	}

	// fmt.Println("this.HSBytes.S0S1S2=", len(this.HSBytes.S0S1S2))
	n, err := this.io.Write(this.HSBytes.S0S1S2)
	if err != nil {
		log.Printf("write s0s1s2 failed")
	} else {
		log.Printf("write s0s1s2 succeed, count=", len(this.HSBytes.S0S1S2))
	}

	if 0 != this.HSBytes.ReadC2() {
		log.Printf("HandShake ReadC2 failed")
		return errors.New("HandShake ReadC2 failed")
	}

	if !this.HSBytes.CheckC2() {
		log.Printf("HandShake CheckC2 failed")
	}

	log.Printf("HandShake Succeed")
	_ = n

	return nil
}

func (this *SrsSimpleHandShake) HandShakeWithServer() error {
	return nil
}
