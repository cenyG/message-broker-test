package log

import (
	"fmt"
	"time"
)

func Err(format string, error error) {
	msg := fmt.Sprintf(format, error)
	printf("ERROR", msg)
}

func Warn(format string, a ...interface{}) {
	msg := fmt.Sprintf(format, a)
	printf("WARN", msg)
}

func Info(format string, a ...interface{}) {
	msg := fmt.Sprintf(format, a)
	printf("INFO", msg)
}

func printf(logLevel string, msg string) {
	dt := time.Now().Format("2006/02/01 15:04:05")
	fmt.Printf("%s [%s] %s\n", dt, logLevel, msg)
}