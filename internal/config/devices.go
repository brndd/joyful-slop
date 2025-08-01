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
		capabilities := map[evdev.EvType][]evdev.EvCode{
			evdev.EV_KEY: makeButtons(deviceConfig.NumButtons, deviceConfig.Buttons),
			evdev.EV_ABS: makeAxes(deviceConfig.NumAxes, deviceConfig.Axes),
			evdev.EV_REL: makeRelativeAxes(deviceConfig.NumRelativeAxes, deviceConfig.RelativeAxes),
		}

		device, err := evdev.CreateDevice(
			name,
			// TODO: who knows what these should actually be
			evdev.InputID{
				BusType: 0x03,
				Vendor:  0x4711,
				Product: 0x0816,
				Version: 1,
			},
			capabilities,
		)

		if err != nil {
			logger.LogIfError(err, "Failed to create virtual device")
			continue
		}

		deviceMap[deviceConfig.Name] = device
		logger.Log(fmt.Sprintf(
			"Created virtual device '%s' with %d buttons, %d axes, and %d relative axes",
			name,
			len(capabilities[evdev.EV_KEY]),
			len(capabilities[evdev.EV_ABS]),
			len(capabilities[evdev.EV_REL]),
		))
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
func (parser *ConfigParser) ConnectPhysicalDevices(lock bool) map[string]*evdev.InputDevice {
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

		if lock {
			err := device.Grab()
			if err != nil {
				logger.LogError(err, "Failed to grab device for exclusive access")
			}
		}

		logger.Log(fmt.Sprintf("Connected to '%s' as '%s'", deviceConfig.DeviceName, deviceConfig.Name))
		deviceMap[deviceConfig.Name] = device
	}

	return deviceMap
}

// TODO: these functions have a lot of duplication; we need to figure out how to refactor it cleanly
// without losing logging context...
func makeButtons(numButtons int, buttonList []string) []evdev.EvCode {
	if numButtons > 0 && len(buttonList) > 0 {
		logger.Log("'num_buttons' and 'buttons' both specified, ignoring 'num_buttons'")
	}

	if numButtons > VirtualDeviceMaxButtons {
		numButtons = VirtualDeviceMaxButtons
		logger.Logf("Limiting virtual device buttons to %d", VirtualDeviceMaxButtons)
	}

	if len(buttonList) > 0 {
		buttons := make([]evdev.EvCode, 0, len(buttonList))
		for _, codeStr := range buttonList {
			code, err := parseCode(codeStr, "BTN")
			if err != nil {
				logger.LogError(err, "Failed to create button, skipping")
				continue
			}
			buttons = append(buttons, code)
		}
		return buttons
	}

	buttons := make([]evdev.EvCode, numButtons)

	for i := 0; i < numButtons; i++ {
		buttons[i] = ButtonFromIndex[i]
	}

	return buttons
}

func makeAxes(numAxes int, axisList []string) []evdev.EvCode {
	if numAxes > 0 && len(axisList) > 0 {
		logger.Log("'num_axes' and 'axes' both specified, ignoring 'num_axes'")
	}

	if len(axisList) > 0 {
		axes := make([]evdev.EvCode, 0, len(axisList))
		for _, codeStr := range axisList {
			code, err := parseCode(codeStr, "ABS")
			if err != nil {
				logger.LogError(err, "Failed to create axis, skipping")
				continue
			}
			axes = append(axes, code)
		}
		return axes
	}

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

func makeRelativeAxes(numAxes int, axisList []string) []evdev.EvCode {
	if numAxes > 0 && len(axisList) > 0 {
		logger.Log("'num_rel_axes' and 'rel_axes' both specified, ignoring 'num_rel_axes'")
	}

	if len(axisList) > 0 {
		axes := make([]evdev.EvCode, 0, len(axisList))
		for _, codeStr := range axisList {
			code, err := parseCode(codeStr, "REL")
			if err != nil {
				logger.LogError(err, "Failed to create axis, skipping")
				continue
			}
			axes = append(axes, code)
		}
		return axes
	}

	if numAxes > 10 {
		numAxes = 10
		logger.Log("Limiting virtual device relative axes to 10")
	}

	axes := make([]evdev.EvCode, numAxes)
	for i := 0; i < numAxes; i++ {
		axes[i] = evdev.EvCode(i)
	}

	return axes
}
