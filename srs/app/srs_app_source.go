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

package app

import (
	"sync"
	"errors"
	"go_srs/srs/protocol/rtmp"
	"go_srs/srs/codec/flv"
	"go_srs/srs/protocol/packet"
	"go_srs/srs/global"
	"go_srs/srs/utils"
)

type ISrsSourceHandler interface {
	OnPublish(s *SrsSource, r *SrsRequest) error
	OnUnpublish(s *SrsSource, r *SrsRequest) error
}

type SrsSHRequester interface {
	GetSH(metaData *rtmp.SrsRtmpMessage, audioSH *rtmp.SrsRtmpMessage, videoSH *rtmp.SrsRtmpMessage)
}

type SrsSource struct {
	source_id 		int64
	handler 		ISrsSourceHandler
	conn			*SrsRtmpConn
	rtmp			*rtmp.SrsRtmpServer
	req 			*SrsRequest

	consumersMtx 	sync.Mutex
	consumers 		[]Consumer
	gopCache		*SrsGopCache
	cacheSHVideo 	*rtmp.SrsRtmpMessage
	cacheSHAudio 	*rtmp.SrsRtmpMessage
	cacheMetaData 	*rtmp.SrsRtmpMessage

	/**
    * atc whether atc(use absolute time and donot adjust time),
    * directly use msg time and donot adjust if atc is true,
    * otherwise, adjust msg time to start from 0 to make flash happy.
    */
	// TODO: FIXME: to support reload atc.
	atc 			bool
	jitterAlgorithm *SrsRtmpJitterAlgorithm
}

var sourcePoolMtx sync.Mutex
var sourcePool map[string]*SrsSource

func init() {
	sourcePool = make(map[string]*SrsSource)
}

func NewSrsSource(c *SrsRtmpConn, r *SrsRequest, h ISrsSourceHandler) *SrsSource {
	source := &SrsSource{
		source_id:c.id,
		req:r,
		conn:c,
		handler:h,
		rtmp:c.rtmp,
		gopCache:NewSrsGopCache(),
		atc:false,
	}

	dvrConsumer := NewSrsDvrConsumer(source, r)
	if dvrConsumer != nil {
		source.AppendConsumer(dvrConsumer)
		go func(){
			dvrConsumer.ConsumeCycle()
		}()
	}

	hlsConsumer := NewSrsHlsConsumer(source, r)
	if dvrConsumer != nil {
		source.AppendConsumer(hlsConsumer)
		go func(){
			hlsConsumer.ConsumeCycle()
		}()
	}

	return source
}

func RemoveSrsSource(s *SrsSource) {
	sourcePoolMtx.Lock()
	defer sourcePoolMtx.Unlock()
	for k, v := range sourcePool {
		if v == s {
			delete(sourcePool,k)
		}
	}
}

func FetchOrCreate(c *SrsRtmpConn, r *SrsRequest, h ISrsSourceHandler) (*SrsSource, error) {
	source := FetchSource(r)
	if source != nil {
		return source, nil
	}

	streamUrl := r.GetStreamUrl()
	vhost := r.vhost
	_ = vhost
	sourcePoolMtx.Lock()
	defer sourcePoolMtx.Unlock()
	if s, ok := sourcePool[streamUrl]; ok {
		return s, errors.New("source already in pool")
	}
	source = NewSrsSource(c, r, h)
	//todo fix return value
	source.Initialize()
	sourcePool[streamUrl] = source
	return source, nil
}

func FetchSource(r *SrsRequest) *SrsSource {
	sourcePoolMtx.Lock()
	defer sourcePoolMtx.Unlock()
	streamUrl := r.GetStreamUrl()
	source, ok := sourcePool[streamUrl]
	if !ok {
		return nil
	}

	//TODO
	// we always update the request of resource, 
    // for origin auth is on, the token in request maybe invalid,
    // and we only need to update the token of request, it's simple.
	//source->req->update_auth(r)
	return source
}

func (this *SrsSource) OnRequestSH(requester SrsSHRequester) error {
	if this.cacheMetaData == nil {
		return errors.New("missing metadata")
	}

	if this.cacheSHAudio == nil {
		return errors.New("missing audio sh")
	}

	if this.cacheSHVideo == nil {
		return errors.New("missing video sh")
	}

	requester.GetSH(this.cacheMetaData, this.cacheSHAudio, this.cacheSHVideo)
	return nil
}


