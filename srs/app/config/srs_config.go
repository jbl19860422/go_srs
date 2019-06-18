package config

import (
	"encoding/json"
	"errors"
	"os"
)

type SrsConfig struct {
	listenPort int
	chunkSize  int
	logConfig  LogConfig
}

var config *SrsConfig

func GetInstance() *SrsConfig {
	if config == nil {
		config = &SrsConfig{}
	}
	return config
}

func (this *SrsConfig) ParseFile(file string) error {
	filePtr, err := os.Open(file)
	if err != nil {
		return err
	}

	defer filePtr.Close()
	decoder := json.NewDecoder(filePtr)

	configData := make(map[string]interface{})
	err = decoder.Decode(&configData)
	if err != nil {
		return err
	}

	var ok bool
	for k, v := range configData {
		switch k {
		case "listen":
			{
				this.listenPort, ok = v.(int)
				if !ok {
					return errors.New("listen port not int type")
				}
			}
		case "chunk_size":
			{
				this.chunkSize, ok = v.(int)
				if !ok {
					return errors.New("chunk size not int type")
				}
			}
		case "log":
			{
				if err = this.logConfig.Parse(v); err != nil {
					return err
				}
			}
		
		}
	}
	return nil
}

func init() {

}
