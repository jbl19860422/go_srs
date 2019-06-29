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

var CRCTable = [256]uint32{
	0x00000000, 0x04c11db7, 0x09823b6e, 0x0d4326d9, 0x130476dc, 0x17c56b6b,
	0x1a864db2, 0x1e475005, 0x2608edb8, 0x22c9f00f, 0x2f8ad6d6, 0x2b4bcb61,
	0x350c9b64, 0x31cd86d3, 0x3c8ea00a, 0x384fbdbd, 0x4c11db70, 0x48d0c6c7,
	0x4593e01e, 0x4152fda9, 0x5f15adac, 0x5bd4b01b, 0x569796c2, 0x52568b75,
	0x6a1936c8, 0x6ed82b7f, 0x639b0da6, 0x675a1011, 0x791d4014, 0x7ddc5da3,
	0x709f7b7a, 0x745e66cd, 0x9823b6e0, 0x9ce2ab57, 0x91a18d8e, 0x95609039,
	0x8b27c03c, 0x8fe6dd8b, 0x82a5fb52, 0x8664e6e5, 0xbe2b5b58, 0xbaea46ef,
	0xb7a96036, 0xb3687d81, 0xad2f2d84, 0xa9ee3033, 0xa4ad16ea, 0xa06c0b5d,
	0xd4326d90, 0xd0f37027, 0xddb056fe, 0xd9714b49, 0xc7361b4c, 0xc3f706fb,
	0xceb42022, 0xca753d95, 0xf23a8028, 0xf6fb9d9f, 0xfbb8bb46, 0xff79a6f1,
	0xe13ef6f4, 0xe5ffeb43, 0xe8bccd9a, 0xec7dd02d, 0x34867077, 0x30476dc0,
	0x3d044b19, 0x39c556ae, 0x278206ab, 0x23431b1c, 0x2e003dc5, 0x2ac12072,
	0x128e9dcf, 0x164f8078, 0x1b0ca6a1, 0x1fcdbb16, 0x018aeb13, 0x054bf6a4,
	0x0808d07d, 0x0cc9cdca, 0x7897ab07, 0x7c56b6b0, 0x71159069, 0x75d48dde,
	0x6b93dddb, 0x6f52c06c, 0x6211e6b5, 0x66d0fb02, 0x5e9f46bf, 0x5a5e5b08,
	0x571d7dd1, 0x53dc6066, 0x4d9b3063, 0x495a2dd4, 0x44190b0d, 0x40d816ba,
	0xaca5c697, 0xa864db20, 0xa527fdf9, 0xa1e6e04e, 0xbfa1b04b, 0xbb60adfc,
	0xb6238b25, 0xb2e29692, 0x8aad2b2f, 0x8e6c3698, 0x832f1041, 0x87ee0df6,
	0x99a95df3, 0x9d684044, 0x902b669d, 0x94ea7b2a, 0xe0b41de7, 0xe4750050,
	0xe9362689, 0xedf73b3e, 0xf3b06b3b, 0xf771768c, 0xfa325055, 0xfef34de2,
	0xc6bcf05f, 0xc27dede8, 0xcf3ecb31, 0xcbffd686, 0xd5b88683, 0xd1799b34,
	0xdc3abded, 0xd8fba05a, 0x690ce0ee, 0x6dcdfd59, 0x608edb80, 0x644fc637,
	0x7a089632, 0x7ec98b85, 0x738aad5c, 0x774bb0eb, 0x4f040d56, 0x4bc510e1,
	0x46863638, 0x42472b8f, 0x5c007b8a, 0x58c1663d, 0x558240e4, 0x51435d53,
	0x251d3b9e, 0x21dc2629, 0x2c9f00f0, 0x285e1d47, 0x36194d42, 0x32d850f5,
	0x3f9b762c, 0x3b5a6b9b, 0x0315d626, 0x07d4cb91, 0x0a97ed48, 0x0e56f0ff,
	0x1011a0fa, 0x14d0bd4d, 0x19939b94, 0x1d528623, 0xf12f560e, 0xf5ee4bb9,
	0xf8ad6d60, 0xfc6c70d7, 0xe22b20d2, 0xe6ea3d65, 0xeba91bbc, 0xef68060b,
	0xd727bbb6, 0xd3e6a601, 0xdea580d8, 0xda649d6f, 0xc423cd6a, 0xc0e2d0dd,
	0xcda1f604, 0xc960ebb3, 0xbd3e8d7e, 0xb9ff90c9, 0xb4bcb610, 0xb07daba7,
	0xae3afba2, 0xaafbe615, 0xa7b8c0cc, 0xa379dd7b, 0x9b3660c6, 0x9ff77d71,
	0x92b45ba8, 0x9675461f, 0x8832161a, 0x8cf30bad, 0x81b02d74, 0x857130c3,
	0x5d8a9099, 0x594b8d2e, 0x5408abf7, 0x50c9b640, 0x4e8ee645, 0x4a4ffbf2,
	0x470cdd2b, 0x43cdc09c, 0x7b827d21, 0x7f436096, 0x7200464f, 0x76c15bf8,
	0x68860bfd, 0x6c47164a, 0x61043093, 0x65c52d24, 0x119b4be9, 0x155a565e,
	0x18197087, 0x1cd86d30, 0x029f3d35, 0x065e2082, 0x0b1d065b, 0x0fdc1bec,
	0x3793a651, 0x3352bbe6, 0x3e119d3f, 0x3ad08088, 0x2497d08d, 0x2056cd3a,
	0x2d15ebe3, 0x29d4f654, 0xc5a92679, 0xc1683bce, 0xcc2b1d17, 0xc8ea00a0,
	0xd6ad50a5, 0xd26c4d12, 0xdf2f6bcb, 0xdbee767c, 0xe3a1cbc1, 0xe760d676,
	0xea23f0af, 0xeee2ed18, 0xf0a5bd1d, 0xf464a0aa, 0xf9278673, 0xfde69bc4,
	0x89b8fd09, 0x8d79e0be, 0x803ac667, 0x84fbdbd0, 0x9abc8bd5, 0x9e7d9662,
	0x933eb0bb, 0x97ffad0c, 0xafb010b1, 0xab710d06, 0xa6322bdf, 0xa2f33668,
	0xbcb4666d, 0xb8757bda, 0xb5365d03, 0xb1f740b4,
}

