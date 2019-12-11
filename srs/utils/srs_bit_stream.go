
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

package utils

import (
	"errors"
	"math"
)
type SrsBitStream struct {
	data    []byte
	currBit uint32
}

func NewSrsBitStream(d []byte) *SrsBitStream {
	return &SrsBitStream{
		data:    d,
		currBit: 0,
	}
}

func (this *SrsBitStream) Empty() bool {
	bytePos := this.currBit / 8
	if bytePos >= uint32(len(this.data)) {
		return true
	}
	return false
}

func (this *SrsBitStream) ReadBit() (int8, error) {
	bytePos := this.currBit / 8
	if bytePos >= uint32(len(this.data)) {
		return 0, errors.New("no enough data")
	}

	bitOff := this.currBit % 8
	this.currBit++
	return int8((this.data[bytePos] >> (7 - bitOff)) & 0x01), nil
}

func (this *SrsBitStream) ReadUEV() (int32, error) {
	if this.Empty() {
		return -1, errors.New("no enougth data")
	}
	// 哥伦布熵编码解码
	// ue(v) in 9.1 Parsing process for Exp-Golomb codes
    // H.264-AVC-ISO_IEC_14496-10-2012.pdf, page 227.
    // Syntax elements coded as ue(v), me(v), or se(v) are Exp-Golomb-coded.
    //      leadingZeroBits = -1;
    //      for( b = 0; !b; leadingZeroBits++ )
    //          b = read_bits( 1 )
    // The variable codeNum is then assigned as follows:
	//      codeNum = (2<<leadingZeroBits) - 1 + read_bits( leadingZeroBits )
	
	var leadingZeroBits int = -1
	var b int8 = 0
	var err error
	
	for b = 0; b == 0 && !this.Empty(); leadingZeroBits++ {
		b, err = this.ReadBit()
		if err != nil {
			return -1, err
		}
	}

	if leadingZeroBits >= 31 {
		return -1, errors.New("")
	}

	var v int32 = 0
	v = (1 << uint(leadingZeroBits)) - 1
	for i := 0; i < leadingZeroBits; i++ {
		b, err = this.ReadBit()
		if err != nil {
			return -1, err
		}
		v += int32(b) << uint(leadingZeroBits - 1 - i)
	}
	return v, nil
}

func (this *SrsBitStream) ReadSEV() (int32, error) {
	codeNum, err := this.ReadUEV()
	if err != nil {
		return 0, err
	}
	// H.264-AVC-ISO_IEC_14496-10-2012.pdf, page 229
	//(−1)k+1 Ceil( k÷2 )
	var v int32 = 0
	v = int32(math.Ceil(float64(codeNum)/2))
	if codeNum%2 == 0 {
		v = (-1)*v
	}
	return v, nil
}