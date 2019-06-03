package protocol

import (
	"bytes"
	"encoding/binary"
	"net/url"
	"strings"
	"math"
)

func IntToBytes(n int) []byte {
	x := int32(n)
	bytesBuffer := bytes.NewBuffer([]byte{})
	binary.Write(bytesBuffer, binary.LittleEndian, x)
	return bytesBuffer.Bytes()
}

func Float64ToByte(float float64) []byte {
	bits := math.Float64bits(float)
	bytes := make([]byte, 8)
	binary.BigEndian.PutUint64(bytes, bits)

	return bytes
}

type SrsAmf0Size struct {
}

func (s SrsAmf0Size) utf8(value string) int {
	return 2 + len(value)
}

func (s SrsAmf0Size) str(value string) int {
	return 1 + s.utf8(value)
}

func (s SrsAmf0Size) number() int {
	return 1 + 8
}

func (s SrsAmf0Size) date() int {
	return 1 + 8 + 2
}

func (s SrsAmf0Size) null() int {
	return 1
}

func (s SrsAmf0Size) undefined() int {
	return 1
}

func (s SrsAmf0Size) boolean() int {
	return 1 + 1
}

func (s SrsAmf0Size) object(obj *SrsAmf0Object) int {
	if obj == nil {
		return 0
	}

	return obj.total_size()
}

func (s SrsAmf0Size) object_eof() int {
	return 2 + 1
}

func (s SrsAmf0Size) any(v interface{}) int {
	var size int = 1
	switch v.(type) {
	case string:

	}
	return size
}

func Srs_discovery_tc_url(tcUrl string) (schema string, host string, vhost string, app string, stream string, port string, param string, err error) {
	var err1 error
	u, err1 := url.Parse(tcUrl)
	if err1 != nil {
		err = err1
		return
	}

	schema = u.Scheme
	host = u.Host
	port = SRS_CONSTS_RTMP_DEFAULT_PORT
	if len(u.Port()) >= 0 {
		port = u.Port()
	}

	m, _ := url.ParseQuery(u.RawQuery)
	vhost_params, ok := m["vhost"]
	if ok {
		vhost = vhost_params[0]
	}

	p := strings.Split(u.Path, "/")
	if len(p) >= 2 {
		app = p[1]
	}

	if len(p) >= 3 {
		stream = p[2]
	}

	param = u.RawQuery
	err = nil
	return
}