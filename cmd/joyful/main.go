package main

import (
	"fmt"
	"syscall"
	"time"

	"github.com/holoplot/go-evdev"
)

func main() {
	// Define virtual device
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

		time.Sleep(1000)
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
