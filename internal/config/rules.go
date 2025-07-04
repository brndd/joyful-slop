package config

import (
	"fmt"
	"slices"
	"strings"

	"git.annabunches.net/annabunches/joyful/internal/logger"
	"git.annabunches.net/annabunches/joyful/internal/mappingrules"
	"github.com/holoplot/go-evdev"
)

// TODO: At some point it would *very likely* make sense to map each rule to all of the physical devices that can
// trigger it, and return that instead. Something like a map[*evdev.InputDevice][]mappingrule.MappingRule.
// This would speed up rule matching by only checking relevant rules for a given input event.
// We could take this further and make it a map[<struct of *inputdevice, type, and code>][]rule
// For very large rule-bases this may be helpful for staying performant.
func (parser *ConfigParser) BuildRules(pDevs map[string]*evdev.InputDevice, vDevs map[string]*evdev.InputDevice) []mappingrules.MappingRule {
	rules := make([]mappingrules.MappingRule, 0)
	modes := parser.GetModes()

	for _, ruleConfig := range parser.config.Rules {
		var newRule mappingrules.MappingRule
		var err error

		baseParams, err := setBaseRuleParameters(ruleConfig, vDevs, modes)
		if err != nil {
			logger.LogError(err, "couldn't set output parameters, skipping rule")
			continue
		}

		switch strings.ToLower(ruleConfig.Type) {
		case RuleTypeSimple:
			newRule, err = makeSimpleRule(ruleConfig, pDevs, baseParams)
		case RuleTypeCombo:
			newRule, err = makeComboRule(ruleConfig, pDevs, baseParams)
		case RuleTypeLatched:
			newRule, err = makeLatchedRule(ruleConfig, pDevs, baseParams)
		}

		if err != nil {
			logger.LogError(err, "failed to build rule")
			continue
		}

		rules = append(rules, newRule)
	}

	return rules
}

func setBaseRuleParameters(ruleConfig RuleConfig, vDevs map[string]*evdev.InputDevice, modes []string) (mappingrules.MappingRuleBase, error) {
	output, err := makeRuleTarget(ruleConfig.Output, vDevs)
	if err != nil {
		return mappingrules.MappingRuleBase{}, err
	}
	ruleModes := verifyModes(ruleConfig, modes)

	return mappingrules.MappingRuleBase{
		Output: output,
		Modes:  ruleModes,
		Name:   ruleConfig.Name,
	}, nil
}

func makeSimpleRule(ruleConfig RuleConfig, pDevs map[string]*evdev.InputDevice, base mappingrules.MappingRuleBase) (*mappingrules.SimpleMappingRule, error) {
	input, err := makeRuleTarget(ruleConfig.Input, pDevs)
	if err != nil {
		return nil, err
	}

	return &mappingrules.SimpleMappingRule{
		MappingRuleBase: base,
		Input:           input,
	}, nil
}

func makeComboRule(ruleConfig RuleConfig, pDevs map[string]*evdev.InputDevice, base mappingrules.MappingRuleBase) (*mappingrules.ComboMappingRule, error) {
	inputs := make([]mappingrules.RuleTarget, 0)
	for _, inputConfig := range ruleConfig.Inputs {
		input, err := makeRuleTarget(inputConfig, pDevs)
		if err != nil {
			return nil, err
		}
		inputs = append(inputs, input)
	}

	return &mappingrules.ComboMappingRule{
		MappingRuleBase: base,
		Inputs:          inputs,
		State:           0,
	}, nil
}

func makeLatchedRule(ruleConfig RuleConfig, pDevs map[string]*evdev.InputDevice, base mappingrules.MappingRuleBase) (*mappingrules.LatchedMappingRule, error) {
	input, err := makeRuleTarget(ruleConfig.Input, pDevs)
	if err != nil {
		return nil, err
	}

	return &mappingrules.LatchedMappingRule{
		MappingRuleBase: base,
		Input:           input,
		State:           false,
	}, nil
}

// makeInputRuleTarget takes an Input declaration from the YAML and returns a fully formed RuleTarget.
func makeRuleTarget(targetConfig RuleTargetConfig, devs map[string]*evdev.InputDevice) (mappingrules.RuleTarget, error) {
	if len(targetConfig.ModeSelect) > 0 {
		return &mappingrules.RuleTargetModeSelect{
			ModeSelect: targetConfig.ModeSelect,
		}, nil
	}

	device, ok := devs[targetConfig.Device]
	if !ok {
		return nil, fmt.Errorf("couldn't build rule due to non-existent device '%s'", targetConfig.Device)
	}

	eventType, eventCode, err := decodeRuleTargetValues(targetConfig)
	if err != nil {
		return nil, err
	}

	baseParams := mappingrules.RuleTargetBase{
		DeviceName: targetConfig.Device,
		Device:     device,
		Inverted:   targetConfig.Inverted,
		Code:       eventCode,
	}

	switch eventType {
	case evdev.EV_KEY:
		return &mappingrules.RuleTargetButton{
			RuleTargetBase: baseParams,
		}, nil

	case evdev.EV_ABS:
		return &mappingrules.RuleTargetAxis{
			RuleTargetBase: baseParams,
			AxisStart:      targetConfig.AxisStart,
			AxisEnd:        targetConfig.AxisEnd,
		}, nil

	default:
		return nil, fmt.Errorf("skipping rule due to unsupported event type '%d'", eventType)
	}
}

// decodeRuleTargetValues returns the appropriate evdev.EvType and evdev.EvCode values
// for a given RuleTargetConfig, converting the config file strings into appropriate constants
//
// Todo: support different formats for key specification
func decodeRuleTargetValues(target RuleTargetConfig) (evdev.EvType, evdev.EvCode, error) {
	var eventType evdev.EvType
	var eventCode evdev.EvCode
	var ok bool

	if target.Button != "" {
		eventType = evdev.EV_KEY
		eventCode, ok = evdev.KEYFromString[target.Button]
		if !ok {
			return 0, 0, fmt.Errorf("skipping rule due to invalid button code '%s'", target.Button)
		}
	} else if target.Axis != "" {
		eventType = evdev.EV_ABS
		eventCode, ok = evdev.ABSFromString[target.Axis]
		if !ok {
			return 0, 0, fmt.Errorf("skipping rule due to invalid axis code '%s'", target.Button)
		}
	}

	return eventType, eventCode, nil
}

func verifyModes(ruleConfig RuleConfig, modes []string) []string {
	verifiedModes := make([]string, 0)

	for _, configMode := range ruleConfig.Modes {
		if !slices.Contains(modes, configMode) {
			logger.Logf("rule '%s' specifies undefined mode '%s', skipping", ruleConfig.Name, configMode)
			continue
		}

		verifiedModes = append(verifiedModes, configMode)
	}
	if len(verifiedModes) == 0 {
		verifiedModes = []string{"*"}
	}

	return verifiedModes
}
