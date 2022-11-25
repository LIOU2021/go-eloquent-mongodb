package core

import (
	"github.com/LIOU2021/go-eloquent-mongodb/logger"
)

func Cleanup() {
	logger.Close()
}
