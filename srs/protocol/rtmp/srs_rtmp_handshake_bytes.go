/*
The MIT License (MIT)

Copyright (c) 2013-2015 GOSRS(gosrs)

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

import(
	"errors"
	"math/rand"
	"time"
	"encoding/binary"
	"bytes"
	"go_srs/srs/protocol/skt"
	"go_srs/srs/utils"
)

type SrsHandshakeBytes struct {
	C0C1 []byte
	S0S1S2 []byte
	C2 []byte
	io *skt.SrsIOReadWriter
}

func NewSrsHandshakeBytes(io_ *skt.SrsIOReadWriter) *SrsHandshakeBytes {
	return &SrsHandshakeBytes{
		io: io_,
	}
}

func (this *SrsHandshakeBytes) ReadC0C1() error {
	if len(this.C0C1) > 0 {
		err := errors.New("handshake read c0c1 failed, already read")
		return err
	}

	this.C0C1 = make([]byte, 1537)
	left := 1537
	for {
		n, err := this.io.Read(this.C0C1[1537-left:1537])
		if err != nil {
			return err
		}
		
		left = left - n
		if left <= 0 {
			return nil
		}
	}
}

func (this *SrsHandshakeBytes) CreateS0S1S2() error {
	if len(this.S0S1S2) > 0 {
		return errors.New("already create")
	}
	rand.Seed(time.Now().UnixNano())
	this.S0S1S2 = make([]byte, 3073)
	//s0 = version
	this.S0S1S2[0] = 0x3
	//s1 for bytes(timestamp)
	b := utils.Int32ToBytes(int32(time.Now().Unix()), binary.LittleEndian)
	copy(this.S0S1S2[1:5], b)
	// binary.Write(bytes.NewBuffer(this.S0S1S2[1:5]), binary.LittleEndian, )
	//s1 rand bytes
	if n, err := rand.Read(this.S0S1S2[9:1537]); err != nil || n != 1528 {
		return errors.New("create rand number failed")
	}
	//s2=c1
	copy(this.S0S1S2[1537:], this.C0C1[1:])
	return nil
}

func (this *SrsHandshakeBytes) ReadC2() int {
	if len(this.C2) > 0 {
		return -1
	}

	this.C2 = make([]byte, 1536)
	left := 1536
	for {
		n, err := this.io.Read(this.C2[1536-left:1536])
		if err != nil {
			return -1
		}
		left = left - n
		if left <= 0 {
			return 0
		}
	}
}

func (this *SrsHandshakeBytes) CheckC2() bool {
	return bytes.Equal(this.C2, this.S0S1S2[1:1537])
}
