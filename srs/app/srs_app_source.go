package app

import (
	"sync"
	"errors"
	"fmt"
	"go_srs/srs/protocol/rtmp"
)

type ISrsSourceHandler interface {
	OnPublish(s *SrsSource, r *SrsRequest) error
	OnUnpublish(s *SrsSource, r *SrsRequest) error
}

var sourcePoolMtx sync.Mutex
var sourcePool map[string]*SrsSource

type SrsSource struct {
	handler 	ISrsSourceHandler
	req 		*SrsRequest
	consumers 	[]*SrsConsumer
}

func NewSrsSource() *SrsSource {
	return &SrsSource{}
}

func (this *SrsSource) Initialize(r *SrsRequest, h ISrsSourceHandler) error {
	this.handler = h
	this.req = r
	return nil
}

func (this *SrsSource) OnAudio(msg *rtmp.SrsRtmpMessage) error {
	for i := 0; i < len(this.consumers); i++ {
		this.consumers[i].Enqueue(msg, false)
		fmt.Println("***********************************************send audio**************************************")
	}
	return nil
}

func (this *SrsSource) OnVideo(msg *rtmp.SrsRtmpMessage) error {
	for i := 0; i < len(this.consumers); i++ {
		this.consumers[i].Enqueue(msg, false)
		fmt.Println("***********************************************send video**************************************")
	}
	return nil
}

//TODO
func (this *SrsSource) SetCache(cache bool) {
	
}

func FetchOrCreate(r *SrsRequest, h ISrsSourceHandler) (*SrsSource, error) {
	fmt.Println("**********FetchOrCreate**********")
	source := FetchSource(r)
	if source != nil {
		return source, nil
	}

	streamUrl := r.GetStreamUrl()
	vhost := r.vhost
	_ = vhost
	fmt.Println("**********streamUrl=", streamUrl)
	sourcePoolMtx.Lock()
	defer sourcePoolMtx.Unlock()
	if s, ok := sourcePool[streamUrl]; ok {
		return s, errors.New("source already in pool")
	}

	source = NewSrsSource()
	if err := source.Initialize(r, h); err != nil {
		return nil, err
	}

	sourcePool[streamUrl] = source
	return source, nil
}

/**
* create consumer and dumps packets in cache.
* @param consumer, output the create consumer.
* @param ds, whether dumps the sequence header.
* @param dm, whether dumps the metadata.
* @param dg, whether dumps the gop cache.
*/
	
func (this *SrsSource) CreateConsumer(conn *SrsRtmpConn, ds bool, dm bool, db bool) *SrsConsumer {
	consumer := NewSrsConsumer(this, conn)
	this.consumers = append(this.consumers, consumer)
	//todo set queue size
	//todo process atc
	//todo copy meta data
	//todo cppy sequence header
	//todo copy gop to consumers queue
	//many things todo 
	return consumer
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

func init() {
	sourcePool = make(map[string]*SrsSource)
}