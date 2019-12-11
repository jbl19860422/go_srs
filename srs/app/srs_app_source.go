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
	// "os"
	"sync"
	"errors"
	"fmt"
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
	handler 		ISrsSourceHandler
	conn			*SrsRtmpConn
	rtmp			*rtmp.SrsRtmpServer
	req 			*SrsRequest
	recvThread		*SrsRecvThread

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

	//record
	//dvr				*SrsDvr
	//dvr				*SrsDvrConsumer
	//hls				*SrsHls
	tsContext		*SrsTsContext
}

var sourcePoolMtx sync.Mutex
var sourcePool map[string]*SrsSource

func init() {
	sourcePool = make(map[string]*SrsSource)
}

func NewSrsSource(c *SrsRtmpConn, r *SrsRequest, h ISrsSourceHandler) *SrsSource {
	tsCtx := NewSrsTsContext()
	source := &SrsSource{
		req:r,
		conn:c,
		handler:h,
		rtmp:c.rtmp,
		gopCache:NewSrsGopCache(),
		atc:false,
		//hls:NewSrsHls(tsCtx),
		tsContext: tsCtx,
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

	source.recvThread = NewSrsRecvThread(c.rtmp, source, 1000)
	return source
}

func RemoveSrsSource(s *SrsSource) {
	sourcePoolMtx.Lock()
	defer sourcePoolMtx.Unlock()
	for k, v := range sourcePool {
		if v == s {
			fmt.Println("source removed")
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
	fmt.Println("create source url=", streamUrl)
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

	//if this.hls != nil {
	//	err := this.hls.onPublish(this.req, false)
	//	if err != nil {
	//		return err
	//	}
	//}

	if this.handler != nil {
		err := this.handler.OnPublish(this, this.req)
		if err != nil {
			return err
		}
	}
	return nil
}

func (this *SrsSource) onHlsStart() error {
	//if this.cacheSHVideo != nil {
	//	err := this.hls.on_video(this.cacheSHVideo)
	//	if err != nil {
	//		return err
	//	}
	//}
	//
	//if this.cacheSHAudio != nil {
	//	err := this.hls.on_audio(this.cacheSHAudio)
	//	if err != nil {
	//		return err
	//	}
	//}
	return nil
}

func (this *SrsSource) Handle(msg *rtmp.SrsRtmpMessage) error {
	if msg.GetHeader().IsAmf0Command() || msg.GetHeader().IsAmf3Command() {
		pkt, err := this.rtmp.DecodeMessage(msg)
		if err != nil {
			return err
		}
		_ = pkt
		//todo isfmle process
	}

	return this.ProcessPublishMessage(msg)
}

func (this *SrsSource) Initialize() {
	//this.dvr.Initialize(this, this.req)
	//this.hls.Initialize(this, this.req)
}

func (this *SrsSource) ProcessPublishMessage(msg *rtmp.SrsRtmpMessage) error {
	//todo fix edge process
	if msg.GetHeader().IsAudio() {
		// process audio
		if err := this.OnAudio(msg); err != nil {

		}
	}

	if msg.GetHeader().IsVideo() {
		if err := this.OnVideo(msg); err != nil {
		}
		//process video
	}
	//todo fix aggregate message
	//todo fix amf0 or amf3 data

	// process onMetaData
    if (msg.GetHeader().IsAmf0Data() || msg.GetHeader().IsAmf3Data()) {
		pkt, err := this.rtmp.DecodeMessage(msg)
		if err != nil {
			return err
		}

		switch pkt.(type) {
			case *packet.SrsOnMetaDataPacket: {
				err := this.OnMetaData(msg, pkt.(*packet.SrsOnMetaDataPacket))
				if err != nil {
					return err
				}
			}
		}
    }
	return nil
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
		fmt.Println("***********************AudioIsSequenceHeader len=", len(msg.GetPayload()), "*************************")
		this.cacheSHAudio = msg
	}

	for i := 0; i < len(this.consumers); i++ {
		this.consumers[i].Enqueue(msg, false, this.jitterAlgorithm)
		// fmt.Println("***********************************************send audio**************************************")
	}

	if err := this.gopCache.cache(msg); err != nil {
	}

	//if err := this.dvr.on_audio(msg); err != nil {
	//	return err
	//}
	//
	//if err := this.hls.on_audio(msg); err != nil {
	//	return err
	//}

	return nil
}

func (this *SrsSource) OnVideo(msg *rtmp.SrsRtmpMessage) error {
	isSequenceHeader := flvcodec.VideoIsSequenceHeader(msg.GetPayload())
	if isSequenceHeader {
		fmt.Println("***********************VideoIsSequenceHeader*************************")
		this.cacheSHVideo = msg
	}
	//fmt.Println("*************onVideo", len(this.consumers))
	for i := 0; i < len(this.consumers); i++ {
		this.consumers[i].Enqueue(msg, false, this.jitterAlgorithm)
		// fmt.Println("***********************************************send video**************************************")
	}

	if err := this.gopCache.cache(msg); err != nil {
		return err
	}

	//if err := this.dvr.on_video(msg); err != nil {
	//	return err
	//}

	//if err := this.hls.on_video(msg); err != nil {
	//	fmt.Println(err)
	//	return err
	//}

	// tsMsg := &SrsTsMessage{
	// 	payload:
	// }
	// this.tsContext.Encode(ts, codec.SrsCodecVideoAVC, codec.SrsCodecAudioAAC)

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
	fmt.Println("width=", width)
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
	fmt.Println("CreateConsumer len=", len(this.consumers))
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

func (this *SrsSource) CyclePublish() error {
	this.recvThread.Start()
	this.recvThread.Join()
	this.StopPublish()
	return nil
}

func (this *SrsSource) StopPublish() {
	//this.dvr.Close()
	for i := 0; i < len(this.consumers); i++ {
		this.consumers[i].OnUnpublish()
	}
	this.recvThread.Stop()
}
