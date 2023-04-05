package logger

import (
	"io/ioutil"

	"github.com/google/logger"
)

var LogDebug *logger.Logger

func Init() {
	if LogDebug == nil {
		LogDebug = logger.Init("DebugLog", true, false, ioutil.Discard)
	}
}

func Close() {
	if LogDebug != nil {
		LogDebug.Close()
	}
}
