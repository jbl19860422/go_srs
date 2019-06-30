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
	"fmt"
	"net/http"
	"strings"
	"net/url"
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
		m, _ := url.ParseQuery(r.URL.RawQuery)
		vHostParams, ok := m["vhost"]
		vhost := "__defaultVhost__"
		if ok {
			vhost = vHostParams[0]
		}
		source, ok := sourcePool[vhost + s]
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
		m, _ := url.ParseQuery(r.URL.RawQuery)
		vHostParams, ok := m["vhost"]
		vhost := "__defaultVhost__"
		if ok {
			vhost = vHostParams[0]
		}

		source, ok := sourcePool[vhost + s]
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
