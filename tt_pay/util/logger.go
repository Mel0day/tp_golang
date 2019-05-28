package util

import (
	"fmt"
	"log"
)

var (
	DebugMode	bool
)

func init() {
	DebugMode = false
}

func SetDebugMode(b bool) {
	DebugMode = b
	fmt.Println("Set to Debug Mode: ", DebugMode)
}

func Debug(format string, msg ...interface{}) {
	if DebugMode {
		log.Printf(format, msg...)
	}
}