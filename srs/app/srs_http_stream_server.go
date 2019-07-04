package app

import (
	"fmt"
	"io"
	"net/http"
)

type SrsHttpStreamServer struct {
	sources map[string]*SrsSource
}

func NewSrsHttpStreamServer() *SrsHttpStreamServer {
	return &SrsHttpStreamServer{
		sources:make(map[string]*SrsSource),
	}
}

func (this *SrsHttpStreamServer) Mount(r *SrsRequest, s *SrsSource) {
	path := r.GetStreamUrl()
	path += ".flv"
	this.sources[path] = s
}

func (this *SrsHttpStreamServer) CreateFlvConsumer(s *SrsSource, w http.ResponseWriter, r *http.Request) Consumer {
	c := NewSrsHttpFlvConsumer(s, w, r)
	if err := s.AppendConsumer(c); err != nil {
		return nil
	}
	return c
}

func (this *SrsHttpStreamServer) CreateTsConsumer(s *SrsSource, w http.ResponseWriter, r *http.Request) Consumer {
	c := NewSrsHttpFlvConsumer(s, w, r)
	if err := s.AppendConsumer(c); err != nil {
		return nil
	}
	return c
}

func (this *SrsHttpStreamServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Println("url=", r.URL.Path)
	source, ok := this.sources[r.URL.Path]
	if !ok {
		fmt.Println("not find for", r.URL.Path)
		for k, _ := range this.sources {
			fmt.Println("k=", k)
		}
		io.WriteString(w, "404")
		return
	}

	fmt.Println("*****************create consumer for", r.URL.Path)
	consumer := this.CreateFlvConsumer(source, w, r)
	err := consumer.PlayCycle()
	_ = err
}
