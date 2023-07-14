package src

import logger "log"

var (
	Verbose = false
)

func log(format string, a ...any) {
	if !Verbose {
		return
	}

	logger.SetFlags(logger.Ldate | logger.Ltime)
	logger.Printf(format, a...)
}
