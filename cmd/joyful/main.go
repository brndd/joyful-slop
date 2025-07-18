package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"sync"

	"git.annabunches.net/annabunches/joyful/internal/config"
	"git.annabunches.net/annabunches/joyful/internal/logger"
	"git.annabunches.net/annabunches/joyful/internal/mappingrules"
	"git.annabunches.net/annabunches/joyful/internal/virtualdevice"
	"github.com/holoplot/go-evdev"
)

func getConfigDir() string {
	configFlag := flag.String("config", "~/.config/joyful", "Directory to read configuration from.")
	flag.Parse()
	configDir := strings.ReplaceAll(*configFlag, "~", "${HOME}")
	return os.ExpandEnv(configDir)
}

func readConfig(configDir string) *config.ConfigParser {
	parser := &config.ConfigParser{}
	err := parser.Parse(configDir)
	logger.FatalIfError(err, "Failed to parse config")
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
		devices[name] = buffer.Device.(*evdev.InputDevice)
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
	configDir := getConfigDir()
	config := readConfig(configDir)

	// Initialize virtual devices with event buffers
	vBuffersByName, vBuffersByDevice := initVirtualBuffers(config)

	// Initialize physical devices
	pDevices := initPhysicalDevices(config)

	rules, eventChannel, doneChannel, wg := loadRules(config, pDevices, getVirtualDevices(vBuffersByName))

	// initialize the mode variable
	mode := config.GetModes()[0]

	fmt.Println("Joyful Running! Press Ctrl+C to quit. Press Enter to reload rules.")
	if len(config.GetModes()) > 1 {
		logger.Logf("Initial mode set to '%s'", mode)
	}

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
			config := readConfig(configDir) // reload the config
			rules, eventChannel, doneChannel, wg = loadRules(config, pDevices, getVirtualDevices(vBuffersByName))
			fmt.Println("Config re-loaded. Only rule changes applied. Device and Mode changes require restart.")
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
