package utils

import (
	"bytes"
	"encoding/binary"
	"net/url"
	"strings"
	"go_srs/srs/global"
	"runtime"
	"log"
)

func numberToBytes(data interface{}, order binary.ByteOrder) []byte {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, order, data)
	if err != nil {
		return nil
	}
	return buf.Bytes()
}

func bytesToNumber(data []byte, order binary.ByteOrder, v interface{}) error {
	buf := bytes.NewReader(data)
	err := binary.Read(buf, order, v)
	return err
}

func UInt16ToBytes(data uint16, order binary.ByteOrder) []byte {
	return numberToBytes(data, order)
}

func UInt32ToBytes(data uint32, order binary.ByteOrder) []byte {
	return numberToBytes(data, order)
}

func UInt64ToBytes(data uint64, order binary.ByteOrder) []byte {
	return numberToBytes(data, order)
}

func Int16ToBytes(data int16, order binary.ByteOrder) []byte {
	return numberToBytes(data, order)
}

func Int32ToBytes(data int32, order binary.ByteOrder) []byte {
	return numberToBytes(data, order)
}

func Int64ToBytes(data int64, order binary.ByteOrder) []byte {
	return numberToBytes(data, order)
}

func Float32ToBytes(data float32, order binary.ByteOrder) []byte {
	return numberToBytes(data, order)
}

func Float64ToBytes(data float64, order binary.ByteOrder) []byte {
	return numberToBytes(data, order)
}

func BytesToUInt16(data []byte, order binary.ByteOrder) (uint16, error) {
	var v uint16 = 0
	err := bytesToNumber(data, order, &v)
	return v, err
}

func BytesToUInt32(data []byte, order binary.ByteOrder) (uint32, error) {
	var v uint32 = 0
	err := bytesToNumber(data, order, &v)
	return v, err
}

func BytesToUInt64(data []byte, order binary.ByteOrder) (uint64, error) {
	var v uint64 = 0
	err := bytesToNumber(data, order, &v)
	return v, err
}

func BytesToInt16(data []byte, order binary.ByteOrder) (int16, error) {
	var v int16 = 0
	err := bytesToNumber(data, order, &v)
	return v, err
}

func BytesToInt32(data []byte, order binary.ByteOrder) (int32, error) {
	var v int32 = 0
	err := bytesToNumber(data, order, &v)
	return v, err
}

func BytesToInt64(data []byte, order binary.ByteOrder) (int64, error) {
	var v int64 = 0
	err := bytesToNumber(data, order, &v)
	return v, err
}

func BytesToFloat32(data []byte, order binary.ByteOrder) (float32, error) {
	var v float32 = 0
	err := bytesToNumber(data, order, &v)
	return v, err
}

func BytesToFloat64(data []byte, order binary.ByteOrder) (float64, error) {
	var v float64 = 0
	err := bytesToNumber(data, order, &v)
	return v, err
}

func MinInt32(a int32, b int32) int32 {
	if a < b {
		return a
	}
	return b
}

func MinUInt32(a uint32, b uint32) uint32 {
	if a < b {
		return a
	}
	return b
}

func SrsDiscoveryTCUrl(tcUrl string) (schema string, host string, vhost string, app string, stream string, port string, param string, err error) {
	var err1 error
	u, err1 := url.Parse(tcUrl)
	if err1 != nil {
		err = err1
		return
	}

	schema = u.Scheme
	host = u.Host
	port = global.SRS_CONSTS_RTMP_DEFAULT_PORT
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

func SrsGenerateStreamUrl(vhost string, app string, stream string) string {
	var url string
	if vhost != global.SRS_CONSTS_RTMP_DEFAULT_VHOST {
		url += vhost
	}

    url += "/"
    url += app
    url += "/"
    url += stream

    return url
}

func TraceMemStats() {
    var ms runtime.MemStats
    runtime.ReadMemStats(&ms)
    log.Printf("Alloc:%d(bytes) HeapIdle:%d(bytes) HeapReleased:%d(bytes), RoutineNum:%d", ms.Alloc, ms.HeapIdle, ms.HeapReleased, runtime.NumGoroutine())
}

func Srs_path_build_stream(template_path string, vhost string, app string, stream string) string {
    path := template_path
    
    // variable [vhost]
    path = strings.Replace(path, "[vhost]", vhost, -1)
    // variable [app]
    path = strings.Replace(path, "[app]", app, -1)
    // variable [stream]
    path = strings.Replace(path, "[stream]", stream, -1)
    return path
}