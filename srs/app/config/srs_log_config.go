package config

import "errors"

type LogConfig struct {
	logLevel string
	logDir   string
	logFile  string
}

func (this *LogConfig) Parse(root interface{}) error {
	m, ok := root.(map[string]interface{})
	if !ok {
		return errors.New("parse log config failed")
	}

	if this.logLevel, ok = m["level"].(string); ok != true {
		this.logLevel = "debug"
	}

	if this.logDir, ok = m["dir"].(string); ok != true {
		this.logDir = "./"
	}

	if this.logFile, ok = m["file"].(string); ok != true {
		this.logFile = "log.txt"
	}
	return nil
}