func (this *SrsSource) on_dvr_request_sh() error {
	//if this.cacheMetaData != nil {
	//	if err := this.dvr.OnMetaData(this.cacheMetaData); err != nil {
	//		return err
	//	}
	//}
	//
	//if this.cacheSHVideo != nil {
	//	if err := this.dvr.on_video(this.cacheSHVideo); err != nil {
	//		return err
	//	}
	//}
	//
	//if this.cacheSHAudio != nil {
	//	fmt.Println("on_dvr_request_sh audio len=", len(this.cacheSHAudio.GetPayload()))
	//	if err := this.dvr.on_audio(this.cacheSHAudio); err != nil {
	//		return err
	//	}
	//}
	return nil
}

func (this *SrsSource) onPublish() error {
	for i := 0; i < len(this.consumers); i++ {
		this.consumers[i].OnPublish()
	}

	if this.handler != nil {
		err := this.handler.OnPublish(this, this.req)
		if err != nil {
			return err
		}
	}

	stat := GetStatisticInstance()
	stat.OnStreamPublish(this.req, this.source_id)
	return nil
}

func (this *SrsSource) Initialize() {
}

func (this *SrsSource) OnRecvError(err error) {
	RemoveSrsSource(this)
}

func (this *SrsSource) RemoveConsumers() {
	this.consumersMtx.Lock()
	defer this.consumersMtx.Unlock()

	for i := 0; i < len(this.consumers); i++ {
		this.consumers[i].StopConsume()
	}

	this.consumers = this.consumers[0:0]
}

func (this *SrsSource) OnAudio(msg *rtmp.SrsRtmpMessage) error {
	isSequenceHeader := flvcodec.AudioIsSequenceHeader(msg.GetPayload())
	if isSequenceHeader {
		this.cacheSHAudio = msg
	}

	for i := 0; i < len(this.consumers); i++ {
		this.consumers[i].Enqueue(msg, false, this.jitterAlgorithm)
	}

	if err := this.gopCache.cache(msg); err != nil {
	}
	return nil
}

func (this *SrsSource) OnVideo(msg *rtmp.SrsRtmpMessage) error {
	isSequenceHeader := flvcodec.VideoIsSequenceHeader(msg.GetPayload())
	if isSequenceHeader {
		this.cacheSHVideo = msg
	}

	for i := 0; i < len(this.consumers); i++ {
		this.consumers[i].Enqueue(msg, false, this.jitterAlgorithm)
	}

	if err := this.gopCache.cache(msg); err != nil {
		return err
	}

	return nil
}

func (this *SrsSource) OnMetaData(msg *rtmp.SrsRtmpMessage, pkt *packet.SrsOnMetaDataPacket) error {
    // SrsAmf0Any* prop = NULL;
	
	//todo
    // when exists the duration, remove it to make ExoPlayer happy.
    // if (metadata->metadata->get_property("duration") != NULL) {
    //     metadata->metadata->remove("duration");
    // }
    
    // generate metadata info to print
    // std::stringstream ss;
    // if ((prop = metadata->metadata->ensure_property_number("width")) != NULL) {
    //     ss << ", width=" << (int)prop->to_number();
    // }
    // if ((prop = metadata->metadata->ensure_property_number("height")) != NULL) {
    //     ss << ", height=" << (int)prop->to_number();
    // }
    // if ((prop = metadata->metadata->ensure_property_number("videocodecid")) != NULL) {
    //     ss << ", vcodec=" << (int)prop->to_number();
    // }
    // if ((prop = metadata->metadata->ensure_property_number("audiocodecid")) != NULL) {
    //     ss << ", acodec=" << (int)prop->to_number();
    // }
    // srs_trace("got metadata%s", ss.str().c_str());
	var width float64
	_ = pkt.Get("width", &width)
	pkt.Set("server", global.RTMP_SIG_SRS_SERVER)
	pkt.Set("srs_primary", global.RTMP_SIG_SRS_PRIMARY)
	pkt.Set("srs_authors", global.RTMP_SIG_SRS_AUTHROS)
	// version, for example, 1.0.0
	// add version to metadata, please donot remove it, for debug.
	pkt.Set("server_version", global.RTMP_SIG_SRS_VERSION)
	
	// if allow atc_auto and bravo-atc detected, open atc for vhost.
	//todo
    // atc = _srs_config->get_atc(_req->vhost);
    // if (_srs_config->get_atc_auto(_req->vhost)) {
    //     if ((prop = metadata->metadata->get_property("bravo_atc")) != NULL) {
    //         if (prop->is_string() && prop->to_str() == "true") {
    //             atc = true;
    //         }
    //     }
    // }
    
	// encode the metadata to payload
	d := make([]byte, 0)
	stream := utils.NewSrsStream(d)
	if err := pkt.Encode(stream); err != nil {
		return err
	}
	
	//this.cacheMetaData = rtmp.NewSrsRtmpMessage()
	//this.cacheMetaData.SetHeader(*(msg.GetHeader()))
	//
	//this.cacheMetaData.GetHeader().SetLength(int32(len(stream.Data())))
	//this.cacheMetaData.GetHeader().Print()
	//this.cacheMetaData.SetPayload(stream.Data())
	this.cacheMetaData = msg
	for i := 0; i < len(this.consumers); i++ {
		this.consumers[i].Enqueue(msg, false, this.jitterAlgorithm)
	}

	//if err := this.dvr.OnMetaData(msg); err != nil {
		//return err
    //}
    return nil
}

