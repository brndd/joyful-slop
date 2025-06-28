package logger

import (
	"fmt"
	"os"
)

func Log(msg string) {
	fmt.Println(msg)
}

func LogIfError(err error, msg string) {
	if err == nil {
		return
	}

	fmt.Printf("%s: %s\n", msg, err.Error())
}

func FatalIfError(err error, msg string) {
	if err == nil {
		return
	}

	LogIfError(err, msg)
	os.Exit(1)
}
