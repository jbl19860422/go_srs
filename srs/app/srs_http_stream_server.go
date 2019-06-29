package app

import (
	"fmt"
	"net/http"
	"strings"
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
	c := NewSrsHttpTsConsumer(s, w, r)
	if err := s.AppendConsumer(c); err != nil {
		return nil
	}
	return c
}

func (this *SrsHttpStreamServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Println("url=", r.URL.Path)
	if strings.HasSuffix(r.URL.Path, ".ts") {
		s := strings.Replace(r.URL.Path, ".ts", "", -1)
		source, ok := sourcePool[s]
		if !ok {
			return
		}
		fmt.Println("Create Ts Consumer)")
		consumer := this.CreateTsConsumer(source, w, r)
		err := consumer.PlayCycle()
		_ = err
		return
	} else if strings.HasSuffix(r.URL.Path, ".flv") {
		s := strings.Replace(r.URL.Path, ".flv", "", -1)
		source, ok := sourcePool[s]
		if !ok {
			return
		}
		fmt.Println("Create flv Consumer)")
		consumer := this.CreateFlvConsumer(source, w, r)
		err := consumer.PlayCycle()
		_ = err
		return
	}
}
