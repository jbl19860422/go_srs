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

func (this *SrsHttpStreamServer) Mount(path string, s *SrsSource) {
	this.sources[path] = s
}

func (this *SrsHttpStreamServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Println("url=", r.URL.Path)
	source, ok := this.sources[r.URL.Path]
	if !ok {
		io.WriteString(w, "404")
		return
	}

	consumer := source.CreateConsumer(nil, true, true, true)
	go func() {
		for {
			
		}
	}()
}
