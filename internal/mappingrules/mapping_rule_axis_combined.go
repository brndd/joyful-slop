package mappingrules

import (
	"git.annabunches.net/annabunches/joyful/internal/configparser"
	"git.annabunches.net/annabunches/joyful/internal/logger"
	"github.com/holoplot/go-evdev"
)

type MappingRuleAxisCombined struct {
	MappingRuleBase
	InputLower *RuleTargetAxis
	InputUpper *RuleTargetAxis
	Output     *RuleTargetAxis
}

func NewMappingRuleAxisCombined(ruleConfig configparser.RuleConfigAxisCombined,
	pDevs map[string]Device,
	vDevs map[string]Device,
	base MappingRuleBase) (*MappingRuleAxisCombined, error) {

	inputLower, err := NewRuleTargetAxisFromConfig(ruleConfig.InputLower, pDevs)
	if err != nil {
		return nil, err
	}

	inputUpper, err := NewRuleTargetAxisFromConfig(ruleConfig.InputUpper, pDevs)
	if err != nil {
		return nil, err
	}

	output, err := NewRuleTargetAxisFromConfig(ruleConfig.Output, vDevs)
	if err != nil {
		return nil, err
	}

	inputLower.OutputMax = 0
	inputUpper.OutputMin = 0
	return &MappingRuleAxisCombined{
		MappingRuleBase: base,
		InputLower:      inputLower,
		InputUpper:      inputUpper,
		Output:          output,
	}, nil
}

func (rule *MappingRuleAxisCombined) MatchEvent(device Device, event *evdev.InputEvent, mode *string) (*evdev.InputDevice, *evdev.InputEvent) {
	if !rule.MappingRuleBase.modeCheck(mode) ||
		!(rule.InputLower.MatchEvent(device, event) ||
			rule.InputUpper.MatchEvent(device, event)) {

		return nil, nil
	}

	// Since lower and upper are guaranteed to have opposite signs,
	// we can just sum them.
	var value int32
	value += getValueFromAbs(rule.InputLower)
	value += getValueFromAbs(rule.InputUpper)

	return rule.Output.Device.(*evdev.InputDevice), rule.Output.CreateEvent(value, mode)
}

func getValueFromAbs(ruleTarget *RuleTargetAxis) int32 {
	absInfo, err := ruleTarget.Device.AbsInfos()
	if err != nil {
		logger.LogErrorf(err, "WARNING: Couldn't get axis data for device '%s'", ruleTarget.DeviceName)
		return 0
	}

	return ruleTarget.NormalizeValue(absInfo[ruleTarget.Axis].Value)
}
