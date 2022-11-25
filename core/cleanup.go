package core

import (
	"github/LIOU2021/go-eloquent-mongodb/logger"
)

func Cleanup() {
	logger.Close()
}
