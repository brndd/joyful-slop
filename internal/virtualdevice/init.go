package virtualdevice

import (
	"fmt"

	"git.annabunches.net/annabunches/joyful/internal/configparser"
	"git.annabunches.net/annabunches/joyful/internal/eventcodes"
	"git.annabunches.net/annabunches/joyful/internal/logger"
	"github.com/holoplot/go-evdev"
	"strconv"
)

func ParseUInt16(s string, def uint16) uint16 {
    if s == "" {
        return def
    }
    v, err := strconv.ParseUint(s, 0, 16)
    if err != nil {
        return def
    }
    return uint16(v)
}

// NewEventBuffer takes a virtual device config specification, creates the underlying
// evdev.InputDevice, and wraps it in a buffered event emitter.
func NewEventBuffer(config configparser.DeviceConfigVirtual) (*EventBuffer, error) {
	deviceMap := make(map[string]*evdev.InputDevice)

	name := fmt.Sprintf("joyful-%s", config.Name)

	var capabilities map[evdev.EvType][]evdev.EvCode

	// todo: add tests for presets
	switch config.Preset {
	case DevicePresetGamepad:
		capabilities = CapabilitiesPresetGamepad
	case DevicePresetKeyboard:
		capabilities = CapabilitiesPresetKeyboard
	case DevicePresetJoystick:
		capabilities = CapabilitiesPresetJoystick
	case DevicePresetMouse:
		capabilities = CapabilitiesPresetMouse
	default:
		capabilities = map[evdev.EvType][]evdev.EvCode{
			evdev.EV_KEY: makeButtons(config.NumButtons, config.Buttons),
			evdev.EV_ABS: makeAxes(config.NumAxes, config.Axes),
			evdev.EV_REL: makeRelativeAxes(config.NumRelativeAxes, config.RelativeAxes),
		}
	}

	device, err := evdev.CreateDevice(
		name,
		// TODO: placeholders. BusType = differentiate between input bus (3 = usb), Vendor = uint16 vendor id, Product: uint16 devoce id, Version = device version (differentiate between differing versions of exact same device)
		evdev.InputID{
			BusType: 0x03,
			Vendor:  ParseUInt16(config.VendorId, VirtualDeviceVendorId),
			Product: ParseUInt16(config.DeviceId, VirtualDeviceDeviceId),
			Version: 1,
		},
		capabilities,
	)

	if err != nil {
		return nil, err
	}

	deviceMap[config.Name] = device
	logger.Log(fmt.Sprintf(
		"Created virtual device '%s' with %d buttons, %d axes, and %d relative axes",
		name,
		len(capabilities[evdev.EV_KEY]),
		len(capabilities[evdev.EV_ABS]),
		len(capabilities[evdev.EV_REL]),
	))

	return &EventBuffer{
		events: make([]*evdev.InputEvent, 0, 100),
		Device: device,
		Name:   config.Name,
	}, nil
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
			code, err := eventcodes.ParseCode(codeStr, "BTN")
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
		buttons[i] = eventcodes.ButtonFromIndex[i]
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
			code, err := eventcodes.ParseCode(codeStr, "ABS")
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
			code, err := eventcodes.ParseCode(codeStr, "REL")
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
