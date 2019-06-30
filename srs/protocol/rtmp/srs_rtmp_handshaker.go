package rtmp

import (
	"errors"
	"go_srs/srs/protocol/skt"
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
		return errors.New("only support rtmp plain text.")
	}

	if err = this.HSBytes.CreateS0S1S2(); err != nil {
		return err
	}

	n, err := this.io.Write(this.HSBytes.S0S1S2)
	if err != nil {
		return err
	}

	if 0 != this.HSBytes.ReadC2() {
		return errors.New("HandShake ReadC2 failed")
	}

	if !this.HSBytes.CheckC2() {
		return errors.New("HandShake CheckC2 failed")
	}

	_ = n
	return nil
}

func (this *SrsSimpleHandShake) HandShakeWithServer() error {
	return nil
}
