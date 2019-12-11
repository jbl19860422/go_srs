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

type VHostConf struct {
	Enabled              string          `json:"enabled"`
	MinLatency           string          `json:"min_latency"`
	GopCache             string          `json:"gop_cache"`
	QueueLength          uint32          `json:"queue_length"`
	SendMinInterval      uint32          `json:"send_min_interval"`
	ReduceSequenceHeader string          `json:"reduce_sequence_header"`
	Publish1stPktTimeout uint32          `json:"publish_1stpkt_timeout"`
	PublishNormalTimeout uint32          `json:"publish_normal_timeout"`
	Forward              []string        `json:"forward"`
	ChunkSize            uint32          `json:"chunk_size"`
	TimerJitter          string          `json:"time_jitter"`
	MixCorrect           string          `json:"mix_correct"`
	Atc                  string          `json:"atc"`
	AtcAuto              string          `json:"act_auto"`
	HeartBeat            *HeartBeatConf  `json:"heartbeat"`
	Stats                *StatsConf      `json:"stats"`
	HttpApi              *HttpApiConf    `json:"http_api"`
	HttpServer           *HttpServerConf `json:"http_server"`
	Security             *SecurityConf   `json:"security"`
	Dvr                  *DvrConf        `json:"dvr"`
	HttpStatic           *HttpStaticConf `json:"http_static"`
	HttpRemux            *HttpRemuxConf  `json:"http_remux"`
	Hls                  *HlsConf        `json:"hls"`
	HttpHooks            *HttpHooksConf  `json:"http_hooks"`
	Publish              *PublishConf    `json:"publish"`
}

func (this *VHostConf) initDefault() {
	this.MinLatency = "on"
	if this.GopCache == "" {
		this.GopCache = "on"
	}

	if this.QueueLength == 0 {
		this.QueueLength = 10
	}

	if this.ReduceSequenceHeader == "" {
		this.ReduceSequenceHeader = "off"
	}

	if this.Publish1stPktTimeout == 0 {
		this.Publish1stPktTimeout = 20000
	}

	if this.PublishNormalTimeout == 0 {
		this.PublishNormalTimeout = 7000
	}

	if this.ChunkSize == 0 {
		this.ChunkSize = 65000
	}

	if this.TimerJitter == "" {
		this.TimerJitter = "full"
	}

	if this.MixCorrect == "" {
		this.MixCorrect = "off"
	}

	if this.Atc == "" {
		this.Atc = "off"
	}

	if this.AtcAuto == "" {
		this.AtcAuto = "on"
	}

	if this.HeartBeat != nil {
		this.HeartBeat.initDefault()
	}

	if this.Stats != nil {
		this.Stats.initDefault()
	}

	if this.HttpApi != nil {
		this.HttpApi.initDefault()
	}

	if this.HttpServer != nil {
		this.HttpServer.initDefault()
	}

	if this.Security != nil {
		this.Security.initDefault()
	}

	if this.Dvr != nil {
		this.Dvr.initDefault()
	}

	if this.HttpStatic != nil {
		this.HttpStatic.initDefault()
	}

	if this.HttpRemux != nil {
		this.HttpRemux.initDefault()
	}

	if this.Hls != nil {
		this.Hls.initDefault()
	}

	if this.HttpHooks != nil {
		this.HttpHooks.initDefault()
	}

	if this.Publish != nil {
		this.Publish.initDefault()
	}
}