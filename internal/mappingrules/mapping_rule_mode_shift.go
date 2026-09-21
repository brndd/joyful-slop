package mappingrules

import (
	"errors"

	"git.annabunches.net/annabunches/joyful/internal/configparser"
	"git.annabunches.net/annabunches/joyful/internal/logger"
	"github.com/holoplot/go-evdev"
)

type MappingRuleModeShift struct {
	MappingRuleBase
	Input        *RuleTargetButton
	Mode         string
	previousMode string
	active       bool
}

func (*MappingRuleModeShift) SilentModeChange() {}

func NewMappingRuleModeShift(ruleConfig configparser.RuleConfigModeShift, pDevs map[string]Device, modes []string, base MappingRuleBase) (*MappingRuleModeShift, error) {
	if !validateModes([]string{ruleConfig.Mode}, modes) {
		return nil, errors.New("mode-shift specifies undefined mode")
	}
	input, err := NewRuleTargetButtonFromConfig(ruleConfig.Input, pDevs)
	if err != nil {
		return nil, err
	}
	return &MappingRuleModeShift{MappingRuleBase: base, Input: input, Mode: ruleConfig.Mode}, nil
}

func (rule *MappingRuleModeShift) MatchEvent(device Device, event *evdev.InputEvent, mode *string) (*evdev.InputDevice, *evdev.InputEvent) {
	if !rule.Input.MatchEvent(device, event) {
		return nil, nil
	}
	value := rule.Input.NormalizeValue(event.Value)
	if value != 0 {
		if rule.active || !rule.MappingRuleBase.modeCheck(mode) {
			return nil, nil
		}
		rule.previousMode = *mode
		rule.active = true
		*mode = rule.Mode
		logger.Logf("Mode changed to '%s'", *mode)
	} else if rule.active {
		rule.active = false
		*mode = rule.previousMode
		logger.Logf("Mode changed to '%s'", *mode)
	}
	return nil, nil
}
