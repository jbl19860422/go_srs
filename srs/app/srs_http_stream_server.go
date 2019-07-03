package app

import (
	"fmt"
	"net/http"
)

type SrsHttpStreamServer struct {

}

func (this *SrsHttpStreamServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Println("url=", r.URL.Path)
}
