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

	return rule.Output.Device.(*evdev.InputDevice), rule.Output.CreateEvent(rule.Input.NormalizeValue(event.Value), mode)
}
