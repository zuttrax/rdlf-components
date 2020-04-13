package logs

import (
	"encoding/json"
	"time"

	"go.uber.org/zap"
)

type Log struct {
	AppName string
}

type Logger interface {
	Info(string)
	Error(error)
	Fatal(error)
}

func (l Log) Info(info string) {
	logger := initialConfiguration(l.AppName)
	//nolint:errcheck
	defer logger.Sync()
	logger.Info(info)

}

func (l Log) Error(err error) {
	logger := initialConfiguration(l.AppName)
	//nolint:errcheck
	defer logger.Sync()
	logger.Error(err.Error())

}

func (l Log) Fatal(err error) {
	logger := initialConfiguration(l.AppName)
	//nolint:errcheck
	defer logger.Sync()
	logger.Fatal(err.Error())

}

func initialConfiguration(app string) *zap.Logger {
	date := time.Now().Format("01-02-2006")

	rawJSON := []byte(`{
	  "level": "debug",
	  "encoding": "json",
	  "outputPaths": ["stdout", "/rdlf/log-` + app + `-` + date + `.log"],
	  "errorOutputPaths": ["stderr"],
	  "encoderConfig": {
	    "messageKey": "message",
	    "levelKey": "level",
		"levelEncoder": "lowercase",
		"timeKey": "time",
	    "timeEncoder": "ISO8601"
	  }
	}`)

	var cfg zap.Config

	if err := json.Unmarshal(rawJSON, &cfg); err != nil {
		panic(err)
	}

	logger, err := cfg.Build()

	if err != nil {
		panic(err)
	}

	return logger
}
