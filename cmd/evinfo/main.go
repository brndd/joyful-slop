package main

import (
	"fmt"
	"slices"

	// TODO: using config here feels like bad coupling... ButtonFromIndex might need a refactor / move
	"git.annabunches.net/annabunches/joyful/internal/config"
	"git.annabunches.net/annabunches/joyful/internal/logger"
	"github.com/holoplot/go-evdev"
	flag "github.com/spf13/pflag"
)

func isJoystickLike(device *evdev.InputDevice) bool {
	types := device.CapableTypes()
	if slices.Contains(types, evdev.EV_ABS) {
		return true
	}

	if slices.Contains(types, evdev.EV_KEY) {
		buttons := device.CapableEvents(evdev.EV_KEY)

		for _, code := range config.ButtonFromIndex {
			if slices.Contains(buttons, code) {
				return true
			}
		}
	}

	return false
}

func printDevice(devPath evdev.InputPath) {
	// Get the device
	device, err := evdev.Open(devPath.Path)
	if err != nil {
		fmt.Printf("ERROR: Couldn't get data for device '%s'\n", devPath.Name)
		return
	}

	// If the device doesn't support any joystick-shaped inputs, don't print it.
	if !isJoystickLike(device) {
		return
	}

	// Get axis info
	var axisOutputs []string
	absInfos, err := device.AbsInfos()
	if err != nil {
		axisOutputs = []string{"ERROR: Failed to get axis data"}
	} else {
		axisOutputs = make([]string, 0, len(absInfos))
		for axisCode, info := range absInfos {
			absStr := fmt.Sprintf("%s: %d - %d", evdev.ABSNames[axisCode], info.Minimum, info.Maximum)
			axisOutputs = append(axisOutputs, absStr)
		}
	}

	fmt.Printf("%s:\n", devPath.Path)
	fmt.Printf("\tName: '%s'\n", devPath.Name)
	if len(axisOutputs) > 0 {
		fmt.Println("\tAxes:")
		for _, str := range axisOutputs {
			fmt.Printf("\t\t%s\n", str)
		}
	}
}

func printDeviceQuiet(devPath evdev.InputPath) {
	// If it's not at least minimally joystick-shaped, skip.
	// If we can't open the device, err on the side of printing it
	// anyway
	device, err := evdev.Open(devPath.Path)
	if err != nil ||
		!isJoystickLike(device) {
		return
	}

	fmt.Printf("'%s'\n", devPath.Name)
}

// TODO: it would be nice to be able to specify a device by name or device file and get axis info
// just for that device
func main() {
	var quietFlag bool
	flag.BoolVarP(&quietFlag, "quiet", "q", false, "Only print device names")
	flag.Parse()

	devices, err := evdev.ListDevicePaths()
	logger.FatalIfError(err, "")

	for _, device := range devices {
		if quietFlag {
			printDeviceQuiet(device)
		} else {
			printDevice(device)
		}
	}
}
