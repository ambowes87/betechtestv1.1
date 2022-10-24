package logger

import (
	"fmt"
	"time"
)

const defaultFormat = "%s: %s\r\n"

func Log(msg string) {
	fmt.Printf(defaultFormat, time.Now().Format(time.Stamp), msg)
}
