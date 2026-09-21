package mappingrules

import (
	"git.annabunches.net/annabunches/joyful/internal/configparser"
	"github.com/holoplot/go-evdev"
)

// A Simple Mapping Rule can map a button to a button or an axis to an axis.
type MappingRuleButton struct {
	MappingRuleBase
	Input  *RuleTargetButton
	Output *RuleTargetButton
	active bool
}

func NewMappingRuleButton(ruleConfig configparser.RuleConfigButton,
	pDevs map[string]Device,
	vDevs map[string]Device,
	base MappingRuleBase) (*MappingRuleButton, error) {

	input, err := NewRuleTargetButtonFromConfig(ruleConfig.Input, pDevs)
	if err != nil {
		return nil, err
	}

	output, err := NewRuleTargetButtonFromConfig(ruleConfig.Output, vDevs)
	if err != nil {
		return nil, err
	}

	return &MappingRuleButton{
		MappingRuleBase: base,
		Input:           input,
		Output:          output,
	}, nil
}

func (rule *MappingRuleButton) MatchEvent(device Device, event *evdev.InputEvent, mode *string) (*evdev.InputDevice, *evdev.InputEvent) {
	if !rule.MappingRuleBase.modeCheck(mode) {
		return nil, nil
	}

	if device != rule.Input.Device ||
		event.Code != rule.Input.Button {
		return nil, nil
	}

	value := rule.Input.NormalizeValue(event.Value)
	rule.active = value != 0
	return rule.Output.Device.(*evdev.InputDevice), rule.Output.CreateEvent(value, mode)
}

func (rule *MappingRuleButton) ModeChanged(newMode string, initiated bool) []OutputEvent {
	if initiated || rule.MappingRuleBase.modeMatches(newMode) {
		return nil
	}
	return rule.deactivate()
}

func (rule *MappingRuleButton) Reset(_ *string) []OutputEvent {
	return rule.deactivate()
}

func (rule *MappingRuleButton) deactivate() []OutputEvent {
	if !rule.active {
		return nil
	}
	rule.active = false
	return []OutputEvent{{Device: rule.Output.Device.(*evdev.InputDevice), Event: rule.Output.CreateEvent(0, nil)}}
}
