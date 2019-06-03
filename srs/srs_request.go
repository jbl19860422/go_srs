package srs

import "go_srs/srs/protocol"

type SrsRequest struct {
	ip             string
	typ            protocol.SrsRtmpConnType
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
