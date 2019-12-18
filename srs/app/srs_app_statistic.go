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

import (
	"go_srs/srs/codec"
	"math/rand"
)

var srs_gvid int64 = rand.Int63n(10000000)

func srs_generate_id() int64 {
	srs_gvid++
	return srs_gvid
}

type SrsStatisticVhost struct {
	id 		int64
	vhost 	string
	nb_streams 	int
	nb_clients 	int
}

func NewSrsStatisticVhost() *SrsStatisticVhost {
	return &SrsStatisticVhost{
		id:srs_generate_id(),
	}
}

type SrsStatisticStreamVideo struct {
	vcodec 			codec.SrsCodecVideo		`json:"vcodec"`
	avc_profile		codec.SrsAvcProfile		`json:"avc_profile"`
	avc_level 		codec.SrsAvcLevel		`json:"avc_level"`
}

type SrsStatisticStreamAudio struct {
	acodec 			codec.SrsCodecAudio				`json:"acodec"`
	asample_rate 	codec.SrsCodecAudioSampleRate	`json:"asample_rate"`
	asound_type 	codec.SrsCodecAudioSoundType	`json:"asound_type"`
	aac_object 		codec.SrsAacObjectType			`json:"aac_object"`
}

type SrsStatisticStream struct {
	id 		int64				`json:"id"`
	vhost 	*SrsStatisticVhost
	app 	string				`json:"app"`
	stream 	string				`json:"name"`
	url 	string				`json:"url"`
	active 	bool				`json:"active"`
	connection_cid 	int			`json:"cid"`
	nb_clients 		int			`json:"clients"`
	nb_frames 		uint64		`json:"frames"`
	//video
	video 			*SrsStatisticStreamVideo 	`json:"video"`
	//audio
	audio 			*SrsStatisticStreamAudio	`json:"audio"`
}

func NewSrsStatisticStream() *SrsStatisticStream {
	return &SrsStatisticStream{
		id:srs_generate_id(),
		vhost:nil,
		connection_cid:-1,
		video:nil,
		audio:nil,
	}
}

func(this *SrsStatisticStream) Publish(cid int) {
	this.connection_cid = cid
	this.active = true
	this.vhost.nb_streams++
}

func(this *SrsStatisticStream) Close() {
	this.video = nil
	this.audio = nil
	this.vhost.nb_streams--
}

type SrsStatisticClient struct {
	stream 			*SrsStatisticStream
	id 				int
	create 			int64
}


type SrsStatistic struct {
	vhosts 		map[int64]*SrsStatisticVhost
	rvhosts 	map[string]*SrsStatisticVhost
	streams 	map[int64]*SrsStatisticStream
	rstreams 	map[string]*SrsStatisticStream
	clients 	map[int64]*SrsStatisticClient
}

func(this *SrsStatistic) FindVHost(vid int64) *SrsStatisticVhost {
	v, ok := this.vhosts[vid]
	if !ok {
		return nil
	}
	return v
}

func(this *SrsStatistic) FindStream(sid int64) *SrsStatisticStream {
	s, ok := this.streams[sid]
	if !ok {
		return nil
	}
	return s
}

func(this *SrsStatistic) FindClient(cid int64) *SrsStatisticClient {
	c, ok := this.clients[cid]
	if !ok {
		return nil
	}
	return c
}

func(this *SrsStatistic) OnVideoInfo(req *SrsRequest, vcodec codec.SrsCodecVideo, avc_profile codec.SrsAvcProfile, avc_level codec.SrsAvcLevel) {

}

func(this *SrsStatistic) createVHost(req *SrsRequest) *SrsStatisticVhost {
	v, ok := this.rvhosts[req.vhost]
	if !ok {
		v := NewSrsStatisticVhost()
		v.vhost = req.vhost
		this.rvhosts[req.vhost] = v
		this.vhosts[v.id] = v
		return v
	}
	return v
}

func(this *SrsStatistic) createStream(vhost *SrsStatisticVhost, req *SrsRequest) *SrsStatisticStream {
	url := req.GetStreamUrl()
	s, ok := this.rstreams[url]
	if !ok {
		s = NewSrsStatisticStream()
		s.vhost = vhost
		s.app = req.app
		s.stream = req.stream
		s.url = url
		this.rstreams[url] = s
		this.streams[s.id] = s
	}
	return s
}

var instance *SrsStatistic

func GetInstance() *SrsStatistic {
	if instance == nil {
		instance = &SrsStatistic{}
	}
	return instance
}
