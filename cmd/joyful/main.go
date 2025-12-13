package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/holoplot/go-evdev"
	flag "github.com/spf13/pflag"

	"git.annabunches.net/annabunches/joyful/internal/configparser"
	"git.annabunches.net/annabunches/joyful/internal/logger"
)

func getConfigDir(dir string) string {
	configDir := strings.ReplaceAll(dir, "~", "${HOME}")
	return os.ExpandEnv(configDir)
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
	config, err := configparser.ParseConfig(configDir)
	logger.FatalIfError(err, "Failed to parse configuration")

	// initialize TTS
	tts, err := newTTS(ttsOps)
	logger.LogIfError(err, "Failed to initialize TTS")

	// Initialize virtual devices with event buffers
	vDevicesByName, vBuffersByName, vBuffersByDevice := initVirtualBuffers(config)

	// Initialize physical devices
	pDevices := initPhysicalDevices(config)

	// initialize the mode variables
	var mode string
	modes := config.Modes
	if len(modes) == 0 {
		mode = "*"
	} else {
		mode = config.Modes[0]
	}

	// Load the rules
	rules, eventChannel, cancel, wg := loadRules(config, pDevices, vDevicesByName, modes)

	// initialize TTS phrases for modes
	if !ttsOps.Disabled {
		for _, m := range modes {
			tts.AddMessage(m)
			logger.LogDebugf("Added TTS message '%s'", m)
		}
	}

	fmt.Println("Joyful Running! Press Ctrl+C to quit. Press Enter to reload rules.")
	if len(modes) > 0 {
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
			config, err := configparser.ParseConfig(configDir) // reload the config
			if err != nil {
				logger.LogError(err, "Failed to parse config, no changes made")
				continue
			}

			fmt.Println("Reloading rules.")
			cancel()
			fmt.Println("Waiting for existing listeners to exit. Provide input from each of your devices.")
			wg.Wait()
			fmt.Println("Listeners exited. Loading new rules.")
			rules, eventChannel, cancel, wg = loadRules(config, pDevices, vDevicesByName, modes)
			fmt.Println("Config re-loaded. Only rule changes applied. Device and Mode changes require restart.")
		}

		if lastMode != mode && tts != nil {
			tts.Say(mode)
		}
	}
}
