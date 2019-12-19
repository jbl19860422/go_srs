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

import (
	"encoding/json"
	"io/ioutil"
	"sync"
)

type SrsConfig struct {
	ListenPort     uint32                `json:"listen_port"`
	Pid            string                `json:"pid"`
	ChunkSize      uint32                `json:"chunk_size"`
	MaxConnections uint32                `json:"max_connection"`
	WorkDir        string                `json:"work_dir"`
	VHosts         map[string]*VHostConf `json:"vhosts"`
}

func (this *SrsConfig) GetVHost(name string) *VHostConf {
	h, ok := this.VHosts[name]
	if !ok {
		return nil
	}
	return h
}

func (this *SrsConfig) initDefault() {
	if this.ListenPort == 0 {
		this.ListenPort = 1935
	}

	if this.Pid == "" {
		this.Pid = "./srs.pid"
	}

	if this.ChunkSize == 0 {
		this.ChunkSize = 60000
	}

	if this.MaxConnections == 0 {
		this.MaxConnections = 1000
	}

	if this.WorkDir == "" {
		this.WorkDir = "./"
	}

	for _, v := range this.VHosts {
		v.initDefault()
	}
}

const SRS_CONF_DEFAULT_HLS_FRAGMENT = 10

func GetHlsFragment(vname string) uint32 {
	vhost := GetInstance().GetVHost(vname)
	if vhost == nil {
		return SRS_CONF_DEFAULT_HLS_FRAGMENT
	}

	return vhost.Hls.HlsFragment
}

const SRS_CONF_DEFAULT_HLS_WINDOW = 60
func GetHlsWindow(vname string) uint32 {
	vhost := GetInstance().GetVHost(vname)
	if vhost == nil {
		return SRS_CONF_DEFAULT_HLS_FRAGMENT
	}

	return vhost.Hls.HlsWindow
}

func GetHlsEntryPrefix(vname string) string {
	vhost := GetInstance().GetVHost(vname)
	if vhost == nil {
		return ""
	}

	return vhost.Hls.HlsEntryPrefix
}

const SRS_CONF_DEFAULT_HLS_PATH = "./html"

func GetHlsPath(vname string) string {
	vhost := GetInstance().GetVHost(vname)
	if vhost == nil {
		return SRS_CONF_DEFAULT_HLS_PATH
	}

	return vhost.Hls.HlsPath
}

const SRS_CONF_DEFAULT_HLS_M3U8_FILE = "[app]/[stream].m3u8"

func GetHlsM3u8File(vname string) string {
	vhost := GetInstance().GetVHost(vname)
	if vhost == nil {
		return SRS_CONF_DEFAULT_HLS_M3U8_FILE
	}

	return vhost.Hls.HlsM3u8File
}

const SRS_CONF_DEFAULT_HLS_TS_FILE = "[app]/[stream]-[seq].ts"

func GetHlsTsFile(vname string) string {
	vhost := GetInstance().GetVHost(vname)
	if vhost == nil {
		return SRS_CONF_DEFAULT_HLS_TS_FILE
	}

	return vhost.Hls.HlsTsFile
}

const SRS_CONF_DEFAULT_HLS_CLEANUP = true

func GetHlsCleanup(vname string) bool {
	vhost := GetInstance().GetVHost(vname)
	if vhost == nil {
		return SRS_CONF_DEFAULT_HLS_CLEANUP
	}

	return vhost.Hls.HlsCleanup == "on"
}

const SRS_CONF_DEFAULT_HLS_WAIT_KEYFRAME = true

func GetHlsWaitKeyframe(vname string) bool {
	vhost := GetInstance().GetVHost(vname)
	if vhost == nil {
		return SRS_CONF_DEFAULT_HLS_WAIT_KEYFRAME
	}

	return vhost.Hls.HlsWaitKeyframe == "on"
}

func (this *SrsConfig) GetChunkSize(vhost string) uint32 {
	h, ok := this.VHosts[vhost]
	if !ok {
		return this.ChunkSize
	}

	if h.Enabled != "on" {
		return this.ChunkSize
	}

	return h.ChunkSize
}

func GetDvrPath(vhost string) string {
	h := GetInstance().GetVHost(vhost)
	if h == nil {
		return SRS_CONF_DEFAULT_DVR_PATH
	}

	if h.Enabled != "on" || h.Dvr == nil || h.Dvr.Enabled != "on" {
		return SRS_CONF_DEFAULT_DVR_PATH
	}

	return h.Dvr.DvrPath
}

const SRS_CONF_DEFAULT_DVR_PLAN_SESSION = "session"
const SRS_CONF_DEFAULT_DVR_PLAN_SEGMENT = "segment"
const SRS_CONF_DEFAULT_DVR_PLAN_APPEND = "append"

const SRS_CONF_DEFAULT_DVR_PLAN = SRS_CONF_DEFAULT_DVR_PLAN_SESSION

func GetDvrPlan(vhost string) string {
	h := GetInstance().GetVHost(vhost)
	if h == nil {
		return SRS_CONF_DEFAULT_DVR_PLAN
	}

	if h.Enabled != "on" || h.Dvr == nil || h.Dvr.Enabled != "on" {
		return SRS_CONF_DEFAULT_DVR_PLAN
	}

	return h.Dvr.DvrPlan
}

const SRS_CONF_DEFAULT_1STPKT_TIMEOUT = 2000
func GetPublish1stpktTimeout(vhost string) uint32 {
	h := GetInstance().GetVHost(vhost)
	if h == nil {
		return SRS_CONF_DEFAULT_1STPKT_TIMEOUT
	}

	if h.Enabled != "on" {
		return h.Publish1stPktTimeout
	}

	return SRS_CONF_DEFAULT_1STPKT_TIMEOUT
}

const SRS_CONF_DEFAULT_NORPKT_TIMEOUT = 5000
	func GetPublishNormalPktTimeout(vhost string) uint32 {
	h := GetInstance().GetVHost(vhost)
	if h == nil {
		return SRS_CONF_DEFAULT_NORPKT_TIMEOUT
	}

	if h.Enabled != "on" {
		return h.PublishNormalTimeout
	}

	return SRS_CONF_DEFAULT_NORPKT_TIMEOUT
}

var config *SrsConfig
var once sync.Once
func GetInstance() *SrsConfig {
	once.Do(func() {
		config = &SrsConfig{}
	})

	return config
}

func (this *SrsConfig) Init(file string) error {
	this.ListenPort = 1935
	this.Pid = "./srs.pid"
	this.ChunkSize = 60000
	this.MaxConnections = 1000
	this.WorkDir = "./"
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}

	err = json.Unmarshal(data, this)
	if err != nil {
		return err
	}

	this.initDefault()
	return nil
}

func init() {
	config = &SrsConfig{}
}
