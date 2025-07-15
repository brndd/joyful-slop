package config

import (
	"errors"
	"fmt"

	"git.annabunches.net/annabunches/joyful/internal/mappingrules"
	"github.com/holoplot/go-evdev"
)

func makeRuleTargetButton(targetConfig RuleTargetConfig, devs map[string]*evdev.InputDevice) (*mappingrules.RuleTargetButton, error) {
	device, ok := devs[targetConfig.Device]
	if !ok {
		return nil, fmt.Errorf("non-existent device '%s'", targetConfig.Device)
	}

	eventCode, ok := evdev.KEYFromString[targetConfig.Button]
	if !ok {
		return nil, fmt.Errorf("invalid button code '%s'", targetConfig.Button)
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

	eventCode, ok := evdev.ABSFromString[targetConfig.Axis]
	if !ok {
		return nil, fmt.Errorf("invalid button code '%s'", targetConfig.Button)
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

	eventCode, ok := evdev.RELFromString[targetConfig.Axis]
	if !ok {
		return nil, fmt.Errorf("invalid button code '%s'", targetConfig.Button)
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