func MpegtsCRC32(data []byte) uint32 {
	var crc uint32 = 0xffffffff
	
	for i := 0; i < len(data); i++ {
		crc = (crc << 8) ^ CRCTable[((crc >> 24) ^ uint32(data[i])) & 0xff]
	}

    return crc;
}


func GetNalu(stream *SrsStream) ([]([]byte), error) {
	payload := stream.GetPayload()
	var i := 0
	nalus := make([]([]byte), 0)
	for (i + 3) < len(payload) {
		if payload[i] != 0x00 || payload[i + 1] != 0x00 {
			return nalus, nil
		}

		if payload[i + 2] == 0x01 {
			
		}
	}
}

func FindAnnexbStartCodePos(stream *SrsStream, start int, startCodePos *int, startCodeCount *int) error {
	p := stream.GetPayload()
	payload := p[start:]
	var i int := 0
	for (i + 3 ) < len(payload) {
		if payload[i] == 0x00 && payload[i + 1] == 0x00 {
			if payload[i + 2] == 0x01 {
				*startCodePos = start + i
				*startCodeCount = 3
				return nil
			} else if payload[i + 2] == 0x00 {
				*startcodePos = start + i
				*startCodeCount = 4
				return nil
			}
		}
	}
	return len(payload)
}

func AvcStartswithAnnexb(stream *SrsStream, startCodeByteCount *int) bool {
	b, err := stream.PeekBytes(4)
	if err != nil {
		return false
	}

	if b[0] == 0x00 && b[1] == 0x00 && b[2] == 0x00 && b[3] == 0x01 {
		*startCodeByteCount = 4
		return true
	} else if b[0] == 0x00 && b[1] == 0x00 && b[2] == 0x01 {
		*startCodeByteCount = 3
		return true	
	}
	return false
}