package mappingrules

import (
	"slices"

	"git.annabunches.net/annabunches/joyful/internal/logger"
	"github.com/holoplot/go-evdev"
)

func (target *RuleTargetButton) NormalizeValue(value int32) int32 {
	if value == 0 {
		return 1
	}
	return 0
}

func (target *RuleTargetAxis) NormalizeValue(value int32) int32 {
	if !target.Inverted {
		return value
	}

	axisRange := target.AxisEnd - target.AxisStart
	axisMid := target.AxisEnd - axisRange/2
	delta := value - axisMid
	if delta < 0 {
		delta = -delta
	}

	if value < axisMid {
		return axisMid + delta
	} else if value > axisMid {
		return axisMid - delta
	}

	// If we reach here, we're either exactly at the midpoint or something
	// strange has happened. Either way, just return the value.
	return value
}

// RuleTargetModeSelect doesn't make sense as an input type
func (target *RuleTargetModeSelect) NormalizeValue(value int32) int32 {
	return -1
}

func (target *RuleTargetButton) CreateEvent(value int32, mode *string) *evdev.InputEvent {
	return &evdev.InputEvent{
		Type:  evdev.EV_KEY,
		Code:  target.Code,
		Value: value,
	}
}

func (target *RuleTargetAxis) CreateEvent(value int32, mode *string) *evdev.InputEvent {
	return &evdev.InputEvent{
		Type:  evdev.EV_ABS,
		Code:  target.Code,
		Value: value,
	}
}

func (target *RuleTargetModeSelect) CreateEvent(value int32, mode *string) *evdev.InputEvent {
	if value == 0 {
		return nil
	}

	index := 0
	if currentMode := slices.Index(target.ModeSelect, *mode); currentMode != -1 {
		// find the next mode
		index = (currentMode + 1) % len(target.ModeSelect)
	}

	*mode = target.ModeSelect[index]
	logger.Logf("Mode changed to '%s'", *mode)
	return nil
}
