package config

import (
	"fmt"
	"strings"

	"git.annabunches.net/annabunches/joyful/internal/logger"
	"github.com/holoplot/go-evdev"
)

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

		name := fmt.Sprintf("joyful-%s", deviceConfig.Name)
		device, err := evdev.CreateDevice(
			name,
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

		deviceMap[deviceConfig.Name] = device
		logger.Log(fmt.Sprintf("Created virtual device '%s'", name))
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
func (parser *ConfigParser) ConnectPhysicalDevices() map[string]*evdev.InputDevice {
	deviceMap := make(map[string]*evdev.InputDevice)

	for _, deviceConfig := range parser.config.Devices {
		if strings.ToLower(deviceConfig.Type) != DeviceTypePhysical {
			continue
		}

		device, err := evdev.OpenByName(deviceConfig.DeviceName)
		if err != nil {
			logger.LogError(err, "Failed to open physical device, skipping. Confirm the device name with 'evlist'. Watch out for spaces.")
			continue
		}

		// TODO: grab exclusive access to device (add config option)

		logger.Log(fmt.Sprintf("Connected to '%s' as '%s'", deviceConfig.DeviceName, deviceConfig.Name))
		deviceMap[deviceConfig.Name] = device
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
