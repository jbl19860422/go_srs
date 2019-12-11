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

type HlsConf struct {
	Enabled         string  `json:"enabled"`
	HlsFragment     uint32  `json:"hls_fragment"`      //the hls fragment in seconds, the duration of a piece of ts.
	HlsTdRatio      float64 `json:"hls_td_ratio"`      //the hls m3u8 target duration ratio
	HlsAofRatio     float64 `json:"hls_aof_ration"`    //the audio overflow ratio.
	HlsWindow       uint32  `json:"hls_window"`        //the hls window in seconds, the number of ts in m3u8.
	HlsOnError      string  `json:"hls_on_error"`      //the error strategy
	HlsPath         string  `json:"hls_path"`          //the hls output path.
	HlsM3u8File     string  `json:"hls_m3u8_file"`     //the hls m3u8 file name.
	HlsTsFile       string  `json:"hls_ts_file"`       //the hls ts file name.
	HlsTsFloor      string  `json:"hls_ts_floor"`      //whether use floor for the hls_ts_file path generation.
	HlsEntryPrefix  string  `json:"hls_entry_prefix"`  //the hls entry prefix, which is base url of ts url.
	HlsAcodec       string  `json:"hls_acodec"`        //the default audio codec of hls.
	HlsVcodec       string  `json:"hls_vcodec"`        //the default video codec of hls.
	HlsCleanup      string  `json:"hls_cleanup"`       //whether cleanup the old expired ts files.
	HlsDispose      uint32  `json:"hls_dispose"`       //the timeout in seconds to dispose the hls,dispose is to remove all hls files, m3u8 and ts files.
	HlsNbNotify     uint32  `json:"hls_nb_notify"`     //the max size to notify hls,to read max bytes from ts of specified cdn network,
	HlsWaitKeyframe string  `json:"hls_wait_keyframe"` //whether wait keyframe to reap segment,
}

func (this *HlsConf) initDefault() {
	if this.Enabled == "" {
		this.Enabled = "off"
	}

	if this.HlsFragment == 0 {
		this.HlsFragment = 10
	}

	if this.HlsTdRatio <= 0 {
		this.HlsTdRatio = 1.5
	}

	if this.HlsAofRatio <= 0 {
		this.HlsAofRatio = 2.0
	}

	if this.HlsWindow <= 0 {
		this.HlsWindow = 60
	}

	if this.HlsOnError == "" {
		this.HlsOnError = "continue"
	}

	if this.HlsPath == "" {
		this.HlsPath = "html"
	}

	if this.HlsM3u8File == "" {
		this.HlsM3u8File = "[app]/[stream].m3u8"
	}

	if this.HlsTsFile == "" {
		this.HlsTsFile = "[app]/[stream]-[seq].ts"
	}

	if this.HlsTsFloor == "" {
		this.HlsTsFloor = "off"
	}

	if this.HlsAcodec == "" {
		this.HlsAcodec = "aac"
	}

	if this.HlsVcodec == "" {
		this.HlsVcodec = "avc"
	}

	if this.HlsCleanup == "" {
		this.HlsCleanup = "on"
	}

	if this.HlsNbNotify == 0 {
		this.HlsNbNotify = 64
	}

	if this.HlsWaitKeyframe == "" {
		this.HlsWaitKeyframe = "on"
	}
}