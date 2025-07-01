package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"git.annabunches.net/annabunches/joyful/internal/logger"
	"github.com/goccy/go-yaml"
	"github.com/holoplot/go-evdev"
)

type ConfigParser struct {
	config Config
}

// Parse all the config files and store the config data for further use
func (parser *ConfigParser) Parse(directory string) error {
	parser.config = Config{}

	// Find the config files in the directory
	dirEntries, err := os.ReadDir(directory)
	if err != nil {
		err = os.Mkdir(directory, 0755)
		if err != nil {
			return errors.New("Failed to create config directory at " + directory)
		}
	}

	// Open each yaml file and add its contents to the global config
	for _, file := range dirEntries {
		name := file.Name()
		if file.IsDir() || !(strings.HasSuffix(name, ".yaml") || strings.HasSuffix(name, ".yml")) {
			continue
		}

		filePath := filepath.Join(directory, name)
		if strings.HasSuffix(filePath, ".yaml") || strings.HasSuffix(filePath, ".yml") {
			data, err := os.ReadFile(filePath)
			if err != nil {
				logger.LogError(err, "Error while opening config file")
				continue
			}
			newConfig := Config{}
			err = yaml.Unmarshal(data, &newConfig)
			logger.LogIfError(err, "Error parsing YAML")
			parser.config.Rules = append(parser.config.Rules, newConfig.Rules...)
			parser.config.Devices = append(parser.config.Devices, newConfig.Devices...)
			// parser.config.Groups = append(parser.config.Groups, newConfig.Groups...)
		}
	}

	if len(parser.config.Devices) == 0 {
		return errors.New("Found no devices in configuration. Please add configuration at " + directory)
	}

	return nil
}

// CreateVirtualDevices will register any configured devices with type = virtual
// using /dev/uinput, and return a map of those devices.
//
// This function assumes you have already called Parse() on the config directory.
//
// This function should only be called once, unless you want to create duplicate devices for some reason.
func (parser *ConfigParser) CreateVirtualDevices() map[string]*evdev.InputDevice {
	deviceMap := make(map[string]*evdev.InputDevice)

	for _, deviceConfig := range parser.config.Devices {
		if strings.ToLower(deviceConfig.Type) != DeviceTypeVirtual {
			continue
		}

		vDevice, err := evdev.CreateDevice(
			fmt.Sprintf("joyful-%s", deviceConfig.Name),
			// TODO: who knows what these should actually be
			evdev.InputID{
				BusType: 0x03,
				Vendor:  0x4711,
				Product: 0x0816,
				Version: 1,
			},
			map[evdev.EvType][]evdev.EvCode{
				evdev.EV_KEY: makeButtons(int(deviceConfig.Buttons)),
				evdev.EV_ABS: makeAxes(int(deviceConfig.Axes)),
			},
		)

		if err != nil {
			logger.LogIfError(err, "Failed to create virtual device")
			continue
		}

		deviceMap[deviceConfig.Name] = vDevice
	}

	return deviceMap
}

// ConnectPhysicalDevices will create InputDevices corresponding to any registered
// devices with type = physical. It will also attempt to acquire exclusive access
// to those devices, to prevent the same inputs from being read on multiple devices.
//
// This function assumes you have already called Parse() on the config directory.
//
// This function should only be called once.
//
// STUB: this function does not yet function.
func (parser *ConfigParser) ConnectPhysicalDevices() map[string]*evdev.InputDevice {
	deviceMap := make(map[string]*evdev.InputDevice)

	for _, deviceConfig := range parser.config.Devices {
		if strings.ToLower(deviceConfig.Type) != DeviceTypePhysical {
			continue
		}

		vDevice, err := evdev.Open("/dev/input/foo")

		if err != nil {
			logger.LogIfError(err, "Failed to open physical device")
			continue
		}

		deviceMap[deviceConfig.Name] = vDevice
	}

	return deviceMap
}

func makeButtons(numButtons int) []evdev.EvCode {
	if numButtons > 56 {
		numButtons = 56
		logger.Log("Limiting virtual device buttons to 56")
	}

	buttons := make([]evdev.EvCode, numButtons)

	startCode := 0x120
	for i := 0; i < numButtons && i < 16; i++ {
		buttons[i] = evdev.EvCode(startCode + i)
	}

	if numButtons > 16 {
		startCode = 0x2c0
		for i := 0; i < numButtons-16; i++ {
			buttons[16+i] = evdev.EvCode(startCode + i)
		}
	}

	return buttons
}

func makeAxes(numAxes int) []evdev.EvCode {
	if numAxes > 8 {
		numAxes = 8
		logger.Log("Limiting virtual device axes to 8")
	}

	axes := make([]evdev.EvCode, numAxes)
	for i := 0; i < numAxes; i++ {
		axes[i] = evdev.EvCode(i)
	}

	return axes
}
