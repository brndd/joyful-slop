package main

import (
	"fmt"
	"os"
	"time"

	"git.annabunches.net/annabunches/joyful/internal/virtualdevice"
	"github.com/holoplot/go-evdev"
)

func main() {
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
	fatalIfError(err, "Failed to create virtual device")

	buffer := virtualdevice.NewEventBuffer(vDevice)

	name, err := vDevice.Name()
	if err != nil {
		name = "Unknown"
	}
	fmt.Printf("Virtual device created as %s.\n", name)

	pDevice, err := evdev.Open("/dev/input/event12")
	fatalIfError(err, "Couldn't open physical device")

	name, err = pDevice.Name()
	if err != nil {
		name = "Unknown"
	}
	fmt.Printf("Connected to physical device %s\n", name)

	var combo int32 = 0

	for {
		last := combo

		event, err := pDevice.ReadOne()
		logIfError(err, "Error while reading event")

		// FIXME: test code
		for event.Code != evdev.SYN_REPORT {
			if event.Type == evdev.EV_KEY {
				switch event.Code {
				case evdev.BTN_TRIGGER:
					if event.Value == 0 {
						fmt.Println("Trigger 0")
						combo++
					}
					if event.Value == 1 {
						fmt.Println("Trigger 1")
						combo--
					}

				case evdev.BTN_THUMB:
					if event.Value == 0 {
						fmt.Println("Thumb 0")
						combo--
					}
					if event.Value == 1 {
						fmt.Println("Thumb 1")
						combo++
					}

				case evdev.BTN_THUMB2:
					if event.Value == 0 {
						fmt.Println("Thumb2 0")
						combo--
					}
					if event.Value == 1 {
						fmt.Println("Thumb2 1")
						combo++
					}

				}
			}

			event, err = pDevice.ReadOne()
			logIfError(err, "Error while reading event")
		}

		if combo > last && combo == 3 {
			buffer.AddEvent(&evdev.InputEvent{
				Type:  evdev.EV_KEY,
				Code:  evdev.BTN_TRIGGER,
				Value: 1,
			})
		}
		if combo < last && combo == 2 {
			buffer.AddEvent(&evdev.InputEvent{
				Type:  evdev.EV_KEY,
				Code:  evdev.BTN_TRIGGER,
				Value: 0,
			})
		}

		buffer.SendEvents()
		// FIXME: end test code

		time.Sleep(1 * time.Millisecond)
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

func logIfError(err error, msg string) {
	if err == nil {
		return
	}

	fmt.Printf("%s: %s\n", msg, err.Error())
}

func fatalIfError(err error, msg string) {
	if err == nil {
		return
	}

	logIfError(err, msg)
	os.Exit(1)
}
