package logger

import (
	"fmt"
	"os"
)

var IsDebugMode = false

func Log(msg string) {
	fmt.Println(msg)
}

func Logf(msg string, params ...interface{}) {
	fmt.Printf(msg+"\n", params...)
}

func LogDebugf(msg string, params ...interface{}) {
	if IsDebugMode {
		fmt.Printf("DEBUG: %s\n", fmt.Sprintf(msg, params...))
	}
}

func LogError(err error, msg string) {
	if msg == "" {
		fmt.Printf("%s\n", err.Error())
	} else {
		fmt.Printf("%s: %s\n", msg, err.Error())
	}
}

func LogErrorf(err error, msg string, params ...interface{}) {
	if msg == "" {
		fmt.Printf("%s\n", err.Error())
	} else {
		fmt.Printf("%s: %s\n", fmt.Sprintf(msg, params...), err.Error())
	}
}

func LogIfError(err error, msg string) {
	if err == nil {
		return
	}

	LogError(err, msg)
}

func FatalIfError(err error, msg string) {
	if err == nil {
		return
	}

	LogError(err, msg)
	os.Exit(1)
}
