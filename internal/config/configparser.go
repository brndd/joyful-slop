package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"git.annabunches.net/annabunches/joyful/internal/logger"
	"github.com/goccy/go-yaml"
	"github.com/holoplot/go-evdev"
)

type ConfigParser struct {
	config      Config
	configFiles []string
}

// Parse all the config files and store the config data for further use
func (parser *ConfigParser) Parse(directory string) {
	// Find the config files in the directory
	dirEntries, err := os.ReadDir(directory)
	if err != nil {
		err = os.Mkdir(directory, 0755)
		logger.FatalIfError(err, "Failed to create config directory at "+directory)
	}

	for _, file := range dirEntries {
		name := file.Name()
		if file.IsDir() || !(strings.HasSuffix(name, ".yaml") || strings.HasSuffix(name, ".yml")) {
			continue
		}

		filePath := filepath.Join(directory, name)
		if strings.HasSuffix(filePath, ".yaml") || strings.HasSuffix(filePath, ".yml") {
			parser.configFiles = append(parser.configFiles, filePath)
		}
	}

	rawData := parser.parseConfigFiles()
	err = yaml.Unmarshal(rawData, &parser.config)
	logger.FatalIfError(err, "Failed to parse config")
}

// Open each config file and concatenate their contents
func (parser *ConfigParser) parseConfigFiles() []byte {
	var rawData []byte

	for _, filePath := range parser.configFiles {
		data, err := os.ReadFile(filePath)
		if err != nil {
			logger.LogIfError(err, fmt.Sprintf("Failed to read config file '%s'", filePath))
			continue
		}

		rawData = append(rawData, data...)
		rawData = append(rawData, '\n')
	}

	if len(rawData) == 0 {
		logger.Log("No config data found. Write .yml config files in ~/.config/joyful")
		return nil
	}

	return rawData
}

func (parser *ConfigParser) CreateVirtualDevices() map[string]*evdev.InputDevice {
	deviceMap := make(map[string]*evdev.InputDevice)

	for _, deviceConfig := range parser.config.Devices.Virtual {
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
		for i := 0; i < numButtons - 16; i++ {
			buttons[16+i] = evdev.EvCode(startCode+i)
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