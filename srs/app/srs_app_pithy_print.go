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

package app

import "go_srs/srs/app/config"

type SrsStageInfo struct {
	*config.SrsAppSubscriber
	stage_id            int64
	nb_clients          int64
	age                 int64
	pithy_print_time_ms int64
}

func NewSrsStageInfo(stage_id int64) *SrsStageInfo {
	return &SrsStageInfo{
		stage_id:            stage_id,
		nb_clients:          0,
		age:                 0,
		pithy_print_time_ms: config.GetInstance().GetPithyPrintMs(),
	}
}

func (this *SrsStageInfo) Elapse(diff int64) {

}

func (this *SrsStageInfo) CanPrint() bool {
	return true
}

func (this *SrsStageInfo) OnReloadPithyPrint() {

}

type SrsPithyPrint struct {
	client_id     int64
	stage_id      int64
	age           int64
	previout_tick int64
}
