/*
The MIT License (MIT)

Copyright (c) 2019 GOSRS(gosrs)

Permission is hereby granted, free of charge, to any person obtaining a copy of
this software and associated documentation files (the "Software"), to deal in
the Software without restriction, including without limitation the rights to
use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
the Software, and to permit persons to whom the Software is furnished to do so,
subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/
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
