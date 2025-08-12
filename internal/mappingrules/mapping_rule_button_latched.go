package mappingrules

import (
	"git.annabunches.net/annabunches/joyful/internal/configparser"
	"github.com/holoplot/go-evdev"
)

type MappingRuleButtonLatched struct {
	MappingRuleBase
	Input  *RuleTargetButton
	Output *RuleTargetButton
	State  bool
}

func NewMappingRuleButtonLatched(ruleConfig configparser.RuleConfigButtonLatched,
	pDevs map[string]Device,
	vDevs map[string]Device,
	base MappingRuleBase) (*MappingRuleButtonLatched, error) {

	input, err := NewRuleTargetButtonFromConfig(ruleConfig.Input, pDevs)
	if err != nil {
		return nil, err
	}

	output, err := NewRuleTargetButtonFromConfig(ruleConfig.Output, vDevs)
	if err != nil {
		return nil, err
	}

	return &MappingRuleButtonLatched{
		MappingRuleBase: base,
		Input:           input,
		Output:          output,
		State:           false,
	}, nil
}

func (rule *MappingRuleButtonLatched) MatchEvent(device Device, event *evdev.InputEvent, mode *string) (*evdev.InputDevice, *evdev.InputEvent) {
	if !rule.MappingRuleBase.modeCheck(mode) {
		return nil, nil
	}

	if device != rule.Input.Device ||
		event.Code != rule.Input.Button ||
		rule.Input.NormalizeValue(event.Value) == 0 {
		return nil, nil
	}

	// Input is pressed, so toggle state and emit event
	var value int32
	rule.State = !rule.State
	if rule.State {
		value = 1
	} else {
		value = 0
	}

	return rule.Output.Device.(*evdev.InputDevice), rule.Output.CreateEvent(value, mode)
}
