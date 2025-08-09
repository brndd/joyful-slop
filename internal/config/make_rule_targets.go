package config

import (
	"errors"
	"fmt"

	"git.annabunches.net/annabunches/joyful/internal/mappingrules"
	"github.com/holoplot/go-evdev"
)

func makeRuleTargetButton(targetConfig RuleTargetConfigButton, devs map[string]Device) (*mappingrules.RuleTargetButton, error) {
	device, ok := devs[targetConfig.Device]
	if !ok {
		return nil, fmt.Errorf("non-existent device '%s'", targetConfig.Device)
	}

	eventCode, err := parseCodeButton(targetConfig.Button)
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

func makeRuleTargetAxis(targetConfig RuleTargetConfigAxis, devs map[string]Device) (*mappingrules.RuleTargetAxis, error) {
	device, ok := devs[targetConfig.Device]
	if !ok {
		return nil, fmt.Errorf("non-existent device '%s'", targetConfig.Device)
	}

	if targetConfig.DeadzoneEnd < targetConfig.DeadzoneStart {
		return nil, errors.New("deadzone_end must be greater than deadzone_start")
	}

	eventCode, err := parseCode(targetConfig.Axis, CodePrefixAxis)
	if err != nil {
		return nil, err
	}

	deadzoneStart, deadzoneEnd, err := calculateDeadzones(targetConfig, device, eventCode)
	if err != nil {
		return nil, err
	}

	return mappingrules.NewRuleTargetAxis(
		targetConfig.Device,
		device,
		eventCode,
		targetConfig.Inverted,
		deadzoneStart,
		deadzoneEnd,
	)
}

func makeRuleTargetRelaxis(targetConfig RuleTargetConfigRelaxis, devs map[string]Device) (*mappingrules.RuleTargetRelaxis, error) {
	device, ok := devs[targetConfig.Device]
	if !ok {
		return nil, fmt.Errorf("non-existent device '%s'", targetConfig.Device)
	}

	eventCode, err := parseCode(targetConfig.Axis, CodePrefixRelaxis)
	if err != nil {
		return nil, err
	}

	return mappingrules.NewRuleTargetRelaxis(
		targetConfig.Device,
		device,
		eventCode,
	)
}

func makeRuleTargetModeSelect(targetConfig RuleTargetConfigModeSelect, allModes []string) (*mappingrules.RuleTargetModeSelect, error) {
	if ok := validateModes(targetConfig.Modes, allModes); !ok {
		return nil, errors.New("undefined mode in mode select list")
	}

	return mappingrules.NewRuleTargetModeSelect(targetConfig.Modes)
}

// hasError exists solely to switch on errors in case statements
func hasError(_ any, err error) bool {
	return err != nil
}

// calculateDeadzones produces the deadzone start and end values in absolute terms
// TODO: on the one hand, this logic feels betten encapsulated in mappingrules. On the other hand,
// passing even more parameters to NewRuleTargetAxis feels terrible
func calculateDeadzones(targetConfig RuleTargetConfigAxis, device Device, axis evdev.EvCode) (int32, int32, error) {

	var deadzoneStart, deadzoneEnd int32
	deadzoneStart = 0
	deadzoneEnd = 0

	if targetConfig.DeadzoneStart != 0 || targetConfig.DeadzoneEnd != 0 {
		return targetConfig.DeadzoneStart, targetConfig.DeadzoneEnd, nil
	}

	var min, max int32
	absInfoMap, err := device.AbsInfos()

	if err != nil {
		min = mappingrules.AxisValueMin
		max = mappingrules.AxisValueMax
	} else {
		absInfo := absInfoMap[axis]
		min = absInfo.Minimum
		max = absInfo.Maximum
	}

	if targetConfig.DeadzoneCenter < min || targetConfig.DeadzoneCenter > max {
		return 0, 0, fmt.Errorf("deadzone_center '%d' is out of bounds", targetConfig.DeadzoneCenter)
	}

	switch {
	case targetConfig.DeadzoneSize != 0:
		deadzoneStart = targetConfig.DeadzoneCenter - targetConfig.DeadzoneSize/2
		deadzoneEnd = targetConfig.DeadzoneCenter + targetConfig.DeadzoneSize/2
	case targetConfig.DeadzoneSizePercent != 0:
		deadzoneSize := (max - min) / targetConfig.DeadzoneSizePercent
		deadzoneStart = targetConfig.DeadzoneCenter - deadzoneSize/2
		deadzoneEnd = targetConfig.DeadzoneCenter + deadzoneSize/2
	}

	deadzoneStart, deadzoneEnd = clampAndShift(deadzoneStart, deadzoneEnd, min, max)
	return deadzoneStart, deadzoneEnd, nil
}

func clampAndShift(start, end, min, max int32) (int32, int32) {
	if start < min {
		end += min - start
		start = min
	}
	if end > max {
		start -= end - max
		end = max
	}

	return start, end
}
