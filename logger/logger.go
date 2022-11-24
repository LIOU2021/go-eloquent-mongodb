package logger

import (
	"io/ioutil"

	"github.com/google/logger"
)

var LogDebug *logger.Logger

func Init() {
	LogDebug = logger.Init("DebugLog", true, false, ioutil.Discard)
}

func Close() {
	LogDebug.Close()
}
