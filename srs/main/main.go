package main

import (
	"flag"
	"go_srs/srs/app"
	"net/http"
	// "log"
	"bytes"
    "io/ioutil"
    "math/rand"
	_ "net/http/pprof"
)

var (
	port = flag.Int("p", 1935, "set port `port`")
)

func main() {
	// go func() {
	// 	http.HandleFunc("/test", handler)
    // 	log.Fatal(http.ListenAndServe(":9876", nil))
	// }()

	flag.Parse()
	server := app.NewSrsServer()
	_ = server.StartProcess(*port)
}

func handler(w http.ResponseWriter, r *http.Request) {
    err := r.ParseForm()
    if nil != err {
        w.Write([]byte(err.Error()))
        return
    }
    doSomeThingOne(10000)
    buff := genSomeBytes()
    b, err := ioutil.ReadAll(buff)
    if nil != err {
        w.Write([]byte(err.Error()))
        return
    }
    w.Write(b)
}

func doSomeThingOne(times int) {
    for i := 0; i < times; i++ {
        for j := 0; j < times; j++ {

        }
    }
}

func genSomeBytes() *bytes.Buffer {
    var buff bytes.Buffer
    for i := 1; i < 20000; i++ {
        buff.Write([]byte{'0' + byte(rand.Intn(10))})
    }
    return &buff
}
