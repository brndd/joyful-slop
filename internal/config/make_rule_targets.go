package config

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"git.annabunches.net/annabunches/joyful/internal/mappingrules"
	"github.com/holoplot/go-evdev"
)

func makeRuleTargetButton(targetConfig RuleTargetConfig, devs map[string]*evdev.InputDevice) (*mappingrules.RuleTargetButton, error) {
	device, ok := devs[targetConfig.Device]
	if !ok {
		return nil, fmt.Errorf("non-existent device '%s'", targetConfig.Device)
	}

	var eventCode evdev.EvCode
	buttonConfig := strings.ToUpper(targetConfig.Button)
	switch {
	case strings.HasPrefix(buttonConfig, "BTN_"):
		eventCode, ok = evdev.KEYFromString[buttonConfig]
		if !ok {
			return nil, fmt.Errorf("invalid button specification '%s'", buttonConfig)
		}

	case strings.HasPrefix(buttonConfig, "0X"):
		codeInt, err := strconv.ParseUint(buttonConfig[2:], 16, 0)
		if err != nil {
			return nil, err
		}
		eventCode = evdev.EvCode(codeInt)

	case !hasError(strconv.Atoi(buttonConfig)):
		index, err := strconv.Atoi(buttonConfig)
		if err != nil {
			return nil, err
		}

		if index >= len(ButtonFromIndex) {
			return nil, fmt.Errorf("button index '%d' out of bounds", index)
		}

		eventCode = ButtonFromIndex[index]

	default:
		eventCode, ok = evdev.KEYFromString["BTN_"+buttonConfig]
		if !ok {
			return nil, fmt.Errorf("invalid button specification '%s'", buttonConfig)
		}
	}

	return mappingrules.NewRuleTargetButton(
		targetConfig.Device,
		device,
		eventCode,
		targetConfig.Inverted,
	)
}

func makeRuleTargetAxis(targetConfig RuleTargetConfig, devs map[string]*evdev.InputDevice) (*mappingrules.RuleTargetAxis, error) {
	device, ok := devs[targetConfig.Device]
	if !ok {
		return nil, fmt.Errorf("non-existent device '%s'", targetConfig.Device)
	}

	if targetConfig.DeadzoneEnd < targetConfig.DeadzoneStart {
		return nil, errors.New("deadzone_end must be greater than deadzone_start")
	}

	var eventCode evdev.EvCode
	axisConfig := strings.ToUpper(targetConfig.Axis)
	switch {
	case strings.HasPrefix(axisConfig, "ABS_"):
		eventCode, ok = evdev.ABSFromString[axisConfig]
		if !ok {
			return nil, fmt.Errorf("invalid axis code '%s'", axisConfig)
		}

	case strings.HasPrefix(axisConfig, "0X"):
		codeInt, err := strconv.ParseUint(axisConfig[2:], 16, 32)
		if err != nil {
			return nil, err
		}
		eventCode = evdev.EvCode(codeInt)

	default:
		eventCode, ok = evdev.ABSFromString["ABS_"+axisConfig]
		if !ok {
			return nil, fmt.Errorf("invalid axis code '%s'", axisConfig)
		}
	}

	return mappingrules.NewRuleTargetAxis(
		targetConfig.Device,
		device,
		eventCode,
		targetConfig.Inverted,
		targetConfig.DeadzoneStart,
		targetConfig.DeadzoneEnd,
	)
}

func makeRuleTargetRelaxis(targetConfig RuleTargetConfig, devs map[string]*evdev.InputDevice) (*mappingrules.RuleTargetRelaxis, error) {
	device, ok := devs[targetConfig.Device]
	if !ok {
		return nil, fmt.Errorf("non-existent device '%s'", targetConfig.Device)
	}

	var eventCode evdev.EvCode
	axisConfig := strings.ToUpper(targetConfig.Axis)
	switch {
	case strings.HasPrefix(axisConfig, "REL_"):
		eventCode, ok = evdev.RELFromString[axisConfig]
		if !ok {
			return nil, fmt.Errorf("invalid axis code '%s'", axisConfig)
		}

	case strings.HasPrefix(axisConfig, "0X"):
		codeInt, err := strconv.ParseUint(axisConfig[2:], 16, 32)
		if err != nil {
			return nil, err
		}
		eventCode = evdev.EvCode(codeInt)

	default:
		eventCode, ok = evdev.RELFromString["REL_"+axisConfig]
		if !ok {
			return nil, fmt.Errorf("invalid axis code '%s'", axisConfig)
		}
	}
	return mappingrules.NewRuleTargetRelaxis(
		targetConfig.Device,
		device,
		eventCode,
		targetConfig.Inverted,
	)
}

func makeRuleTargetModeSelect(targetConfig RuleTargetConfig, allModes []string) (*mappingrules.RuleTargetModeSelect, error) {
	if ok := validateModes(targetConfig.Modes, allModes); !ok {
		return nil, errors.New("undefined mode in mode select list")
	}

	return mappingrules.NewRuleTargetModeSelect(targetConfig.Modes)
}

// hasError exists solely to switch on errors in case statements
func hasError(_ any, err error) bool {
	return err != nil
}
