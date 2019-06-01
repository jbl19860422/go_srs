package main

import (
	"net"
	// "fmt"
	_ "time"
)

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:1935")
	if err != nil {
		// handle error
	}

	str := []byte("test client")
	conn.Write(str)
	conn.Close();
}