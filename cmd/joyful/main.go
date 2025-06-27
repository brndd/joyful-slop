package main

import (
	"fmt"
	"syscall"
	"time"

	"git.annabunches.net/annabunches/joyful/internal/virtualdevice"
	"github.com/holoplot/go-evdev"
)

func main() {
	// Check for and destroy any existing joyful devices
	virtualdevice.CleanupStaleVirtualDevices()

	// STUB: parse virtual device config

	// STUB: parse mapping config

	// Define virtual device
	// TODO: create virtual devices from config
	vDevice, err := evdev.CreateDevice(
		"joyful-joystick",
		evdev.InputID{
			BusType: 0x03,
			Vendor:  0x4711,
			Product: 0x0816,
			Version: 1,
		},
		map[evdev.EvType][]evdev.EvCode{
			evdev.EV_KEY: jsButtons(),
			evdev.EV_ABS: {
				evdev.ABS_X,
				evdev.ABS_Y,
				evdev.ABS_Z,
				evdev.ABS_RX,
				evdev.ABS_RY,
				evdev.ABS_RZ,
				evdev.ABS_THROTTLE,
				evdev.ABS_WHEEL,
			},
		},
	)
	if err != nil {
		fmt.Printf("Failed to create vDevice: %s", err.Error())
	}

	location, err := vDevice.PhysicalLocation()
	if err != nil {
		fmt.Printf("Couldn't get virtual device location: %s\n", err.Error())
	}
	fmt.Printf("Device created as %s. Press Ctrl+C to quit and destroy the device.\n", location)

	var value int32 = 1
	for {
		eventTime := syscall.NsecToTimeval(int64(time.Now().Nanosecond()))

		vDevice.WriteOne(&evdev.InputEvent{
			Time:  eventTime,
			Type:  evdev.EV_KEY,
			Code:  evdev.BTN_TRIGGER,
			Value: value,
		})

		if value == 0 {
			value = 1
		} else {
			value = 0
		}

		time.Sleep(1 * time.Second)
	}
}

func jsButtons() []evdev.EvCode {
	buttons := make([]evdev.EvCode, 80)
	i := 0

	for code := 0x120; code <= 0x12f; code++ {
		buttons[i] = evdev.EvCode(code)
		i++
	}

	for code := 0x2c0; code <= 0x2ff; code++ {
		buttons[i] = evdev.EvCode(code)
		i++
	}

	return buttons
}
