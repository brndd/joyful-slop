package main

import (
	"fmt"

	"git.annabunches.net/annabunches/joyful/internal/logger"
	"github.com/holoplot/go-evdev"
)

func main() {
	devices, err := evdev.ListDevicePaths()
	logger.FatalIfError(err, "")

	for _, device := range devices {
		fmt.Printf("%s: '%s'\n", device.Path, device.Name)
	}
}
