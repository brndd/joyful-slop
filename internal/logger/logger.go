package logger

import (
	"fmt"
	"os"
)

func Log(msg string) {
	fmt.Println(msg)
}

func Logf(msg string, params ...interface{}) {
	fmt.Printf(msg, params...)
}

func LogError(err error, msg string) {
	if msg == "" {
		fmt.Printf("%s\n", err.Error())
	} else {
		fmt.Printf("%s: %s\n", msg, err.Error())
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
