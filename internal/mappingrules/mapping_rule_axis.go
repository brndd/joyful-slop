package mappingrules

import (
	"git.annabunches.net/annabunches/joyful/internal/configparser"
	"github.com/holoplot/go-evdev"
)

// A Simple Mapping Rule can map a button to a button or an axis to an axis.
type MappingRuleAxis struct {
	MappingRuleBase
	Input  *RuleTargetAxis
	Output *RuleTargetAxis
}

func NewMappingRuleAxis(ruleConfig configparser.RuleConfigAxis,
	pDevs map[string]Device,
	vDevs map[string]Device,
	base MappingRuleBase) (*MappingRuleAxis, error) {

	input, err := NewRuleTargetAxisFromConfig(ruleConfig.Input, pDevs)
	if err != nil {
		return nil, err
	}

	output, err := NewRuleTargetAxisFromConfig(ruleConfig.Output, vDevs)
	if err != nil {
		return nil, err
	}

	return &MappingRuleAxis{
		MappingRuleBase: base,
		Input:           input,
		Output:          output,
	}, nil
}

func (rule *MappingRuleAxis) MatchEvent(device Device, event *evdev.InputEvent, mode *string) (*evdev.InputDevice, *evdev.InputEvent) {
	if !rule.MappingRuleBase.modeCheck(mode) ||
		!rule.Input.MatchEvent(device, event) {
		return nil, nil
	}

	// The cast here is safe because the interface is only ever different for unit tests
	return rule.Output.Device.(*evdev.InputDevice), rule.Output.CreateEvent(rule.Input.NormalizeValue(event.Value), mode)
}
