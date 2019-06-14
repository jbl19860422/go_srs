package main

import(
	"fmt"
	"bytes"
	"encoding/binary"
)

func UInt16ToBytes(data uint16, order binary.ByteOrder) []byte {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, order, data)
	if err != nil {
		return nil
	}
	return buf.Bytes()
}

func NumberToBytes(data interface{}, order binary.ByteOrder) []byte {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, order, data)
	if err != nil {
		return nil
	}
	return buf.Bytes()
}

func main() {
	var d int16 = 18
	b := NumberToBytes(d, binary.LittleEndian)
	fmt.Println(b)
}
