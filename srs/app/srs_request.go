package app

import (
	"go_srs/srs/protocol/rtmp"
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
