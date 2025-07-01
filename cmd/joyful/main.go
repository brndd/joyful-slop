package main

import (
	"fmt"
	"maps"
	"os"
	"path/filepath"
	"slices"
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
	if len(vDevices) == 0 {
		logger.Log("Warning: no virtual devices found in configuration. No rules will work.")
	}

	vBuffers := make(map[string]*virtualdevice.EventBuffer)
	for name, device := range vDevices {
		vBuffers[name] = virtualdevice.NewEventBuffer(device)
	}
	return vBuffers
}

func initPhysicalDevices(config *config.ConfigParser) map[string]*evdev.InputDevice {
	pDeviceMap := config.ConnectPhysicalDevices()
	if len(pDeviceMap) == 0 {
		logger.Log("Warning: no physical devices found in configuration. No rules will work.")
	}
	return pDeviceMap
}

func main() {
	// parse configs
	config := readConfig()

	// Initialize virtual devices and event buffers
	vBuffers := initVirtualDevices(config)

	// Initialize physical devices
	pDevices := initPhysicalDevices(config)

	// TEST CODE
	testDriver(vBuffers, pDevices)
}

func testDriver(vBuffers map[string]*virtualdevice.EventBuffer, pDevices map[string]*evdev.InputDevice) {
	pDevice := slices.Collect(maps.Values(pDevices))[0]

	name, err := pDevice.Name()
	if err != nil {
		name = "Unknown"
	}
	fmt.Printf("Test Driver using physical device %s\n", name)

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
						combo++
					}
					if event.Value == 1 {
						combo--
					}

				case evdev.BTN_THUMB:
					if event.Value == 0 {
						combo--
					}
					if event.Value == 1 {
						combo++
					}

				case evdev.BTN_THUMB2:
					if event.Value == 0 {
						combo--
					}
					if event.Value == 1 {
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
