package log

import (
	"fmt"
	"log"
	"path/filepath"
	"runtime"
	"strconv"
)

func init() {
	log.SetFlags(log.LstdFlags)
}

func Fatalf(format string, args ...interface{}) {
	printf("FATAL", format, args)
	panic(fmt.Sprintf(format, args...))
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
		format = filepath.Base(file) + ":" + strconv.Itoa(line) + ": " + level + ": " + format
	} else {
		format = level + ": " + format
	}

	log.Printf(format, args...)
}
