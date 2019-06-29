package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

/*
{
	"hls":{
		"record":{
			"path":"aaaa",
			"time":12345
		},
		"public":false,
		"url_refix":"http://172.19.5.107/"
	}
}
*/

type RecordConf struct {
	Path string `json:"path"`
	time int    `json:"time"`
}

type HlsConf struct {
	Record    RecordConf `json:"record"`
	Public    bool       `json:"public"`
	UrlPrefix string     `json:"url_prefix"`
}

type Config struct {
	Hls *HlsConf `json:"hls"`
}

func main() {
	var c Config
	data, err := ioutil.ReadFile("conf.json")
	if err != nil {
		return
	}

	err = json.Unmarshal(data, &c)
	_ = err

	fmt.Println("paht=", c.Hls.Record.Path)
	if err != nil {
		fmt.Println("json decode failed")
	} else {
		fmt.Println("json decode succeed")
	}
}
