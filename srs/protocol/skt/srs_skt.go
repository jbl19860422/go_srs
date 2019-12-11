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
package skt

import (
	"bufio"
	"io"
	"net"
	"time"
	_ "fmt"
)

type SrsIOErrListener interface {
	OnRecvError(err error)
}

type SrsIOReadWriter struct {
	conn     net.Conn
	IOReader *bufio.Reader
	IOWriter *bufio.Writer
}

func NewSrsIOReadWriter(c net.Conn) *SrsIOReadWriter {
	rw := &SrsIOReadWriter{
		conn:c,
	}
	rw.IOReader = bufio.NewReader(rw.conn)
	rw.IOWriter = bufio.NewWriter(rw.conn)
	return rw
}

func (this *SrsIOReadWriter) GetClientIP() string {
	return this.conn.RemoteAddr().String()
}

func (this *SrsIOReadWriter) Read(b []byte) (int, error) {
	return this.IOReader.Read(b)
}

func (this *SrsIOReadWriter) Close() {
	this.conn.Close()
}

func (this *SrsIOReadWriter) ReadWithTimeout(b []byte, timeoutms uint32) (int, error) {
	this.conn.SetReadDeadline(time.Now().Add(time.Millisecond * time.Duration(timeoutms)))
	return this.IOReader.Read(b)
}

func (this *SrsIOReadWriter) ReadFully(b []byte, timeoutms uint32) (int, error) {
	count := len(b)
	left := count
	for {
		n, err := this.IOReader.Read(b[count-left:count])
		if err != nil {
			return 0, err
		}

		left = left - n
		if left <= 0 {
			return count, nil
		}
	}
}

func (this *SrsIOReadWriter) ReadFullyWithTimeout(b []byte, timeoutms uint32) (int, error) {
	this.conn.SetReadDeadline(time.Now().Add(time.Millisecond * time.Duration(timeoutms)))
	return io.ReadFull(this.conn, b)
}

func (this *SrsIOReadWriter) Write(b []byte) (int, error) {
	n, err := this.IOWriter.Write(b)
	_ = this.IOWriter.Flush()
	return n, err
}

func (this *SrsIOReadWriter) WriteWithTimeout(b []byte, timeoutms uint32) (int, error) {
	this.conn.SetWriteDeadline(time.Now().Add(time.Millisecond * time.Duration(timeoutms)))
	return this.IOWriter.Write(b)
}
