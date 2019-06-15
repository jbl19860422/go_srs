package main

import(
	"fmt"
	"bytes"
	"encoding/binary"
)

func BoolToBytes(data bool) []byte {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.BigEndian, data)
	if err != nil {
		return nil
	}
	return buf.Bytes()
}

func main() {
	b := BoolToBytes(false)
	fmt.Println(b)
}