//TODO
func (this *SrsSource) SetCache(cache bool) {
	
}

/**
* create consumer and dumps packets in cache.
* @param consumer, output the create consumer.
* @param ds, whether dumps the sequence header.
* @param dm, whether dumps the metadata.
* @param dg, whether dumps the gop cache.
*/
func (this *SrsSource) CreateConsumer(conn *SrsRtmpConn, ds bool, dm bool, db bool) Consumer {
	this.consumersMtx.Lock()
	consumer := NewSrsConsumer(this, conn)
	this.consumers = append(this.consumers, consumer)
	this.consumersMtx.Unlock()
	//todo set queue size
	//todo process atc
	//todo copy meta data
	//todo cppy sequence header
	//todo copy gop to consumers queue
	//many things todo
	if this.cacheMetaData != nil {
		consumer.Enqueue(this.cacheMetaData, false, this.jitterAlgorithm)
	}

	if this.cacheSHVideo != nil {
		consumer.Enqueue(this.cacheSHVideo, false, this.jitterAlgorithm)
	}
	
	if this.cacheSHAudio != nil {
		consumer.Enqueue(this.cacheSHAudio, false, this.jitterAlgorithm)
	}

	if err := this.gopCache.dump(consumer, false, this.jitterAlgorithm); err != nil {
		return nil
	}

	return consumer
}

func (this *SrsSource) AppendConsumer(consumer Consumer) error {
	this.consumersMtx.Lock()
	this.consumers = append(this.consumers, consumer)
	this.consumersMtx.Unlock()
	//todo set queue size
	//todo process atc
	//todo copy meta data
	//todo cppy sequence header
	//todo copy gop to consumers queue
	//many things todo 
	if this.cacheMetaData != nil {
		consumer.Enqueue(this.cacheMetaData, false, this.jitterAlgorithm)
	}

	if this.cacheSHVideo != nil {
		consumer.Enqueue(this.cacheSHVideo, false, this.jitterAlgorithm)
	}
	
	if this.cacheSHAudio != nil {
		consumer.Enqueue(this.cacheSHAudio, false, this.jitterAlgorithm)
	}

	if err := this.gopCache.dump(consumer, false, this.jitterAlgorithm); err != nil {
		return err
	}
	return nil
}

func (this *SrsSource) RemoveConsumer(consumer Consumer) {
	this.consumersMtx.Lock()
	defer this.consumersMtx.Unlock()
	for i := 0; i < len(this.consumers); i++ {
		if this.consumers[i] == consumer {
			this.consumers = append(this.consumers[:i], this.consumers[i+1:]...)
		}
	}
}

func (this *SrsSource) StopPublish() {
	for i := 0; i < len(this.consumers); i++ {
		this.consumers[i].OnUnpublish()
	}

	stat := GetStatisticInstance()
	stat.OnStreamClose(this.req, this.source_id)
}
