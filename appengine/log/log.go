package log

import (
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
)

func init() {
	log.SetFlags(log.LstdFlags)
}

func Fatalf(format string, args ...interface{}) {
	printf("FATAL", format, args)
	os.Exit(1)
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

func printf(level string, format string, args []interface{}) {
	_, file, line, ok := runtime.Caller(2)
	if ok {
		format = level + ": " + filepath.Base(file) + ":" + strconv.Itoa(line) + ": " + format
	} else {
		format = level + ": " + format
	}

	log.Printf(format, args...)
}
