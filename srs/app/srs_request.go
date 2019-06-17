package app

import (
	"go_srs/srs/protocol/rtmp"
	"go_srs/srs/utils"
)

type SrsRequest struct {
	ip             string
	typ            rtmp.SrsRtmpConnType
	tcUrl          string
	pageUrl        string
	swfUrl         string
	schema         string
	vhost          string
	host           string
	port           string
	app            string
	param          string
	stream         string
	duration       float64
	objectEncoding float64
}

func NewSrsRequest() *SrsRequest {
	return &SrsRequest{}
}

func (this SrsRequest) GetStreamUrl() string {
    return utils.SrsGenerateStreamUrl(this.vhost, this.app, this.stream);
}

