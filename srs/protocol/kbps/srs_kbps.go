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

package kbps

import "go_srs/srs/utils"

type SrsKbpsSample struct {
	bytes int64
	time  int64
	kbps  int64
}

type SrsKbpsSlice struct {
	io ISrsIOStatistic
	// session startup bytes
	// @remark, use total_bytes() to get the total bytes of slice.
	bytes         int64
	starttime     int64
	io_bytes_base int64
	last_bytes    int64

	sample_30s SrsKbpsSample
	sample_1m  SrsKbpsSample
	sample_5m  SrsKbpsSample
	sample_60m SrsKbpsSample

	delta_bytes int64
}

func NewSrsKbpsSlice() *SrsKbpsSlice {
	return &SrsKbpsSlice{}
}

func (this SrsKbpsSlice) GetTotalBytes() int64 {
	return this.bytes + this.last_bytes - this.io_bytes_base
}

func (this *SrsKbpsSlice) Sample() {
	now := utils.SystemTimeMs()
	total_bytes := this.GetTotalBytes()
	// if the sample is not initialized, initialize first.
	if this.sample_30s.time <= 0 {
		this.sample_30s.kbps = 0
		this.sample_30s.time = now
		this.sample_30s.bytes = total_bytes
	}

	if this.sample_1m.time <= 0 {
		this.sample_1m.kbps = 0
		this.sample_1m.time = now
		this.sample_1m.bytes = total_bytes
	}

	if this.sample_5m.time <= 0 {
		this.sample_5m.kbps = 0
		this.sample_5m.time = now
		this.sample_5m.bytes = 0
	}

	if this.sample_60m.time <= 0 {
		this.sample_60m.kbps = 0
		this.sample_60m.time = now
		this.sample_60m.bytes = 0
	}
	//caculate the result
	if now-this.sample_30s.time >= 30*1000 {
		this.sample_30s.kbps = int64((total_bytes - this.sample_30s.bytes) * 8 / (now - this.sample_30s.time))
		this.sample_30s.time = now
		this.sample_30s.bytes = total_bytes
	}

	if now-this.sample_1m.time >= 60*1000 {
		this.sample_1m.kbps = int64((total_bytes - this.sample_1m.bytes) * 8 / (now - this.sample_1m.time))
		this.sample_1m.time = now
		this.sample_1m.bytes = total_bytes
	}

	if now-this.sample_5m.time >= 5*60*1000 {
		this.sample_5m.kbps = int64((total_bytes - this.sample_5m.bytes) * 8 / (now - this.sample_5m.time))
		this.sample_5m.time = now
		this.sample_5m.bytes = total_bytes
	}

	if now-this.sample_60m.time >= 60*60*1000 {
		this.sample_60m.kbps = int64((total_bytes - this.sample_60m.bytes) * 8 / (now - this.sample_60m.time))
		this.sample_60m.time = now
		this.sample_60m.bytes = total_bytes
	}
}

type SrsKbps struct {
	is SrsKbpsSlice
	os SrsKbpsSlice
}

func NewSrsKbps() *SrsKbps {
	return &SrsKbps{}
}

func (this *SrsKbps) SetIO(in ISrsIOStatistic, out ISrsIOStatistic) error {
	if this.is.starttime == 0 {
		this.is.starttime = utils.GetCurrentMs()
	}
	if this.is.io != nil {
		this.is.bytes = this.is.GetTotalBytes() - this.is.io_bytes_base
	}
	this.is.io = in
	this.is.last_bytes = 0
	this.is.io_bytes_base = 0
	if in != nil {
		this.is.last_bytes = in.GetRecvBytes()
		this.is.io_bytes_base = in.GetRecvBytes()
	}

	this.is.Sample()

	if this.os.starttime == 0 {
		this.os.starttime = utils.GetCurrentMs()
	}

	if this.os.io != nil {
		this.os.bytes = this.os.GetTotalBytes() - this.os.io_bytes_base
	}

	this.os.io = out
	this.os.last_bytes = 0
	this.os.io_bytes_base = 0
	if out != nil {
		this.os.last_bytes = out.GetRecvBytes()
		this.os.io_bytes_base = out.GetRecvBytes()
	}

	this.os.Sample()
	return nil
}

func (this *SrsKbps) GetSendKbps() int64 {
	duration := utils.GetCurrentMs() - this.is.starttime
	if duration <= 0 {
		return 0
	}
	bytes := this.GetSendBytes() - this.os.io_bytes_base
	return (bytes*8)/duration
}

func (this *SrsKbps) GetSendBytes() int64 {
	b := this.os.bytes
	if this.os.io != nil {
		b += this.os.io.GetSendBytes() - this.os.io_bytes_base
	}
	b += this.os.last_bytes - this.os.io_bytes_base
	return b
}

func (this *SrsKbps) GetRecvKbps() int64 {
	duration := utils.GetCurrentMs() - this.is.starttime
	if duration <= 0 {
		return 0
	}
	bytes := this.GetRecvBytes() - this.is.io_bytes_base
	return (bytes*8)/duration
}

func (this *SrsKbps) GetRecvBytes() int64 {
	b := this.is.bytes
	if this.is.io != nil {
		b += this.is.io.GetRecvBytes() - this.is.io_bytes_base
	}
	b += this.is.last_bytes - this.is.io_bytes_base
	return b
}

func (this *SrsKbps) Resample() {
	this.sample()
}

func (this *SrsKbps) sample() {
	if this.os.io != nil {
		this.os.last_bytes = this.os.io.GetSendBytes()
	}

	if this.is.io != nil {
		this.is.last_bytes = this.is.io.GetRecvBytes()
	}

	this.os.Sample()
	this.is.Sample()
}

func (this *SrsKbps) GetSendKbps30s() int64 {
	return this.os.sample_30s.kbps
}

func (this *SrsKbps) GetRecvKbps30s() int64 {
	return this.is.sample_30s.kbps
}

func (this *SrsKbps) GetSendKbps1m() int64 {
	return this.os.sample_1m.kbps
}

func (this *SrsKbps) GetRecvKbps1m() int64 {
	return this.is.sample_1m.kbps
}

func (this *SrsKbps) GetSendKbps5m() int64 {
	return this.os.sample_5m.kbps
}

func (this *SrsKbps) GetRecvKbps5m() int64 {
	return this.is.sample_5m.kbps
}

func (this *SrsKbps) GetSendKbps60m() int64 {
	return this.os.sample_60m.kbps
}

func (this *SrsKbps) GetRecvKbps60m() int64 {
	return this.is.sample_60m.kbps
}
