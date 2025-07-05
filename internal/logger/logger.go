package logger

import (
	"fmt"
	"os"
)

func Log(msg string) {
	fmt.Println(msg)
}

func Logf(msg string, params ...interface{}) {
	fmt.Printf(msg+"\n", params...)
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
		fmt.Printf("%s: %s\n", err.Error(), fmt.Sprintf(msg, params...))
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
