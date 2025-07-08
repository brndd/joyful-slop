package mappingrules

import (
	"errors"
	"slices"

	"git.annabunches.net/annabunches/joyful/internal/logger"
	"github.com/holoplot/go-evdev"
)

type RuleTargetModeSelect struct {
	Modes []string
}

func NewRuleTargetModeSelect(modes []string) (*RuleTargetModeSelect, error) {
	if len(modes) == 0 {
		return nil, errors.New("cannot create RuleTargetModeSelect: mode list is empty")
	}
	return &RuleTargetModeSelect{
		Modes: modes,
	}, nil
}

// RuleTargetModeSelect doesn't make sense as an input type
func (target *RuleTargetModeSelect) NormalizeValue(_ int32) int32 {
	return -1
}

func (target *RuleTargetModeSelect) CreateEvent(_ int32, mode *string) *evdev.InputEvent {
	index := 0
	if currentMode := slices.Index(target.Modes, *mode); currentMode != -1 {
		// find the next mode
		index = (currentMode + 1) % len(target.Modes)
	}

	*mode = target.Modes[index]
	logger.Logf("Mode changed to '%s'", *mode)
	return nil
}

func (target *RuleTargetModeSelect) MatchEvent(_ *evdev.InputDevice, _ *evdev.InputEvent) bool {
	return false
}
