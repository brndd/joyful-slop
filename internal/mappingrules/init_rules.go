package mappingrules

import (
	"errors"
	"fmt"
	"slices"

	"git.annabunches.net/annabunches/joyful/internal/configparser"
	"git.annabunches.net/annabunches/joyful/internal/logger"
	"github.com/holoplot/go-evdev"
)

func ConvertDeviceMap(inputDevs map[string]*evdev.InputDevice) map[string]Device {
	// Golang can't inspect the concrete map type to determine interface conformance,
	// so we handle that here.
	devices := make(map[string]Device)
	for name, dev := range inputDevs {
		devices[name] = dev
	}
	return devices
}

// NewRule parses a RuleConfig struct and creates and returns the appropriate rule type.
// You can remap a map[string]*evdev.InputDevice to our interface type with ConvertDeviceMap
func NewRule(config configparser.RuleConfig, pDevs map[string]Device, vDevs map[string]Device, modes []string) (MappingRule, error) {
	var newRule MappingRule
	var err error

	if !validateModes(config.Modes, modes) {
		return nil, errors.New("mode list specifies undefined mode")
	}

	base := NewMappingRuleBase(config.Name, config.Modes)

	switch config.Type {
	case configparser.RuleTypeButton:
		newRule, err = NewMappingRuleButton(config.Config.(configparser.RuleConfigButton), pDevs, vDevs, base)
	case configparser.RuleTypeButtonCombo:
		newRule, err = NewMappingRuleButtonCombo(config.Config.(configparser.RuleConfigButtonCombo), pDevs, vDevs, base)
	case configparser.RuleTypeButtonLatched:
		newRule, err = NewMappingRuleButtonLatched(config.Config.(configparser.RuleConfigButtonLatched), pDevs, vDevs, base)
	case configparser.RuleTypeButtonTempo:
		newRule, err = NewMappingRuleButtonTempo(config.Config.(configparser.RuleConfigButtonTempo), pDevs, vDevs, modes, base)
	case configparser.RuleTypeAxis:
		newRule, err = NewMappingRuleAxis(config.Config.(configparser.RuleConfigAxis), pDevs, vDevs, base)
	case configparser.RuleTypeAxisCombined:
		newRule, err = NewMappingRuleAxisCombined(config.Config.(configparser.RuleConfigAxisCombined), pDevs, vDevs, base)
	case configparser.RuleTypeAxisToButton:
		newRule, err = NewMappingRuleAxisToButton(config.Config.(configparser.RuleConfigAxisToButton), pDevs, vDevs, base)
	case configparser.RuleTypeAxisToRelaxis:
		newRule, err = NewMappingRuleAxisToRelaxis(config.Config.(configparser.RuleConfigAxisToRelaxis), pDevs, vDevs, base)
	case configparser.RuleTypeModeSelect:
		newRule, err = NewMappingRuleModeSelect(config.Config.(configparser.RuleConfigModeSelect), pDevs, modes, base)
	case configparser.RuleTypeModeShift:
		newRule, err = NewMappingRuleModeShift(config.Config.(configparser.RuleConfigModeShift), pDevs, modes, base)
	case configparser.RuleTypeHat:
		newRule, err = NewMappingRuleHat(config.Config.(configparser.RuleConfigHat), pDevs, vDevs, base)
	default:
		// Shouldn't actually be possible to get here...
		err = fmt.Errorf("bad rule type '%s' for rule '%s'", config.Type, config.Name)
	}

	if err != nil {
		logger.LogErrorf(err, "Failed to build rule '%s'", config.Name)
		return nil, err
	}

	return newRule, nil
}

// validateModes checks the provided modes against a larger subset of modes (usually all defined ones)
// and returns false if any of the modes are not defined.
func validateModes(modes []string, allModes []string) bool {
	if len(modes) == 0 {
		return true
	}

	for _, mode := range modes {
		if !slices.Contains(allModes, mode) {
			return false
		}
	}

	return true
}
