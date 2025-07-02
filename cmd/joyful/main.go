package main

import (
	"os"
	"path/filepath"
	"time"

	"git.annabunches.net/annabunches/joyful/internal/config"
	"git.annabunches.net/annabunches/joyful/internal/logger"
	"git.annabunches.net/annabunches/joyful/internal/mappingrules"
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

func initVirtualBuffers(config *config.ConfigParser) map[string]*virtualdevice.EventBuffer {
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

func getVirtualDevices(buffers map[string]*virtualdevice.EventBuffer) map[string]*evdev.InputDevice {
	devices := make(map[string]*evdev.InputDevice)
	for name, buffer := range buffers {
		devices[name] = buffer.Device
	}
	return devices
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

	// Initialize virtual devices with event buffers
	vBuffers := initVirtualBuffers(config)

	// Initialize physical devices
	pDevices := initPhysicalDevices(config)

	// Initialize rules
	rules := config.BuildRules(pDevices, getVirtualDevices(vBuffers))

	// TEST CODE
	testDriver(vBuffers, pDevices, rules)
}

func testDriver(vBuffers map[string]*virtualdevice.EventBuffer, pDevices map[string]*evdev.InputDevice, rules []mappingrules.MappingRule) {
	pDevice := pDevices["right-stick"]
	buffer := vBuffers["main"]
	for {
		// Get the first event for this report
		event, err := pDevice.ReadOne()
		logger.LogIfError(err, "Error while reading event")

		for event.Code != evdev.SYN_REPORT {
			for _, rule := range rules {
				event := rule.MatchEvent(pDevice, event)
				if event == nil {
					continue
				}

				buffer.AddEvent(event)
			}

			// Get the next event
			event, err = pDevice.ReadOne()
			logger.LogIfError(err, "Error while reading event")
		}

		// We've received a SYN_REPORT, so now we can send all of our events
		// TODO: how shall we handle this when dealing with multiple devices?
		buffer.SendEvents()

		time.Sleep(1 * time.Millisecond)
	}
	// END TEST CODE
}
