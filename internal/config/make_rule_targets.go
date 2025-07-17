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

	eventCode, err := parseCode(targetConfig.Button, "BTN")
	if err != nil {
		return nil, err
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

	eventCode, err := parseCode(targetConfig.Axis, "ABS")
	if err != nil {
		return nil, err
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

	eventCode, err := parseCode(targetConfig.Axis, "REL")
	if err != nil {
		return nil, err
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
