package mappingrules

import (
	"git.annabunches.net/annabunches/joyful/internal/configparser"
	"github.com/holoplot/go-evdev"
)

type MappingRuleModeSelect struct {
	MappingRuleBase
	Input  *RuleTargetButton
	Output *RuleTargetModeSelect
}

func NewMappingRuleModeSelect(ruleConfig configparser.RuleConfigModeSelect,
	pDevs map[string]Device,
	modes []string,
	base MappingRuleBase) (*MappingRuleModeSelect, error) {

	input, err := NewRuleTargetButtonFromConfig(ruleConfig.Input, pDevs)
	if err != nil {
		return nil, err
	}

	output, err := NewRuleTargetModeSelectFromConfig(ruleConfig.Output, modes)
	if err != nil {
		return nil, err
	}

	return &MappingRuleModeSelect{
		MappingRuleBase: base,
		Input:           input,
		Output:          output,
	}, nil
}

func (rule *MappingRuleModeSelect) MatchEvent(
	device Device,
	event *evdev.InputEvent,
	mode *string) (*evdev.InputDevice, *evdev.InputEvent) {

	if !rule.MappingRuleBase.modeCheck(mode) {
		return nil, nil
	}

	if device != rule.Input.Device ||
		event.Code != rule.Input.Button ||
		rule.Input.NormalizeValue(event.Value) == 0 {
		return nil, nil
	}

	return nil, rule.Output.CreateEvent(event.Value, mode)
}
