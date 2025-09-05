package main

import (
	"context"
	"sync"

	"git.annabunches.net/annabunches/joyful/internal/configparser"
	"git.annabunches.net/annabunches/joyful/internal/logger"
	"git.annabunches.net/annabunches/joyful/internal/mappingrules"
	"git.annabunches.net/annabunches/joyful/internal/virtualdevice"
	"github.com/holoplot/go-evdev"
)

func initPhysicalDevices(conf *configparser.Config) map[string]*evdev.InputDevice {
	pDeviceMap := make(map[string]*evdev.InputDevice)

	for _, devConfig := range conf.Devices {
		if devConfig.Type != configparser.DeviceTypePhysical {
			continue
		}

		innerConfig := devConfig.Config.(configparser.DeviceConfigPhysical)
		name, device, err := initPhysicalDevice(innerConfig)
		if err != nil {
			logger.LogError(err, "Failed to initialize physical device")
			continue
		}

		pDeviceMap[name] = device

		displayName := innerConfig.DeviceName
		if innerConfig.DevicePath != "" {
			displayName = innerConfig.DevicePath
		}
		logger.Logf("Connected to '%s' as '%s'", displayName, name)
	}

	if len(pDeviceMap) == 0 {
		logger.Log("Warning: no physical devices found in configuration. No rules will work.")
	}
	return pDeviceMap
}

func initPhysicalDevice(config configparser.DeviceConfigPhysical) (string, *evdev.InputDevice, error) {
	name := config.Name
	var device *evdev.InputDevice
	var err error

	if config.DevicePath != "" {
		device, err = evdev.Open(config.DevicePath)
	} else {
		device, err = evdev.OpenByName(config.DeviceName)
	}

	if config.Lock && err == nil {
		grabErr := device.Grab()
		logger.LogIfError(grabErr, "Failed to lock device for exclusive access")
	}

	return name, device, err
}

// TODO: juggling all these maps is a pain. Is there a better solution here?
func initVirtualBuffers(config *configparser.Config) (map[string]*evdev.InputDevice,
	map[string]*virtualdevice.EventBuffer,
	map[*evdev.InputDevice]*virtualdevice.EventBuffer) {

	vDevicesByName := make(map[string]*evdev.InputDevice)
	vBuffersByName := make(map[string]*virtualdevice.EventBuffer)
	vBuffersByDevice := make(map[*evdev.InputDevice]*virtualdevice.EventBuffer)

	for _, devConfig := range config.Devices {
		if devConfig.Type != configparser.DeviceTypeVirtual {
			continue
		}

		vConfig := devConfig.Config.(configparser.DeviceConfigVirtual)
		buffer, err := virtualdevice.NewEventBuffer(vConfig)
		if err != nil {
			logger.LogError(err, "Failed to create virtual device, skipping")
			continue
		}
		vDevicesByName[buffer.Name] = buffer.Device.(*evdev.InputDevice)
		vBuffersByName[buffer.Name] = buffer
		vBuffersByDevice[buffer.Device.(*evdev.InputDevice)] = buffer
	}

	if len(vDevicesByName) == 0 {
		logger.Log("Warning: no virtual devices found in configuration. No rules will work.")
	}

	return vDevicesByName, vBuffersByName, vBuffersByDevice
}

// TODO: At some point it would *very likely* make sense to map each rule to all of the physical devices that can
// trigger it, and return that instead. Something like a map[Device][]mappingrule.MappingRule.
// This would speed up rule matching by only checking relevant rules for a given input event.
// We could take this further and make it a map[<struct of *inputdevice, type, and code>][]rule
// For very large rule-bases this may be helpful for staying performant.
func loadRules(
	config *configparser.Config,
	pDevices map[string]*evdev.InputDevice,
	vDevices map[string]*evdev.InputDevice,
	modes []string) ([]mappingrules.MappingRule, <-chan ChannelEvent, func(), *sync.WaitGroup) {

	var wg sync.WaitGroup
	eventChannel := make(chan ChannelEvent, 1000)
	ctx, cancel := context.WithCancel(context.Background())

	// Setup device mapping for the mappingrules package
	pDevs := mappingrules.ConvertDeviceMap(pDevices)
	vDevs := mappingrules.ConvertDeviceMap(vDevices)

	// Initialize rules
	rules := make([]mappingrules.MappingRule, 0)
	for _, ruleConfig := range config.Rules {
		newRule, err := mappingrules.NewRule(ruleConfig, pDevs, vDevs, modes)
		if err != nil {
			logger.LogError(err, "Failed to create rule, skipping")
			continue
		}
		rules = append(rules, newRule)
	}

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
