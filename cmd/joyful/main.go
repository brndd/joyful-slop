package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"git.annabunches.net/annabunches/joyful/internal/config"
	"git.annabunches.net/annabunches/joyful/internal/logger"
	"git.annabunches.net/annabunches/joyful/internal/virtualdevice"
	"github.com/holoplot/go-evdev"
)

func readConfig() *config.ConfigParser {
	parser := &config.ConfigParser{}
	homeDir, err := os.UserHomeDir()
	logger.FatalIfError(err, "Can't get user home directory, so can't find configuration.")
	err = parser.Parse(filepath.Join(homeDir, ".config/joyful"))
	logger.FatalIfError(err, "")
	return parser
}

func initVirtualDevices(config *config.ConfigParser) map[string]*virtualdevice.EventBuffer {
	vDevices := config.CreateVirtualDevices()
	vBuffers := make(map[string]*virtualdevice.EventBuffer)
	for name, device := range vDevices {
		vBuffers[name] = virtualdevice.NewEventBuffer(device)
	}
	return vBuffers
}

func main() {
	// parse configs
	config := readConfig()

	// Initialize virtual devices and event buffers
	vBuffers := initVirtualDevices(config)

	// Initialize physical devices
	// pDevices := config.ConnectPhysicalDevices()

	// TEST CODE
	pDevice, err := evdev.Open("/dev/input/event11")
	logger.FatalIfError(err, "Couldn't open physical device")

	name, err := pDevice.Name()
	if err != nil {
		name = "Unknown"
	}
	fmt.Printf("Connected to physical device %s\n", name)

	var combo int32 = 0
	buffer := vBuffers["main"]
	for {
		last := combo

		event, err := pDevice.ReadOne()
		logger.LogIfError(err, "Error while reading event")

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
			logger.LogIfError(err, "Error while reading event")
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

		time.Sleep(1 * time.Millisecond)
	}
	// END TEST CODE
}
