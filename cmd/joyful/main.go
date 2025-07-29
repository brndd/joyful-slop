package main

import (
	"context"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/holoplot/go-evdev"
	flag "github.com/spf13/pflag"

	"git.annabunches.net/annabunches/joyful/internal/config"
	"git.annabunches.net/annabunches/joyful/internal/logger"
	"git.annabunches.net/annabunches/joyful/internal/mappingrules"
	"git.annabunches.net/annabunches/joyful/internal/virtualdevice"
)

func getConfigDir(dir string) string {
	configDir := strings.ReplaceAll(dir, "~", "${HOME}")
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
	// parse command-line
	var configFlag string
	flag.BoolVarP(&logger.IsDebugMode, "debug", "d", false, "Output very verbose debug messages.")
	flag.StringVarP(&configFlag, "config", "c", "~/.config/joyful", "Directory to read configuration from.")
	ttsOps := addTTSFlags()
	flag.Parse()

	// parse configs
	configDir := getConfigDir(configFlag)
	config := readConfig(configDir)

	// initialize TTS
	tts, err := newTTS(ttsOps)
	logger.LogIfError(err, "Failed to initialize TTS")

	// Initialize virtual devices with event buffers
	vBuffersByName, vBuffersByDevice := initVirtualBuffers(config)

	// Initialize physical devices
	pDevices := initPhysicalDevices(config)

	// Load the rules
	rules, eventChannel, cancel, wg := loadRules(config, pDevices, getVirtualDevices(vBuffersByName))

	// initialize the mode variable
	mode := config.GetModes()[0]

	// initialize TTS phrases for modes
	for _, m := range config.GetModes() {
		tts.AddMessage(m)
		logger.LogDebugf("Added TTS message '%s'", m)
	}

	fmt.Println("Joyful Running! Press Ctrl+C to quit. Press Enter to reload rules.")
	if len(config.GetModes()) > 1 {
		logger.Logf("Initial mode set to '%s'", mode)
	}

	for {
		lastMode := mode
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
			cancel()
			fmt.Println("Waiting for existing listeners to exit. Provide input from each of your devices.")
			wg.Wait()
			fmt.Println("Listeners exited. Parsing config.")
			config := readConfig(configDir) // reload the config
			rules, eventChannel, cancel, wg = loadRules(config, pDevices, getVirtualDevices(vBuffersByName))
			fmt.Println("Config re-loaded. Only rule changes applied. Device and Mode changes require restart.")
		}

		if lastMode != mode && tts != nil {
			tts.Say(mode)
		}
	}
}

func loadRules(
	config *config.ConfigParser,
	pDevices map[string]*evdev.InputDevice,
	vDevices map[string]*evdev.InputDevice) ([]mappingrules.MappingRule, <-chan ChannelEvent, func(), *sync.WaitGroup) {

	var wg sync.WaitGroup
	eventChannel := make(chan ChannelEvent, 1000)
	ctx, cancel := context.WithCancel(context.Background())

	// Initialize rules
	rules := config.BuildRules(pDevices, vDevices)
	logger.Logf("Created %d mapping rules.", len(rules))

	// start listening for events on devices and timers
	for _, device := range pDevices {
		wg.Add(1)
		go eventWatcher(device, eventChannel, ctx, &wg)
	}

	timerCount := 0
	for _, rule := range rules {
		if timedRule, ok := rule.(mappingrules.TimedEventEmitter); ok {
			wg.Add(1)
			go timerWatcher(timedRule, eventChannel, ctx, &wg)
			timerCount++
		}
	}
	logger.Logf("Registered %d timers.", timerCount)

	go consoleWatcher(eventChannel)

	return rules, eventChannel, cancel, &wg
}
