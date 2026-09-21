package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/holoplot/go-evdev"
	flag "github.com/spf13/pflag"

	"git.annabunches.net/annabunches/joyful/internal/configparser"
	"git.annabunches.net/annabunches/joyful/internal/logger"
	"git.annabunches.net/annabunches/joyful/internal/mappingrules"
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
	initialMode := mode

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
		suppressModeAnnouncement := false
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
				// All rules see the mode that was active when this input arrived.
				// A requested mode change is applied after dispatch so rule order
				// cannot split one event across two modes.
				outputs, requestedMode, modeChangeRule := matchRules(rules, channelEvent.Device, channelEvent.Event, mode)
				for _, output := range outputs {
					vBuffersByDevice[output.Device].AddEvent(output.Event)
				}
				if modeChangeRule != nil {
					mode = requestedMode
					_, suppressModeAnnouncement = modeChangeRule.(mappingrules.SilentModeChangeRule)
					for _, output := range modeTransitionEvents(rules, mode, modeChangeRule) {
						vBuffersByDevice[output.Device].AddEvent(output.Event)
					}
				}
			}

		case ChannelEventTimer:
			// Evaluate timed rules here, never in the timer goroutine. This keeps mode
			// and rule state changes serialized with physical input processing.
			changedBuffers := make(map[*evdev.InputDevice]struct{})
			timerMode := mode
			for _, output := range channelEvent.Rule.TimerEvents(&timerMode) {
				vBuffersByDevice[output.Device].AddEvent(output.Event)
				changedBuffers[output.Device] = struct{}{}
			}
			if timerMode != mode {
				mode = timerMode
				if rule, ok := channelEvent.Rule.(mappingrules.MappingRule); ok {
					_, suppressModeAnnouncement = rule.(mappingrules.SilentModeChangeRule)
					for _, output := range modeTransitionEvents(rules, mode, rule) {
						vBuffersByDevice[output.Device].AddEvent(output.Event)
						changedBuffers[output.Device] = struct{}{}
					}
				}
			}
			for device := range changedBuffers {
				vBuffersByDevice[device].SendEvents()
			}

		case ChannelEventReload:
			// stop existing channels
			config, err := configparser.ParseConfig(configDir) // reload the config
			if err != nil {
				logger.LogError(err, "Failed to parse config, no changes made")
				continue
			}

			fmt.Println("Reloading rules.")
			for _, output := range resetRuleEvents(rules, &mode, initialMode) {
				vBuffersByDevice[output.Device].AddEvent(output.Event)
			}
			for _, buffer := range vBuffersByName {
				buffer.SendEvents()
			}
			cancel()
			fmt.Println("Waiting for existing listeners to exit. Provide input from each of your devices.")
			wg.Wait()
			fmt.Println("Listeners exited. Loading new rules.")
			rules, eventChannel, cancel, wg = loadRules(config, pDevices, vDevicesByName, modes)
			fmt.Println("Config re-loaded. Active outputs cleared and startup mode restored. Device and Mode configuration changes require restart.")
		}

		if shouldAnnounceModeChange(lastMode, mode, suppressModeAnnouncement) && tts != nil {
			tts.Say(mode)
		}
	}
}

func shouldAnnounceModeChange(previous, current string, suppressed bool) bool {
	return previous != current && !suppressed
}

func matchRules(rules []mappingrules.MappingRule, device mappingrules.Device, event *evdev.InputEvent, mode string) ([]mappingrules.OutputEvent, string, mappingrules.MappingRule) {
	var outputs []mappingrules.OutputEvent
	requestedMode := mode
	var modeChangeRule mappingrules.MappingRule
	for _, rule := range rules {
		if modeChangeRule != nil {
			if modeRule, changesMode := rule.(mappingrules.ModeChangingRule); changesMode && !modeRule.ModeChangeActive() {
				continue
			}
		}
		ruleMode := mode
		if multiRule, ok := rule.(mappingrules.MultiEventMappingRule); ok {
			outputs = append(outputs, multiRule.MatchEvents(device, event, &ruleMode)...)
		} else {
			outputDevice, outputEvent := rule.MatchEvent(device, event, &ruleMode)
			if outputDevice != nil && outputEvent != nil {
				outputs = append(outputs, mappingrules.OutputEvent{Device: outputDevice, Event: outputEvent})
			}
		}
		if modeChangeRule == nil && ruleMode != mode {
			requestedMode = ruleMode
			modeChangeRule = rule
		}
	}
	return outputs, requestedMode, modeChangeRule
}

func modeTransitionEvents(rules []mappingrules.MappingRule, newMode string, initiator mappingrules.MappingRule) []mappingrules.OutputEvent {
	var events []mappingrules.OutputEvent
	for _, rule := range rules {
		if transitionRule, ok := rule.(mappingrules.ModeTransitionRule); ok {
			events = append(events, transitionRule.ModeChanged(newMode, rule == initiator)...)
		}
	}
	return events
}

func resetRuleEvents(rules []mappingrules.MappingRule, mode *string, initialMode string) []mappingrules.OutputEvent {
	var events []mappingrules.OutputEvent
	for _, rule := range rules {
		if resettableRule, ok := rule.(mappingrules.ResettableRule); ok {
			events = append(events, resettableRule.Reset(mode)...)
		}
	}
	*mode = initialMode
	return events
}
