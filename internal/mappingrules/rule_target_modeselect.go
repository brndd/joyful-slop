package mappingrules

import (
	"slices"

	"git.annabunches.net/annabunches/joyful/internal/logger"
	"github.com/holoplot/go-evdev"
)

type RuleTargetModeSelect struct {
	RuleTargetBase
	ModeSelect []string
}

func NewRuleTargetModeSelect(modes []string) *RuleTargetModeSelect {
	return &RuleTargetModeSelect{
		RuleTargetBase: NewRuleTargetBase("", nil, 0, false),
		ModeSelect:     modes,
	}
}

// RuleTargetModeSelect doesn't make sense as an input type
func (target *RuleTargetModeSelect) NormalizeValue(value int32) int32 {
	return -1
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
