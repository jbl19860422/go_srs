package main

import (
	"flag"
	"go_srs/srs/app"
)

var (
	port = flag.Int("p", 1935, "set port `port`")
)

func main() {
	flag.Parse()
	server := &app.SrsServer{}
	_ = server.StartProcess(*port)
}
