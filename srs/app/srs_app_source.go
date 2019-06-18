package app

import (
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

type SrsSource struct {
	handler 		ISrsSourceHandler
	conn			*SrsRtmpConn
	rtmp			*rtmp.SrsRtmpServer
	req 			*SrsRequest
	recvThread		*SrsRecvThread

	consumersMtx 	sync.Mutex
	consumers 		[]*SrsConsumer
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
		req:r,
		conn:c,
		handler:h,
		rtmp:c.rtmp,
		gopCache:NewSrsGopCache(),
		atc:false,
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
	vhost := r.vhost
	_ = vhost
	sourcePoolMtx.Lock()
	defer sourcePoolMtx.Unlock()
	if s, ok := sourcePool[streamUrl]; ok {
		return s, errors.New("source already in pool")
	}
	source = NewSrsSource(c, r, h)
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
				err := this.on_meta_data(msg, pkt.(*packet.SrsOnMetaDataPacket))
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
		this.consumers[i].StopPlay()
	}

	this.consumers = this.consumers[0:0]
}

func (this *SrsSource) OnAudio(msg *rtmp.SrsRtmpMessage) error {
	isSequenceHeader := flvcodec.AudioIsSequenceHeader(msg.GetPayload())
	if isSequenceHeader {
		fmt.Println("***********************AudioIsSequenceHeader*************************")
		this.cacheSHAudio = msg
	}

	for i := 0; i < len(this.consumers); i++ {
		this.consumers[i].Enqueue(msg, false, this.jitterAlgorithm)
		// fmt.Println("***********************************************send audio**************************************")
	}

	if err := this.gopCache.cache(msg); err != nil {
		return err
	}

	return nil
}

func (this *SrsSource) OnVideo(msg *rtmp.SrsRtmpMessage) error {
	isSequenceHeader := flvcodec.VideoIsSequenceHeader(msg.GetPayload())
	if isSequenceHeader {
		fmt.Println("***********************VideoIsSequenceHeader*************************")
		this.cacheSHVideo = msg
	}

	for i := 0; i < len(this.consumers); i++ {
		this.consumers[i].Enqueue(msg, false, this.jitterAlgorithm)
		// fmt.Println("***********************************************send video**************************************")
	}

	if err := this.gopCache.cache(msg); err != nil {
		return err
	}

	return nil
}

func (this *SrsSource) on_meta_data(msg *rtmp.SrsRtmpMessage, pkt *packet.SrsOnMetaDataPacket) error {
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
	
	this.cacheMetaData = rtmp.NewSrsRtmpMessage()
	this.cacheMetaData.SetHeader(*(msg.GetHeader()))
	
	this.cacheMetaData.GetHeader().SetLength(int32(len(stream.Data())))
	this.cacheMetaData.GetHeader().Print()
	this.cacheMetaData.SetPayload(stream.Data())

	for i := 0; i < len(this.consumers); i++ {
		this.consumers[i].Enqueue(this.cacheMetaData, false, this.jitterAlgorithm)
	}
    // when already got metadata, drop when reduce sequence header.
    // bool drop_for_reduce = false;
    // if (cache_metadata && _srs_config->get_reduce_sequence_header(_req->vhost)) {
    //     drop_for_reduce = true;
    //     srs_warn("drop for reduce sh metadata, size=%d", msg->size);
    // }
    
    // create a shared ptr message.
    // srs_freep(cache_metadata);
    // cache_metadata = new SrsSharedPtrMessage();
    
    // dump message to shared ptr message.
    // the payload/size managed by cache_metadata, user should not free it.
    // if ((ret = cache_metadata->create(&msg->header, payload, size)) != ERROR_SUCCESS) {
    //     srs_error("initialize the cache metadata failed. ret=%d", ret);
    //     return ret;
    // }
    
	// copy to all consumer
	//todo
    // if (!drop_for_reduce) {
    //     std::vector<SrsConsumer*>::iterator it;
    //     for (it = consumers.begin(); it != consumers.end(); ++it) {
    //         SrsConsumer* consumer = *it;
    //         if ((ret = consumer->enqueue(cache_metadata, atc, jitter_algorithm)) != ERROR_SUCCESS) {
    //             srs_error("dispatch the metadata failed. ret=%d", ret);
    //             return ret;
    //         }
    //     }
    // }
	
	//todo
    // copy to all forwarders
    // if (true) {
    //     std::vector<SrsForwarder*>::iterator it;
    //     for (it = forwarders.begin(); it != forwarders.end(); ++it) {
    //         SrsForwarder* forwarder = *it;
    //         if ((ret = forwarder->on_meta_data(cache_metadata)) != ERROR_SUCCESS) {
    //             srs_error("forwarder process onMetaData message failed. ret=%d", ret);
    //             return ret;
    //         }
    //     }
    // }
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
	
func (this *SrsSource) CreateConsumer(conn *SrsRtmpConn, ds bool, dm bool, db bool) *SrsConsumer {
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
	fmt.Println("CreateConsumer")
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

func (this *SrsSource) RemoveConsumer(consumer *SrsConsumer) {
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
	return nil
}

func (this *SrsSource) StopPublish() {
	this.recvThread.Stop()
}
