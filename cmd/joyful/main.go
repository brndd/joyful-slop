package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

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

func initVirtualBuffers(config *config.ConfigParser) (map[string]*virtualdevice.EventBuffer, map[*evdev.InputDevice]*virtualdevice.EventBuffer) {
	vDevices := config.CreateVirtualDevices()
	if len(vDevices) == 0 {
		logger.Log("Warning: no virtual devices found in configuration. No rules will work.")
	}

	vBuffersByName := make(map[string]*virtualdevice.EventBuffer)
	vBuffersByDevice := make(map[*evdev.InputDevice]*virtualdevice.EventBuffer)
	for name, device := range vDevices {
		vBuffersByName[name] = virtualdevice.NewEventBuffer(device)
		vBuffersByDevice[device] = vBuffersByName[name]
	}
	return vBuffersByName, vBuffersByDevice
}

// Extracts the evdev devices from a list of virtual buffers and returns them.
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
	vBuffersByName, vBuffersByDevice := initVirtualBuffers(config)

	// Initialize physical devices
	pDevices := initPhysicalDevices(config)

	rules, eventChannel, doneChannel, wg := loadRules(config, pDevices, getVirtualDevices(vBuffersByName))

	// initialize the mode variable
	mode := config.GetModes()[0]
	logger.Logf("Initial mode set to '%s'", mode)

	fmt.Println("Joyful Running! Press Ctrl+C to quit.")
	for {
		// Get an event (blocks if necessary)
		channelEvent := <-eventChannel

		switch channelEvent.Type {
		case ChannelEventInput:
			switch channelEvent.Event.Type {
			case evdev.EV_SYN:
				// We've received a SYN_REPORT, so now we send all pending events; since SYN_REPORTs
				// might come from multiple input devices, we'll always flush, just to be sure.
				for _, buffer := range vBuffersByName {
					buffer.SendEvents()
				}

			case evdev.EV_KEY, evdev.EV_ABS:
				// We have a matchable event type. Check all the events
				for _, rule := range rules {
					device, outputEvent := rule.MatchEvent(channelEvent.Device, channelEvent.Event, &mode)
					if device == nil || outputEvent == nil {
						continue
					}
					vBuffersByDevice[device].AddEvent(outputEvent)
				}
			}

		case ChannelEventTimer:
			// Timer events give us the device and event to use directly
			vBuffersByDevice[channelEvent.Device].AddEvent(channelEvent.Event)
			// If we get a timer event, flush the output device buffer immediately
			vBuffersByDevice[channelEvent.Device].SendEvents()

		case ChannelEventReload:
			// stop existing channels
			fmt.Println("Reloading rules.")
			doneChannel <- true
			fmt.Println("Waiting for existing listeners to exit. Provide input from each of your devices.")
			wg.Wait()
			fmt.Println("Listeners exited. Parsing config.")
			config := readConfig() // reload the config
			rules, eventChannel, doneChannel, wg = loadRules(config, pDevices, getVirtualDevices(vBuffersByName))
		}
	}
}

func loadRules(
	config *config.ConfigParser,
	pDevices map[string]*evdev.InputDevice,
	vDevices map[string]*evdev.InputDevice) ([]mappingrules.MappingRule, <-chan ChannelEvent, chan bool, *sync.WaitGroup) {

	var wg sync.WaitGroup
	eventChannel := make(chan ChannelEvent, 1000)
	doneChannel := make(chan bool)

	// Initialize rules
	rules := config.BuildRules(pDevices, vDevices)
	logger.Logf("Created %d mapping rules.", len(rules))

	// start listening for events on devices and timers
	for _, device := range pDevices {
		wg.Add(1)
		go eventWatcher(device, eventChannel, doneChannel, &wg)
	}

	timerCount := 0
	for _, rule := range rules {
		if timedRule, ok := rule.(mappingrules.TimedEventEmitter); ok {
			wg.Add(1)
			go timerWatcher(timedRule, eventChannel, doneChannel, &wg)
			timerCount++
		}
	}
	logger.Logf("registered %d timers", timerCount)

	go consoleWatcher(eventChannel, &wg)

	return rules, eventChannel, doneChannel, &wg
}
