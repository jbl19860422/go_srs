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

package config

type HeartBeatConf struct {
	Enabled   string  `json:"enabled"`
	Interval  float64 `json:"interval"`
	Url       string  `json:"url"`
	DeviceId  string  `json:"device_id"`
	Summeries string  `json:"summeries"`
}

func (this *HeartBeatConf) amendDefault() {
	if this.Enabled != "on" {
		this.Enabled = "off"
	}

	if this.Interval <= 0 {
		this.Interval = 9.3
	}

	if this.Url == "" {
		this.Url = "http://127.0.0.1:8085/api/v1/servers"
	}

	if this.Summeries == "" {
		this.Summeries = "on"
	}
}
