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

package config

type DvrConf struct {
	Enabled         string `json:"enabled"`
	DvrPlan         string `json:"dvr_plan"`
	DvrPath         string `json:"dvr_path"`
	DvrDuration     uint32 `json:"dvr_duration"`
	DvrWaitKeyFrame string `json:"dvr_wait_keyframe"`
	TimerJitter     string `json:"timer_jitter"` //full, zero, off
}

const SRS_CONF_DEFAULT_DVR_PATH = "./html/[app]/[stream].[timestamp].flv"
func (this *DvrConf) initDefault() {
	if this.Enabled == "" {
		this.Enabled = "off"
	}

	if this.DvrPlan == "" {
		this.DvrPlan = "session"
	}

	if this.DvrPath == "" {
		this.DvrPath = SRS_CONF_DEFAULT_DVR_PATH
	}

	if this.DvrDuration == 0 {
		this.DvrDuration = 30
	}

	if this.DvrWaitKeyFrame == "" {
		this.DvrWaitKeyFrame = "on"
	}

	if this.TimerJitter == "" {
		this.TimerJitter = "full"
	}
}