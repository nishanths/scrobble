package log

import (
	"context"
	"log"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func Criticalf(_ context.Context, format string, args ...interface{}) {
	printf("CRITICAL", format, args)
}

func Errorf(_ context.Context, format string, args ...interface{}) {
	printf("ERROR", format, args)
}

func Warningf(_ context.Context, format string, args ...interface{}) {
	printf("WARNING", format, args)
}

func Infof(_ context.Context, format string, args ...interface{}) {
	printf("INFO", format, args)
}

func printf(level string, format string, args ...interface{}) {
	log.Printf(level+": "+format, args...)
}
