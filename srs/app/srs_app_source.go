package app

import (
	"sync"
	"errors"
)

type ISrsSourceHandler interface {
	OnPublish(s *SrsSource, r *SrsRequest) error
	OnUnpublish(s *SrsSource, r *SrsRequest) error
}

var sourcePoolMtx sync.Mutex
var sourcePool map[string]*SrsSource

type SrsSource struct {
	handler ISrsSourceHandler
	req 	*SrsRequest
}

func NewSrsSource() *SrsSource {
	return &SrsSource{}
}

func (this *SrsSource) Initialize(r *SrsRequest, h ISrsSourceHandler) error {
	this.handler = h
	this.req = r
	return nil
}

//TODO
func (this *SrsSource) SetCache(cache bool) {
	
}

func FetchOrCreate(r *SrsRequest, h ISrsSourceHandler) (*SrsSource, error) {
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

	source = NewSrsSource()
	if err := source.Initialize(r, h); err != nil {
		return nil, err
	}

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

func init() {
	sourcePool = make(map[string]*SrsSource)
}