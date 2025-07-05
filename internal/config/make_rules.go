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
			logger.LogErrorf(err, "couldn't set output parameters, skipping rule '%s'", ruleConfig.Name)
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
	// We perform this check here instead of in makeRuleTarget because only Output targets
	// can meaningfully have ModeSelect; this lets us avoid plumbing the modes in on every
	// makeRuleTarget call.
	if len(ruleConfig.Output.ModeSelect) > 0 {
		err := validateModes(ruleConfig.Output.ModeSelect, modes)
		if err != nil {
			return mappingrules.MappingRuleBase{}, err
		}
	}

	output, err := makeRuleTarget(ruleConfig.Output, vDevs)
	if err != nil {
		return mappingrules.MappingRuleBase{}, err
	}

	err = validateModes(ruleConfig.Modes, modes)
	if err != nil {
		return mappingrules.MappingRuleBase{}, err
	}

	ruleModes := ensureModes(ruleConfig.Modes)

	return mappingrules.MappingRuleBase{
		Output: output,
		Modes:  ruleModes,
		Name:   ruleConfig.Name,
	}, nil
}

func makeSimpleRule(ruleConfig RuleConfig, pDevs map[string]*evdev.InputDevice, base mappingrules.MappingRuleBase) (*mappingrules.MappingRuleSimple, error) {
	input, err := makeRuleTarget(ruleConfig.Input, pDevs)
	if err != nil {
		return nil, err
	}

	return &mappingrules.MappingRuleSimple{
		MappingRuleBase: base,
		Input:           input,
	}, nil
}

func makeComboRule(ruleConfig RuleConfig, pDevs map[string]*evdev.InputDevice, base mappingrules.MappingRuleBase) (*mappingrules.MappingRuleCombo, error) {
	inputs := make([]mappingrules.RuleTarget, 0)
	for _, inputConfig := range ruleConfig.Inputs {
		input, err := makeRuleTarget(inputConfig, pDevs)
		if err != nil {
			return nil, err
		}
		inputs = append(inputs, input)
	}

	return &mappingrules.MappingRuleCombo{
		MappingRuleBase: base,
		Inputs:          inputs,
		State:           0,
	}, nil
}

func makeLatchedRule(ruleConfig RuleConfig, pDevs map[string]*evdev.InputDevice, base mappingrules.MappingRuleBase) (*mappingrules.MappingRuleLatched, error) {
	input, err := makeRuleTarget(ruleConfig.Input, pDevs)
	if err != nil {
		return nil, err
	}

	return &mappingrules.MappingRuleLatched{
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

// ensureModes either returns the mode list, or if it is empty returns []string{"*"}
func ensureModes(modes []string) []string {
	if len(modes) == 0 {
		return []string{"*"}
	}
	return modes
}

// validateModes checks the provided modes against a larger subset of modes (usually all defined ones)
// and returns an error if any of the modes are not defined.
func validateModes(modes []string, allModes []string) error {
	if len(modes) == 0 {
		return nil
	}

	for _, mode := range modes {
		if !slices.Contains(allModes, mode) {
			return fmt.Errorf("mode list specifies undefined mode '%s'", mode)
		}
	}

	return nil
}
