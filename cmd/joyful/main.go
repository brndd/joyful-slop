package main

import (
	"fmt"
	"os"
	"path/filepath"

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

	// Listen for events and map them forever
	mapEvents(vBuffers, pDevices, rules)
}

type ChannelEvent struct {
	Device *evdev.InputDevice
	Event  *evdev.InputEvent
}

func mapEvents(vBuffers map[string]*virtualdevice.EventBuffer, pDevices map[string]*evdev.InputDevice, rules []mappingrules.MappingRule) {
	// start listening for events on all devices
	eventChannel := make(chan *ChannelEvent, 1000)
	for _, device := range pDevices {
		go eventWatcher(device, eventChannel)
	}

	fmt.Println("Joyful Running! Press Ctrl+C to quit.")
	for {
		// Get an event (blocks if necessary)
		wrapper := <-eventChannel

		switch wrapper.Event.Type {
		case evdev.EV_SYN:
			// We've received a SYN_REPORT, so now we send all of our pending events
			for _, buffer := range vBuffers {
				buffer.SendEvents()
			}

		case evdev.EV_KEY:
		case evdev.EV_ABS:
			// We have a matchable event type. Check all the events
			for _, rule := range rules {
				outputEvent := rule.MatchEvent(wrapper.Device, wrapper.Event)
				if outputEvent == nil {
					continue
				}

				vBuffers[rule.OutputName()].AddEvent(outputEvent)
			}
		}
	}
}

func eventWatcher(device *evdev.InputDevice, channel chan *ChannelEvent) {
	for {
		event, err := device.ReadOne()
		if err != nil {
			logger.LogError(err, "Error while reading event")
			continue
		}
		channel <- &ChannelEvent{Device: device, Event: event}
	}
}
