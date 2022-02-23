package common

import (
	"io/ioutil"

	"github.com/google/logger"
)

// init Init logger
func init() {
	var (
		name      = "log.txt"
		verbose   = true
		systemLog = false
	)
	logger.Init(name, verbose, systemLog, ioutil.Discard)
}
