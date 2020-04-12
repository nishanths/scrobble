package log

import (
	"log"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func Criticalf(format string, args ...interface{}) {
	printf("CRITICAL", format, args)
}

func Errorf(format string, args ...interface{}) {
	printf("ERROR", format, args)
}

func Warningf(format string, args ...interface{}) {
	printf("WARNING", format, args)
}

func Infof(format string, args ...interface{}) {
	printf("INFO", format, args)
}

func printf(level string, format string, args ...interface{}) {
	log.Printf(level+": "+format, args...)
}
